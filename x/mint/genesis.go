package mint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint/internal/keeper"
)

// GenesisState - minter state
type GenesisState struct {
	Minter             MinterCustom `json:"minter_custom" yaml:"minter_custom"` // minter object
	Params             Params       `json:"params" yaml:"params"`               // inflation params
	OriginalMintedPerBlock sdk.Dec      `json:"original_minted_per_block" yaml:"original_minted_per_block"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(minter MinterCustom, params Params, originalMintedPerBlock sdk.Dec) GenesisState {
	return GenesisState{
		Minter:             minter,
		Params:             params,
		OriginalMintedPerBlock: originalMintedPerBlock,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Minter:             DefaultInitialMinterCustom(),
		Params:             DefaultParams(),
		OriginalMintedPerBlock: keeper.DefaultOriginalMintedPerBlock(),
	}
}

// InitGenesis new mint genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetMinterCustom(ctx, data.Minter)
	keeper.SetParams(ctx, data.Params)
	keeper.SetOriginalMintedPerBlock(data.OriginalMintedPerBlock)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	minter := keeper.GetMinterCustom(ctx)
	params := keeper.GetParams(ctx)
	return NewGenesisState(minter, params, keeper.GetOriginalMintedPerBlock())
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	err := ValidateParams(data.Params)
	if err != nil {
		return err
	}

	err = ValidateMinterCustom(data.Minter)
	if err != nil {
		return err
	}

	err = keeper.ValidateOriginalMintedPerBlock(data.OriginalMintedPerBlock)
	if err != nil {
		return err
	}
	return nil
}
