package baseapp

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	sm "github.com/tendermint/tendermint/state"
	"sync"
	"time"
)

func (app *BaseApp) deliverTxsWithParallel(group map[int][]int, nextTx map[int]int) []*abci.ResponseDeliverTx {
	txs := app.parallelTxManage.txs
	txsBytes := app.parallelTxManage.txsByte

	var validTxs, invalidTxs = 0, 0
	txIndex := 0
	txReps := make([]abci.ExecuteRes, len(app.parallelTxManage.txStatus))
	asCache := NewAsyncCache()
	signal := make(chan int, 1)
	rerunIdx := 0

	tsCommi := time.Duration(0)

	deliverTxsResponse := make([]*abci.ResponseDeliverTx, len(txs), len(txs))
	AsyncCb := func(execRes abci.ExecuteRes) {
		txReps[execRes.GetCounter()] = execRes
		//fmt.Println("zhixingwanbi", execRes.GetCounter(), execRes.GetBase())
		for txReps[txIndex] != nil {
			res := txReps[txIndex]
			cc := res.Conflict(asCache)
			//fmt.Println("tcInfr", txIndex, cc, res.GetCounter(), res.GetBase())
			if cc {
				rerunIdx++
				txReps[txIndex] = nil
				go app.parallelTx(txIndex)
				//check proxy.err?
				return
			}
			txRs := res.GetResponse()
			deliverTxsResponse[txIndex] = &txRs
			tsQ := time.Now()
			res.Collect(asCache)
			res.Commit()
			tsCommi += time.Now().Sub(tsQ)
			if deliverTxsResponse[txIndex].Code == abci.CodeTypeOK {
				validTxs++
			} else {
				invalidTxs++
			}

			app.parallelTxManage.currMergeIndex = txIndex
			if nextTxIndex, ok := nextTx[txIndex]; ok {
				go app.parallelTx(nextTxIndex)
			}
			txIndex++
			//fmt.Println("current", txIndex)
			if txIndex == len(txsBytes) {
				app.logger.Info(fmt.Sprintf("BlockHeight %d With Tx %d : Paralle run %d, Conflected tx %d CommitTs %d",
					app.LastBlockHeight(), len(txsBytes), len(deliverTxsResponse)-rerunIdx, rerunIdx, tsCommi.Microseconds()))
				signal <- 0
				return
			}
		}
	}

	app.parallelTxManage.workgroup.Cb = AsyncCb

	for index := 0; index < len(group); index++ {
		go app.parallelTx(group[index][0])
	}

	if len(txsBytes) > 0 {
		//waiting for call back
		<-signal
		//CheckErr
		receiptsLogs := app.EndParallelTxs()
		for index, v := range receiptsLogs {
			if len(v) != 0 { // only update evm tx result
				deliverTxsResponse[index].Data = v
			}
		}
	}
	return deliverTxsResponse
}

func (app *BaseApp) PrepareParallelTxs(txs [][]byte) []*abci.ResponseDeliverTx {
	ts := time.Now()
	//app.parallelTxManage.isAsyncDeliverTx = true
	sendAccs := make([]ethcmn.Address, 0)
	toAccs := make([]*ethcmn.Address, 0)
	evmIndex := uint32(0)
	for k, v := range txs {
		tx, err := app.txDecoder(v)
		if err != nil {
			panic(err)
		}
		t := &txStatus{
			indexInBlock: uint32(k),
		}
		fee, isEvm, singerCache, from, to := app.getTxFee(app.getContextForTx(runTxModeDeliverInAsync, v), tx)
		if isEvm {
			t.evmIndex = evmIndex
			t.isEvmTx = true
			evmIndex++
		}
		sendAccs = append(sendAccs, from)
		toAccs = append(toAccs, to)
		app.parallelTxManage.singerCaches[string(v)] = singerCache

		vString := string(v)
		app.parallelTxManage.txs = append(app.parallelTxManage.txs, tx)
		app.parallelTxManage.txsByte = append(app.parallelTxManage.txsByte, v)
		app.parallelTxManage.SetFee(vString, fee)

		app.parallelTxManage.txStatus[vString] = t
		app.parallelTxManage.indexMapBytes = append(app.parallelTxManage.indexMapBytes, vString)
	}

	if !viper.GetBool(sm.FlagParalleledTx) {
		return nil

	}

	//fmt.Println("zhunbei fenzu")
	groupList, nextTxInGroup := grouping(sendAccs, toAccs)
	fmt.Println("grouping", time.Now().Sub(ts).Microseconds())
	//fmt.Println("fenzu", len(groupList), groupList, nextTxInGroup)
	res := app.deliverTxsWithParallel(groupList, nextTxInGroup)
	fmt.Println("run-end", time.Now().Sub(ts).Microseconds())
	//fmt.Println("RRRRRR", len(res))
	return res
}

func (app *BaseApp) EndParallelTxs() [][]byte {
	txFeeInBlock := sdk.Coins{}
	feeMap := app.parallelTxManage.GetFeeMap()
	refundMap := app.parallelTxManage.GetRefundFeeMap()
	for tx, v := range feeMap {
		if app.parallelTxManage.txStatus[tx].anteErr != nil {
			continue
		}
		txFeeInBlock = txFeeInBlock.Add(v...)
		if refundFee, ok := refundMap[tx]; ok {
			txFeeInBlock = txFeeInBlock.Sub(refundFee)
		}
	}
	ctx, cache := app.cacheTxContext(app.getContextForTx(runTxModeDeliverInAsync, []byte{}), []byte{})
	if err := app.updateFeeCollectorAccHandler(ctx, txFeeInBlock); err != nil {
		panic(err)
	}
	cache.Write()

	txExecStats := make([][]string, 0)
	for _, v := range app.parallelTxManage.indexMapBytes {
		errMsg := ""
		if err := app.parallelTxManage.txStatus[v].anteErr; err != nil {
			errMsg = err.Error()
		}
		txExecStats = append(txExecStats, []string{v, errMsg})
	}
	app.parallelTxManage.Clear()
	return app.logFix(txExecStats)
}

func (app *BaseApp) parallelTx(index int) {
	if app.parallelTxManage.GetRunningStats(index) {
		return
	}
	if index <= app.parallelTxManage.currMergeIndex {
		return
	}
	defer app.parallelTxManage.SetRunningStats(index, false)
	//fmt.Println("palllllll", index, app.parallelTxManage.currMergeIndex)

	txBytes := app.parallelTxManage.txsByte[index]
	txStd := app.parallelTxManage.txs[index]
	mergedIndex := app.parallelTxManage.currMergeIndex
	txStatus := app.parallelTxManage.txStatus[string(txBytes)]

	if !txStatus.isEvmTx {
		asyncExe := NewExecuteResult(abci.ResponseDeliverTx{}, nil, txStatus.indexInBlock, txStatus.evmIndex, mergedIndex)
		app.parallelTxManage.workgroup.Push(asyncExe)
	}

	go func() {
		mergedIndex := mergedIndex
		var resp abci.ResponseDeliverTx
		g, r, m, e := app.runTx(runTxModeDeliverInAsync, txBytes, txStd, LatestSimulateTxHeight)
		if e != nil {
			resp = sdkerrors.ResponseDeliverTx(e, 0, 0, app.trace)
		} else {
			resp = abci.ResponseDeliverTx{
				GasWanted: int64(g.GasWanted), // TODO: Should type accept unsigned ints?
				GasUsed:   int64(g.GasUsed),   // TODO: Should type accept unsigned ints?
				Log:       r.Log,
				Data:      r.Data,
				Events:    r.Events.ToABCIEvents(),
			}
		}

		txStatus := app.parallelTxManage.txStatus[string(txBytes)]
		asyncExe := NewExecuteResult(resp, m, txStatus.indexInBlock, txStatus.evmIndex, mergedIndex)
		asyncExe.err = e
		app.parallelTxManage.workgroup.Push(asyncExe)
	}()
}
func (app *BaseApp) DeliverTxWithCache(req abci.RequestDeliverTx) abci.ExecuteRes { //TODO delete return
	panic("need delete")
}

type ExecuteResult struct {
	Resp       abci.ResponseDeliverTx
	Ms         sdk.CacheMultiStore
	Counter    uint32
	err        error
	evmCounter uint32
	base       int
}

func (e ExecuteResult) GetResponse() abci.ResponseDeliverTx {
	return e.Resp
}

func (e ExecuteResult) Conflict(cache abci.AsyncCacheInterface) bool {
	if e.base+1 == int(e.Counter) {
		return false
	}
	rerun := false
	if e.Ms == nil {
		return true //TODO fix later
	}

	e.Ms.IteratorCache(func(key, value []byte, isDirty bool) bool {
		//the key we have read was wrote by pre txs
		if cache.Has(e.base, int(e.Counter), key) && !whiteAccountList[hex.EncodeToString(key)] {
			rerun = true
			return false
		}
		return true
	})
	return rerun
}

var (
	whiteAccountList = map[string]bool{
		"01f1829676db577682e944fc3493d451b67ff3e29f": true, //fee
	}
)

func (e ExecuteResult) Collect(cache abci.AsyncCacheInterface) {
	if e.Ms == nil {
		return
	}
	e.Ms.IteratorCache(func(key, value []byte, isDirty bool) bool {
		if isDirty {
			//push every data we have written in current tx
			cache.Push(int(e.Counter), key, value)
		}
		return true
	})
}

func (e ExecuteResult) GetCounter() uint32 {
	return e.Counter
}

func (e ExecuteResult) GetBase() int {
	return e.base
}

func (e ExecuteResult) Commit() {
	if e.Ms == nil {
		return
	}
	e.Ms.Write()
}

func NewExecuteResult(r abci.ResponseDeliverTx, ms sdk.CacheMultiStore, counter uint32, evmCounter uint32, base int) ExecuteResult {
	return ExecuteResult{
		Resp:       r,
		Ms:         ms,
		Counter:    counter,
		evmCounter: evmCounter,
		base:       base,
	}
}

type AsyncWorkGroup struct {
	WorkCh chan ExecuteResult
	Cb     abci.AsyncCallBack
}

func NewAsyncWorkGroup() *AsyncWorkGroup {
	return &AsyncWorkGroup{
		WorkCh: make(chan ExecuteResult, 64),
		Cb:     nil,
	}
}

func (a *AsyncWorkGroup) Push(item ExecuteResult) {
	a.WorkCh <- item
}

func (a *AsyncWorkGroup) Start() {
	go func() {
		for {
			select {
			case exec := <-a.WorkCh:
				if a.Cb != nil {
					a.Cb(exec)
				}
			}
		}
	}()
}

type parallelTxManager struct {
	mu               sync.RWMutex
	isAsyncDeliverTx bool
	workgroup        *AsyncWorkGroup

	fee       map[string]sdk.Coins
	refundFee map[string]sdk.Coins

	txStatus      map[string]*txStatus
	indexMapBytes []string

	txs            []sdk.Tx
	txsByte        [][]byte
	singerCaches   map[string]sdk.SigCache
	currMergeIndex int

	isRunning map[int]bool
}

type txStatus struct {
	isEvmTx      bool
	evmIndex     uint32
	indexInBlock uint32
	anteErr      error
}

func newParallelTxManager() *parallelTxManager {
	return &parallelTxManager{
		isAsyncDeliverTx: false,
		workgroup:        NewAsyncWorkGroup(),
		fee:              make(map[string]sdk.Coins),
		refundFee:        make(map[string]sdk.Coins),

		txStatus:       make(map[string]*txStatus),
		indexMapBytes:  make([]string, 0),
		txs:            make([]sdk.Tx, 0),
		txsByte:        make([][]byte, 0),
		singerCaches:   make(map[string]sdk.SigCache, 0),
		currMergeIndex: -1,
		isRunning:      make(map[int]bool),
	}
}

func (f *parallelTxManager) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fee = make(map[string]sdk.Coins)
	f.refundFee = make(map[string]sdk.Coins)

	f.txStatus = make(map[string]*txStatus)
	f.indexMapBytes = make([]string, 0)
	f.txs = make([]sdk.Tx, 0)
	f.txsByte = make([][]byte, 0)
	f.singerCaches = make(map[string]sdk.SigCache, 0)
	f.currMergeIndex = -1
	f.isRunning = make(map[int]bool, 0)

}
func (f *parallelTxManager) SetFee(key string, value sdk.Coins) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fee[key] = value
}

func (f *parallelTxManager) GetFeeMap() map[string]sdk.Coins {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.fee
}
func (f *parallelTxManager) SetRefundFee(key string, value sdk.Coins) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.refundFee[key] = value
}

func (f *parallelTxManager) GetRefundFeeMap() map[string]sdk.Coins {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.refundFee
}

func (f *parallelTxManager) SetRunningStats(txIndex int, status bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.isRunning[txIndex] = status
}
func (f *parallelTxManager) GetRunningStats(txIndex int) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.isRunning[txIndex]
}
