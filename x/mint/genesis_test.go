package mint

import (
	"github.com/cosmos/cosmos-sdk/x/mint/internal/keeper"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {

	setup := newTestInput(t)

	genesisState := NewGenesisState(
		DefaultInitialMinterCustom(), DefaultParams(), keeper.DefaultInitTokensPerBlock())

	defaultGenesisState := DefaultGenesisState()

	require.Equal(t, genesisState, defaultGenesisState)

	InitGenesis(setup.ctx, setup.mintKeeper, defaultGenesisState)

	require.NoError(t, ValidateGenesis(defaultGenesisState))

	exportedState := ExportGenesis(setup.ctx, setup.mintKeeper)

	require.NotNil(t, exportedState)

	require.Equal(t, defaultGenesisState, exportedState)
}
