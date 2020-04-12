package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/token/types"
	"github.com/spf13/cobra"
	"github.com/pkg/errors"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)


const (
	From          = "from"
)

var (
	errFromNotValid           = errors.New("from not valid")
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	stakingTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Token transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	stakingTxCmd.AddCommand(client.PostCommands(
		//GetCmdCreateValidator(cdc),
		//GetCmdEditValidator(cdc),
		GetCmdTokenBurn(cdc),
	)...)

	return stakingTxCmd
}


// GetCmdSetName is the CLI command for sending a SetName transaction
func GetCmdTokenBurn(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [tokens]",
		Short: "burn some amount of token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := authTypes.NewAccountRetriever(cliCtx).EnsureExists(cliCtx.FromAddress); err != nil {
				return err
			}

			// params check
			flags := cmd.Flags()

			_, err := flags.GetString(From)
			if err != nil {
				return errFromNotValid
			}

			coins, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}
			msg := types.NewMsgTokenBurn(coins, cliCtx.FromAddress)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
	//cmd.LocalFlags()
	return cmd
}
