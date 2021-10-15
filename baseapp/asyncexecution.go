package baseapp

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

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

func (a *AsyncWorkGroup) Reset() {
	//a.WorkCh = make(chan ExecuteResult, 300)
}
