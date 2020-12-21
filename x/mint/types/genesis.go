package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// NewGenesisState creates a new GenesisState object
func NewGenesisState(minter MinterCustom, params Params, originalMintedPerBlock sdk.Dec) *GenesisState {
	return &GenesisState{
		Minter: minter,
		Params: params,

		OriginalMintedPerBlock: originalMintedPerBlock,
	}
}

func DefaultOriginalMintedPerBlock() sdk.Dec {
	return sdk.MustNewDecFromStr("1")
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Minter: DefaultInitialMinterCustom(),
		Params: DefaultParams(),
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	return ValidateMinterCustom(data.Minter)
}
