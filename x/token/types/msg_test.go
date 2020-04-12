package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMsg(t *testing.T)  {

	addr, _ := sdk.AccAddressFromHex("eva1hg40dv5e237qy28vtyum52ygke32ez35syykpz")
	msg := NewMsgTokenBurn(
		sdk.NewCoins(sdk.NewCoin("okt", sdk.NewInt(5))), addr)

	s1 := string(msg.GetSignBytes())
	s2 := string(msg.GetSignBytes2())

	_ = s1
	_ = s2
}

