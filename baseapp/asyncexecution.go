package baseapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type ExecuteResult struct {
	Resp    abci.ResponseDeliverTx
	Ms      sdk.CacheMultiStore
	Counter uint32
}

func (e ExecuteResult) GetResponse() abci.ResponseDeliverTx {
	return e.Resp
}

func (e ExecuteResult) Recheck() bool {
	return true
}

func (e ExecuteResult) GetCounter() uint32 {
	return e.Counter
}

func (e ExecuteResult) Commit() bool {
	e.Ms.Write()
	return true
}

func NewExecuteResult(r abci.ResponseDeliverTx, ms sdk.CacheMultiStore, counter uint32) ExecuteResult {
	return ExecuteResult{
		Resp:    r,
		Ms:      ms,
		Counter: counter,
	}
}

type AsyncWorkGroup struct {
	WorkCh     chan ExecuteResult
	ExecRes    []abci.ExecuteRes
	MaxCounter int
	Cb         abci.AsyncCallBack
}

func NewAsyncWorkGroup() *AsyncWorkGroup {
	return &AsyncWorkGroup{
		WorkCh:     make(chan ExecuteResult, 1),
		ExecRes:    make([]abci.ExecuteRes, 0),
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
		select {
		case exec = <-a.WorkCh:
			a.ExecRes = append(a.ExecRes, exec)
			if len(a.ExecRes) == a.MaxCounter {
				//call tendermint to update the deliver tx response
				if a.Cb != nil {
					a.Cb(a.ExecRes)
				}
			}
		}
	}()
}

func (a *AsyncWorkGroup) IncreaseCounter() {
	a.MaxCounter++
}

func (a *AsyncWorkGroup) Reset() {
	a.ExecRes = make([]abci.ExecuteRes, 0)
	a.MaxCounter = 0
}
