package baseapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (app *BaseApp) PushAnteHandler(ah sdk.AnteHandler) {
	app.anteHandler = ah
}

func (app *BaseApp) GetState(mode runTxMode) *state {
	if mode == runTxModeCheck || mode == runTxModeSimulate {
		return app.checkState
	}

	return app.deliverState
}
