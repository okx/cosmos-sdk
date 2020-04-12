package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// creates a querier for staking REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		//case types.QueryValidators:
		//	return queryValidators(ctx, req, k)
		//case types.QueryValidator:
		//	return queryValidator(ctx, req, k)
		//case types.QueryValidatorDelegations:
		//	return queryValidatorDelegations(ctx, req, k)
		//case types.QueryValidatorUnbondingDelegations:
		//	return queryValidatorUnbondingDelegations(ctx, req, k)
		//case types.QueryDelegation:
		//	return queryDelegation(ctx, req, k)
		//case types.QueryUnbondingDelegation:
		//	return queryUnbondingDelegation(ctx, req, k)
		//case types.QueryDelegatorDelegations:
		//	return queryDelegatorDelegations(ctx, req, k)
		//case types.QueryDelegatorUnbondingDelegations:
		//	return queryDelegatorUnbondingDelegations(ctx, req, k)
		//case types.QueryRedelegations:
		//	return queryRedelegations(ctx, req, k)
		//case types.QueryDelegatorValidators:
		//	return queryDelegatorValidators(ctx, req, k)
		//case types.QueryDelegatorValidator:
		//	return queryDelegatorValidator(ctx, req, k)
		//case types.QueryPool:
		//	return queryPool(ctx, k)
		//case types.QueryParameters:
		//	return queryParameters(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown staking query endpoint")
		}
	}
}
