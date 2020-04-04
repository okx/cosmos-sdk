package types

import (
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MinterCustom struct {
	NextBlockToUpdate uint64         `json:"next_block_to_update" yaml:"next_block_to_update"` // record the block height for next year
	AnnualProvisions  sdk.Dec        `json:"annual_provisions" yaml:"annual_provisions"`       // record the amount of Annual minted
	MintedPerBlock    types.DecCoins `json:"minted_per_block" yaml:"minted_per_block"`         // record the MintedPerBlock per block in this year
}
