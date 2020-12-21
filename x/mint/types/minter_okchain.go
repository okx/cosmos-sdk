package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMinterCustom returns a new Minter object with the given inflation and annual
// provisions values.
func NewMinterCustom(nextBlockToUpdate uint64, mintedPerBlock sdk.DecCoins) MinterCustom {
	return MinterCustom{
		NextBlockToUpdate: nextBlockToUpdate,
		MintedPerBlock:    mintedPerBlock,
	}
}

// InitialMinterCustom returns an initial Minter object with a given inflation value.
func InitialMinterCustom() MinterCustom {
	return NewMinterCustom(
		0,
		sdk.DecCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.ZeroInt())},
	)
}

// DefaultInitialMinterCustom returns a default initial MinterCustom object for a new chain
// which uses an inflation rate of 1%.
func DefaultInitialMinterCustom() MinterCustom {
	return InitialMinterCustom()
}

// ValidateMinterCustom validate minter
func ValidateMinterCustom(minter MinterCustom) error {
	if len(minter.MintedPerBlock) != 1 || minter.MintedPerBlock[0].Denom != sdk.DefaultBondDenom {
		return fmt.Errorf(" MintedPerBlock must contain only %s", sdk.DefaultBondDenom)
	}
	return nil
}
