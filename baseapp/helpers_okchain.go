package baseapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (app *BaseApp) PushAnteHandler(ah sdk.AnteHandler) {
	app.anteHandler = ah
}
