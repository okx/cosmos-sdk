package keeper

import (
	"container/list"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/token/types"
)

const aminoCacheSize = 500

// keeper of the staking store
type Keeper struct {
	storeKey           sdk.StoreKey
	storeTKey          sdk.StoreKey
	cdc                *codec.Codec
	supplyKeeper       types.SupplyKeeper
	paramstore         params.Subspace
	validatorCacheList *list.List

	bankKeeper       bank.Keeper

	feeCollectorName string // name of the FeeCollector ModuleAccount
	// The reference to the Param Keeper to get and set Global Params
	paramsKeeper params.Keeper

	// The reference to the Paramstore to get and set gov specific params
	paramSpace     params.Subspace
	tokenStoreKey  sdk.StoreKey // Unexposed key to access name store from sdk.Context
	freezeStoreKey sdk.StoreKey
	lockStoreKey   sdk.StoreKey

	// codespace
	codespace sdk.CodespaceType
}


// NewKeeper creates a new staking Keeper instance
func NewKeeper(cdc *codec.Codec, key, tkey sdk.StoreKey,
	bankKeeper     bank.Keeper,
	supplyKeeper types.SupplyKeeper,
	paramstore params.Subspace, codespace sdk.CodespaceType) Keeper {

	// ensure bonded and not bonded module accounts are set
	if addr := supplyKeeper.GetModuleAddress(types.BondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	if addr := supplyKeeper.GetModuleAddress(types.NotBondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	return Keeper{
		storeKey:           key,
		storeTKey:          tkey,
		cdc:                cdc,
		bankKeeper:         bankKeeper,
		supplyKeeper:       supplyKeeper,
		paramstore:         paramstore.WithKeyTable(ParamKeyTable()),
		validatorCacheList: list.New(),
		codespace:          codespace,
	}
}

func (k Keeper) SupplyKeeper() types.SupplyKeeper {
	return k.supplyKeeper
}

func (k Keeper) BankKeeper() bank.Keeper {
	return k.bankKeeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}


// return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// Load the last total validator power.
func (k Keeper) GetLastTotalPower(ctx sdk.Context) (power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.LastTotalPowerKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &power)
	return
}

// Set the last total validator power.
func (k Keeper) SetLastTotalPower(ctx sdk.Context, power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(power)
	store.Set(types.LastTotalPowerKey, b)
}
