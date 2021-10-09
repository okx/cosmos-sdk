package baseapp

import (
	"encoding/hex"
	"fmt"
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

func (e ExecuteResult) Recheck(cache abci.AsyncCacheInterface) bool {
	rerun := false
	if e.Ms == nil {
		return true
	}

	e.Ms.IteratorCache(func(key, value []byte, isDirty bool) bool {
		//the key we have read was wrote by pre txs
		if cache.Has(key) && !whiteAccountList[hex.EncodeToString(key)] {
			fmt.Println("conflict", hex.EncodeToString(key), string(key))
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

func (e ExecuteResult) Commit() {
	if e.Ms != nil {
		e.Ms.IteratorCache(func(key, value []byte, isDirty bool) bool {
			if isDirty {
				//fmt.Println("ok.scf.debug", hex.EncodeToString(key), hex.EncodeToString(value))
			}
			return true
		})
		e.Ms.Write()
	} else {
		// TODO delete
		panic("need panic")
	}

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
