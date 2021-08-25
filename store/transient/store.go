package transient

import (
	"context"
	"github.com/tendermint/iavl"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/store/types"

	"github.com/cosmos/cosmos-sdk/store/dbadapter"
)

var _ types.Committer = (*Store)(nil)
var _ types.KVStore = (*Store)(nil)

// Store is a wrapper for a MemDB with Commiter implementation
type Store struct {
	dbadapter.Store
}

// Constructs new MemDB adapter
func NewStore() *Store {
	return &Store{Store: dbadapter.Store{DB: dbm.NewMemDB()}}
}

// Implements CommitStore
// Commit cleans up Store.
func (ts *Store) Commit(context.Context, *iavl.TreeDelta) (ctx context.Context, id types.CommitID, _ iavl.TreeDelta) {
	ts.Store = dbadapter.Store{DB: dbm.NewMemDB()}
	return
}

// Implements CommitStore
func (ts *Store) SetPruning(pruning types.PruningOptions) {
}

// Implements CommitStore
func (ts *Store) LastCommitID() (id types.CommitID) {
	return
}

// Implements Store.
func (ts *Store) GetStoreType() types.StoreType {
	return types.StoreTypeTransient
}
