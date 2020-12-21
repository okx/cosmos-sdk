package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

type ValidateMsgHandler func(ctx sdk.Context, msgs []sdk.Msg) sdk.Result

type IsSystemFreeHandler func(ctx sdk.Context, msgs []sdk.Msg) bool

type ObserverI interface {
	OnAccountUpdated(acc types.AccountI)
}

func (ak *AccountKeeper) SetObserverKeeper(observer ObserverI) {
	ak.observer = observer
}
