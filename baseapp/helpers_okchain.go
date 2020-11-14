package baseapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func (app *BaseApp) GetCommitMultiStore() sdk.CommitMultiStore {
	return app.cms
}

func (app *BaseApp) GetDeliverStateCtx() sdk.Context {
	return app.deliverState.ctx
}

//-------------------------------------------------------
// for protocol engine to invoke
//-------------------------------------------------------
func (app *BaseApp) PushInitChainer(initChainer sdk.InitChainer) {
	app.initChainer = initChainer
}

func (app *BaseApp) PushBeginBlocker(beginBlocker sdk.BeginBlocker) {
	app.beginBlocker = beginBlocker
}

func (app *BaseApp) PushEndBlocker(endBlocker sdk.EndBlocker) {
	app.endBlocker = endBlocker
}

func (app *BaseApp) PushAnteHandler(ah sdk.AnteHandler) {
	app.anteHandler = ah
}

func (app *BaseApp) SetTxDecoder(txDecoder sdk.TxDecoder) {
	app.txDecoder = txDecoder
}
