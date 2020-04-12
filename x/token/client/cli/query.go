package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/token/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	stakingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the staking module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	//stakingQueryCmd.AddCommand(client.GetCommands(
	//	GetCmdQueryDelegation(queryRoute, cdc),
	//	GetCmdQueryDelegations(queryRoute, cdc),
	//	GetCmdQueryUnbondingDelegation(queryRoute, cdc),
	//	GetCmdQueryUnbondingDelegations(queryRoute, cdc),
	//	GetCmdQueryRedelegation(queryRoute, cdc),
	//	GetCmdQueryRedelegations(queryRoute, cdc),
	//	GetCmdQueryPool(queryRoute, cdc))...)

	return stakingQueryCmd

}

// GetCmdQueryValidator implements the validator query command.
