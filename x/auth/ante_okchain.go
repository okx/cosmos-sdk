package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	SystemFee = "0.0125"
)

func GetSystemFee() sdk.Coin {
	return sdk.MustParseCoins(sdk.DefaultBondDenom, SystemFee)[0]
}
func ZeroFee() sdk.Coin {
	return sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt())
}

func GetSysFeeCoins() sdk.Coins {
	return sdk.Coins{GetSystemFee()}
}

type ValidateMsgHandler func(ctx sdk.Context, msgs []sdk.Msg) sdk.Result

type IsSystemFreeHandler func(ctx sdk.Context, msgs []sdk.Msg) bool
