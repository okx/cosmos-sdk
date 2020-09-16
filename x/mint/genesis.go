package mint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint/internal/keeper"
)

// GenesisState - minter state
type GenesisState struct {
	Minter             MinterCustom `json:"minter_custom" yaml:"minter_custom"` // minter object
	Params             Params       `json:"params" yaml:"params"`               // inflation params
	InitTokensPerBlock sdk.Dec      `json:"init_tokens_per_block" yaml:"init_tokens_per_block"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(minter MinterCustom, params Params, initTokensPerBlock sdk.Dec) GenesisState {
	return GenesisState{
		Minter:             minter,
		Params:             params,
		InitTokensPerBlock: initTokensPerBlock,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Minter:             DefaultInitialMinterCustom(),
		Params:             DefaultParams(),
		InitTokensPerBlock: keeper.DefaultInitTokensPerBlock(),
	}
}

// InitGenesis new mint genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetMinterCustom(ctx, data.Minter)
	keeper.SetParams(ctx, data.Params)
	keeper.SetInitTokensPerBlock(data.InitTokensPerBlock)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	minter := keeper.GetMinterCustom(ctx)
	params := keeper.GetParams(ctx)
	return NewGenesisState(minter, params, keeper.GetInitTokensPerBlock())
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

	err = keeper.ValidateInitTokensPerBlock(data.InitTokensPerBlock)
	if err != nil {
		return err
	}
	return nil
}
