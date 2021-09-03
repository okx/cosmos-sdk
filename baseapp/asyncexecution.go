package baseapp

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type ExecuteResult struct {
	Resp       abci.ResponseDeliverTx
	Ms         sdk.CacheMultiStore
	Counter    uint32
	err        error
	reAnte     bool
	evmCounter uint32
}

func (e ExecuteResult) GetResponse() abci.ResponseDeliverTx {
	return e.Resp
}

func (e ExecuteResult) Recheck(cache abci.AsyncCacheInterface) bool {
	rerun := false
	if e.reAnte {
		//if ante failed, it means the same `from` address has sent multi tx in one block , ante may using wrong nonce
		//so we rerun this tx in directly
		return true
	}
	if e.Ms == nil { //means tx was failed, nothing need to commit
		return false
	}
	e.Ms.IteratorCache(func(key, value []byte, isDirty bool) bool {
		//the key we have read was wrote by pre txs
		if cache.Has(key) {
			rerun = true
		}
		return true
	})

	return rerun
}

func (e ExecuteResult) Collect(cache abci.AsyncCacheInterface) {
	if e.Ms == nil {
		return
	}
	e.Ms.IteratorCache(func(key, value []byte, isDirty bool) bool {
		if isDirty {
			//push every data we have wrote in current tx
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

func (e ExecuteResult) Commit() bool {
	if e.Ms == nil {
		fmt.Println("commiting a nil res")
		return false
	}
	e.Ms.Write()
	return true
}

func (e ExecuteResult) GetEvmTxCounter() uint32 {
	return e.evmCounter
}

func (e ExecuteResult) NeedAnte() bool {
	return e.reAnte
}

func NewExecuteResult(r abci.ResponseDeliverTx, ms sdk.CacheMultiStore, counter uint32, evmCounter uint32) ExecuteResult {
	return ExecuteResult{
		Resp:       r,
		Ms:         ms,
		Counter:    counter,
		reAnte:     false,
		evmCounter: evmCounter,
	}
}

type AsyncWorkGroup struct {
	WorkCh     chan ExecuteResult
	ExecRes    map[int]abci.ExecuteRes
	MaxCounter int
	Cb         abci.AsyncCallBack
}

func NewAsyncWorkGroup() *AsyncWorkGroup {
	return &AsyncWorkGroup{
		WorkCh:     make(chan ExecuteResult, 1),
		ExecRes:    make(map[int]abci.ExecuteRes, 0),
		MaxCounter: 0,
		Cb:         nil,
	}
}

func (a *AsyncWorkGroup) Push(item ExecuteResult) {
	a.WorkCh <- item
}

func (a *AsyncWorkGroup) Start() {
	go func() {
		var exec ExecuteResult
		for {
			select {
			case exec = <-a.WorkCh:
				a.ExecRes[int(exec.Counter)] = exec
				if len(a.ExecRes) == a.MaxCounter {
					//call tendermint to update the deliver tx response
					if a.Cb != nil {
						a.Cb(a.ExecRes)
					}
				}
			}
		}
	}()
}

func (a *AsyncWorkGroup) SetMaxCounter(MaxCounter int) {
	a.MaxCounter = MaxCounter
}

func (a *AsyncWorkGroup) Reset() {
	a.ExecRes = make(map[int]abci.ExecuteRes, 0)
	a.MaxCounter = 0
}
