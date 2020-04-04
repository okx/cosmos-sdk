package utils

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

)


func SearchTxs(cliCtx context.CLIContext, cdc *codec.Codec, tags []string, page, limit int) ([]sdk.TxResponse, error) {

	res, err := QueryTxsByEvents(cliCtx, tags, page, limit)
	if err != nil {
		return nil, err
	}

	return res.Txs, nil
}