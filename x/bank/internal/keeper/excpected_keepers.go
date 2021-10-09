package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	ethvm "github.com/ethereum/go-ethereum/core/vm"
)

const (
	COSMOS_CALL_TYPE     = "cosmos"
	SEND_CALL_NAME       = "send"
	DELEGATE_CALL_NAME   = "delegate"
	MULTI_CALL_NAME      = "multi-send"
	UNDELEGATE_CALL_NAME = "undelegate"
	COSMOS_DEPTH         = 0
)

type InnerTxKeeper interface {
	GetInnerBlockData() ethvm.BlockInnerData
	InitInnerBlock(hash string)
	AddInnerTx(hash string, txs []*ethvm.InnerTx)
	GetCodec() *codec.Codec
}
