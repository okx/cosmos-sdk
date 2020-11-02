package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/cosmos/cosmos-sdk/x/supply/internal/types"
)

// Keeper of the supply store
type Keeper struct {
	cdc       *codec.Codec
	storeKey  sdk.StoreKey
	ak        types.AccountKeeper
	bk        types.BankKeeper
	permAddrs map[string]types.PermissionsForAddress
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, ak types.AccountKeeper, bk types.BankKeeper, maccPerms map[string][]string) Keeper {
	// set the addresses
	permAddrs := make(map[string]types.PermissionsForAddress)
	for name, perms := range maccPerms {
		permAddrs[name] = types.NewPermissionsForAddress(name, perms)
	}

	return Keeper{
		cdc:       cdc,
		storeKey:  key,
		ak:        ak,
		bk:        bk,
		permAddrs: permAddrs,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetSupply retrieves the Supply from store
//func (k Keeper) GetSupply(ctx sdk.Context) (supply exported.SupplyI) {
//	store := ctx.KVStore(k.storeKey)
//	b := store.Get(SupplyKey)
//	if b == nil {
//		panic("stored supply should not have been nil")
//	}
//	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &supply)
//	return
//}

// SetSupply sets the Supply to store
//func (k Keeper) SetSupply(ctx sdk.Context, supply exported.SupplyI) {
//	store := ctx.KVStore(k.storeKey)
//	b := k.cdc.MustMarshalBinaryLengthPrefixed(supply)
//	store.Set(SupplyKey, b)
//}

// ValidatePermissions validates that the module account has been granted
// permissions within its set of allowed permissions.
func (k Keeper) ValidatePermissions(macc exported.ModuleAccountI) error {
	permAddr := k.permAddrs[macc.GetName()]
	for _, perm := range macc.GetPermissions() {
		if !permAddr.HasPermission(perm) {
			return fmt.Errorf("invalid module permission %s", perm)
		}
	}
	return nil
}

// GetTotalSupply gets all the total supply of tokens in the ledger
func (k Keeper) GetTotalSupply(ctx sdk.Context) (totalSupply sdk.DecCoins) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PrefixTokenSupplyKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var amount sdk.Dec
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &amount)
		totalSupply = append(totalSupply, sdk.NewDecCoinFromDec(string(iterator.Key()[1:]), amount))
	}

	return
}

// Inflate adds the amount of a token in the store
func (k Keeper) inflate(ctx sdk.Context, tokenSymbol string, amount sdk.Dec) {
	if amount.Equal(sdk.ZeroDec()) {
		return
	}

	originalSupplyAmount := k.GetTokenSupplyAmount(ctx, tokenSymbol)
	k.SetTokenSupplyAmount(ctx, tokenSymbol, originalSupplyAmount.Add(amount))
}

// Deflate subtracts the amount of a token from the original in the store
func (k Keeper) deflate(ctx sdk.Context, tokenSymbol string, deflationAmount sdk.Dec) sdk.Error {
	currentSupplyAmount := k.GetTokenSupplyAmount(ctx, tokenSymbol)
	supplyAmount := currentSupplyAmount.Sub(deflationAmount)
	if supplyAmount.IsNegative() {
		return types.ErrInvalidDeflation(types.DefaultCodespace, deflationAmount, currentSupplyAmount, tokenSymbol)
	}

	k.SetTokenSupplyAmount(ctx, tokenSymbol, supplyAmount)
	return nil
}

// GetTokenSupplyAmount gets the amount of a token supply from the store
func (k Keeper) GetTokenSupplyAmount(ctx sdk.Context, tokenSymbol string) sdk.Dec {
	tokenSupplyAmount := sdk.ZeroDec()
	bytes := ctx.KVStore(k.storeKey).Get(types.GetTokenSupplyKey(tokenSymbol))
	if bytes == nil {
		return tokenSupplyAmount
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bytes, &tokenSupplyAmount)
	return tokenSupplyAmount
}

// SetTokenSupplyAmount sets the supply amount of a token to the store
func (k Keeper) SetTokenSupplyAmount(ctx sdk.Context, tokenSymbol string, supplyAmount sdk.Dec) {
	ctx.KVStore(k.storeKey).Set(types.GetTokenSupplyKey(tokenSymbol), k.cdc.MustMarshalBinaryLengthPrefixed(supplyAmount))
}
