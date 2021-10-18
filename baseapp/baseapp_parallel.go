package baseapp

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"sync"
	"unsafe"
)

func (app *BaseApp) FinalTx() [][]byte {
	txFeeInBlock := sdk.Coins{}
	feeMap := app.feeManage.GetFeeMap()
	refundMap := app.feeManage.GetRefundFeeMap()
	for tx, v := range feeMap {
		if app.feeManage.txDetail[tx].AnteErr != nil {
			continue
		}
		txFeeInBlock = txFeeInBlock.Add(v...)
		if refundFee, ok := refundMap[tx]; ok {
			txFeeInBlock = txFeeInBlock.Sub(refundFee)
		}
	}
	ctx, cache := app.cacheTxContext(app.getContextForTx(runTxModeDeliverInAsync, []byte{}), []byte{})
	app.feeCollectorAccHandler(ctx, true, txFeeInBlock.Add(app.initPoolCoins...))
	cache.Write()

	tmp := make([][]string, 0)
	for _, v := range app.feeManage.indexMapBytes {
		errMsg := ""
		if err := app.feeManage.txDetail[v].AnteErr; err != nil {
			errMsg = err.Error()
		}
		tmp = append(tmp, []string{v, errMsg})
	}

	evmReceipts := app.fixLog(tmp)
	res := make([][]byte, 0)
	txLen := len(app.feeManage.txDetail)
	for index := 0; index < txLen; index++ {
		res = append(res, evmReceipts[index])
	}
	app.feeManage.Clear()
	return res
}

//we reuse the nonce that changed by the last async call
//if last ante handler has been failed, we need rerun it ? or not?
func (app *BaseApp) DeliverTxWithCache(req abci.RequestDeliverTx) abci.ExecuteRes {
	tx, err := app.txDecoder(req.Tx)
	if err != nil {
		return nil
	}
	var (
		gInfo sdk.GasInfo
		resp  abci.ResponseDeliverTx
		mode  runTxMode
	)
	mode = runTxModeDeliverInAsync
	g, r, m, e := app.runTx(mode, req.Tx, tx, LatestSimulateTxHeight)
	if e != nil {
		resp = sdkerrors.ResponseDeliverTx(e, gInfo.GasWanted, gInfo.GasUsed, app.trace)
	} else {
		resp = abci.ResponseDeliverTx{
			GasWanted: int64(g.GasWanted), // TODO: Should type accept unsigned ints?
			GasUsed:   int64(g.GasUsed),   // TODO: Should type accept unsigned ints?
			Log:       r.Log,
			Data:      r.Data,
			Events:    r.Events.ToABCIEvents(),
		}
	}

	txx := bytes2str(req.Tx)

	asyncExe := NewExecuteResult(resp, m, app.feeManage.txDetail[txx].IndexInBlock, app.feeManage.txDetail[txx].EvmIndex)
	asyncExe.err = e
	return asyncExe
}

func (app *BaseApp) SetAsyncDeliverTxCb(cb abci.AsyncCallBack) {
	app.workgroup.Cb = cb
}

func (app *BaseApp) SetAsyncConfig(sw bool, txs [][]byte) {
	app.isAsyncDeliverTx = true
	app.initPoolCoins = app.feeCollectorAccHandler(app.getContextForTx(runTxModeDeliverInAsync, nil), false, sdk.Coins{})

	evmIndex := uint32(0)
	for k, v := range txs {
		tx, err := app.txDecoder(v)
		if err != nil {
			panic(err)
		}
		t := &txIndex{
			IndexInBlock: uint32(k),
		}
		fee, isEvm := app.getTxFee(tx)
		if isEvm {
			t.EvmIndex = evmIndex
			t.isEvmTx = true
			evmIndex++
		}

		app.feeManage.SetFee(bytes2str(v), fee)
		vString := bytes2str(v)

		app.feeManage.txDetail[vString] = t
		app.feeManage.indexMapBytes = append(app.feeManage.indexMapBytes, vString)
	}

}

type ExecuteResult struct {
	Resp       abci.ResponseDeliverTx
	Ms         sdk.CacheMultiStore
	Counter    uint32
	err        error
	evmCounter uint32
}

func (e ExecuteResult) GetResponse() abci.ResponseDeliverTx {
	return e.Resp
}

func (e ExecuteResult) Conflict(cache abci.AsyncCacheInterface) bool {
	rerun := false
	if e.Ms == nil {
		return true //TODO fix later
	}

	e.Ms.IteratorCache(func(key, value []byte, isDirty bool) bool {
		//the key we have read was wrote by pre txs
		if cache.Has(key) && !whiteAccountList[hex.EncodeToString(key)] {
			rerun = true
		}
		return true
	})
	return rerun
}

var (
	whiteAccountList = map[string]bool{
		//"676c6f62616c4163636f756e744e756d626572":     true, //globalAccountNumber
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
			cache.Push(key, value)
		}
		return true
	})
}

func (e ExecuteResult) Error() error {
	return e.err
}

func (e ExecuteResult) GetCounter() uint32 {
	return e.Counter
}

func (e ExecuteResult) Commit() {
	if e.Ms == nil {
		return
	}
	e.Ms.Write()
}

func (e ExecuteResult) GetEvmTxCounter() uint32 {
	return e.evmCounter
}

func NewExecuteResult(r abci.ResponseDeliverTx, ms sdk.CacheMultiStore, counter uint32, evmCounter uint32) ExecuteResult {
	return ExecuteResult{
		Resp:       r,
		Ms:         ms,
		Counter:    counter,
		evmCounter: evmCounter,
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

type feeManager struct {
	mu sync.RWMutex

	fee       map[string]sdk.Coins
	refundFee map[string]sdk.Coins

	txDetail      map[string]*txIndex
	indexMapBytes []string
}

type txIndex struct {
	isEvmTx      bool
	EvmIndex     uint32
	IndexInBlock uint32
	AnteErr      error
}

func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func NewFeeManager() *feeManager {
	return &feeManager{
		fee:       make(map[string]sdk.Coins),
		refundFee: make(map[string]sdk.Coins),

		txDetail:      make(map[string]*txIndex),
		indexMapBytes: make([]string, 0),
	}
}

func (f *feeManager) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fee = make(map[string]sdk.Coins)
	f.refundFee = make(map[string]sdk.Coins)

	f.txDetail = make(map[string]*txIndex)
	f.indexMapBytes = make([]string, 0)

}
func (f *feeManager) SetFee(key string, value sdk.Coins) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fee[key] = value
}

func (f *feeManager) GetFeeMap() map[string]sdk.Coins {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.fee
}
func (f *feeManager) SetRefundFee(key string, value sdk.Coins) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.refundFee[key] = value
}

func (f *feeManager) GetRefundFeeMap() map[string]sdk.Coins {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.refundFee
}
