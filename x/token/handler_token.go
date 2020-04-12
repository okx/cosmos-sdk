package token

import (
	"fmt"
	_ "strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/token/types"
)


func handleMsgTokenBurn(ctx sdk.Context, keeper Keeper, msg types.MsgTokenBurn) sdk.Result {

	subCoins := msg.Amount
	// check balance
	myCoins := keeper.BankKeeper().GetCoins(ctx, msg.Owner)
	for _,  existing:= range myCoins {
		fmt.Printf("%v - %v \n", existing, msg.Amount)

		for _, in := range msg.Amount {

			if existing.Denom != in.Denom {
				continue
			}

			if existing.Amount.LT(in.Amount) {
				return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient coins(need %s)", existing.String())).Result()
			}
			break
		}
	}

	fmt.Printf("subCoins: %v \n", subCoins)

	// send coins to moduleAcc
	err := keeper.SupplyKeeper().SendCoinsFromAccountToModule(ctx, msg.Owner, types.ModuleName, subCoins)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("supply send coins error:%s", err.Error())).Result()
	}

	// set supply
	err = keeper.SupplyKeeper().BurnCoins(ctx, types.ModuleName, subCoins)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("supply burn coins error:%s", err.Error())).Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}
