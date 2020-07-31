package baseapp

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/params"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Paramspace defines the parameter subspace to be used for the paramstore.
const Paramspace = "baseapp"

// Parameter store keys for all the consensus parameter types.
var (
	ParamStoreKeyBlockParams     = []byte("BlockParams")
	ParamStoreKeyEvidenceParams  = []byte("EvidenceParams")
	ParamStoreKeyValidatorParams = []byte("ValidatorParams")
)

// ParamStore defines the interface the parameter store used by the BaseApp must
// fulfill.
type ParamStore interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Has(ctx sdk.Context, key []byte) bool
	Set(ctx sdk.Context, key []byte, param interface{})
}

// ValidateBlockParams defines a stateless validation on BlockParams. This function
// is called whenever the parameters are updated or stored.
func ValidateBlockParams(i interface{}) error {
	v, ok := i.(abci.BlockParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.MaxBytes <= 0 {
		return fmt.Errorf("block maximum bytes must be positive: %d", v.MaxBytes)
	}
	if v.MaxGas < -1 {
		return fmt.Errorf("block maximum gas must be greater than or equal to -1: %d", v.MaxGas)
	}

	return nil
}

// ValidateEvidenceParams defines a stateless validation on EvidenceParams. This
// function is called whenever the parameters are updated or stored.
func ValidateEvidenceParams(i interface{}) error {
	v, ok := i.(abci.EvidenceParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.MaxAge <= 0 {
		return fmt.Errorf("evidence maximum age in blocks must be positive: %d", v.MaxAge)
	}


	return nil
}

// ValidateValidatorParams defines a stateless validation on ValidatorParams. This
// function is called whenever the parameters are updated or stored.
func ValidateValidatorParams(i interface{}) error {
	v, ok := i.(abci.ValidatorParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(v.PubKeyTypes) == 0 {
		return errors.New("validator allowed pubkey types must not be empty")
	}

	return nil
}

// ConsensusParamsKeyTable returns an x/params module keyTable to be used in
// the BaseApp's ParamStore. The KeyTable registers the types along with the
// standard validation functions. Applications can choose to adopt this KeyTable
// or provider their own when the existing validation functions do not suite their
// needs.
func ConsensusParamsKeyTable() params.KeyTable {
	return params.NewKeyTable(
		params.ParamSetPair{
			Key: ParamStoreKeyBlockParams, Value: abci.BlockParams{},
		},
		params.ParamSetPair{
			Key: ParamStoreKeyEvidenceParams, Value: abci.EvidenceParams{},
		},
		params.ParamSetPair{
			Key: ParamStoreKeyValidatorParams, Value: abci.ValidatorParams{},
		},
	)
}