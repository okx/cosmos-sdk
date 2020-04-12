package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)


//func (k Keeper) NewToken(ctx sdk.Context, token types.Token) {
//	// save token info
//	store := ctx.KVStore(k.storeKey)
//	store.Set(types.GetTokenAddress(token.Symbol), k.cdc.MustMarshalBinaryBare(token))
//
//	// update token number
//	var tokenNumber uint64
//	b := store.Get(types.TokenNumberKey)
//	if b == nil {
//		tokenNumber = 0
//	} else {
//		k.cdc.MustUnmarshalBinaryBare(b, &tokenNumber)
//	}
//	b = k.cdc.MustMarshalBinaryBare(tokenNumber + 1)
//	store.Set(types.TokenNumberKey, b)
//}


func (k Keeper) GetNumKeys(ctx sdk.Context) (tokenStoreKeyNum, freezeStoreKeyNum, lockStoreKeyNum int64) {
	{
		store := ctx.KVStore(k.tokenStoreKey)
		iter := store.Iterator(nil, nil)
		defer iter.Close()
		for ; iter.Valid(); iter.Next() {
			tokenStoreKeyNum++
		}
	}
	{
		store := ctx.KVStore(k.freezeStoreKey)
		iter := store.Iterator(nil, nil)
		defer iter.Close()
		for ; iter.Valid(); iter.Next() {
			freezeStoreKeyNum++
		}
	}
	{
		store := ctx.KVStore(k.lockStoreKey)
		iter := store.Iterator(nil, nil)
		defer iter.Close()
		for ; iter.Valid(); iter.Next() {
			lockStoreKeyNum++
		}
	}

	return
}

