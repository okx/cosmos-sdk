package trie

import (
	"github.com/cosmos/cosmos-sdk/store/cachekv"
	"github.com/cosmos/cosmos-sdk/store/tracekv"
	"github.com/cosmos/cosmos-sdk/store/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/trie"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
	"io"
)

const (
	defaultIAVLCacheSize = 10000
)

var (
	_         types.KVStore       = (*Store)(nil)
	_         types.CommitStore   = (*Store)(nil)
	_         types.CommitKVStore = (*Store)(nil)
	_         types.Queryable     = (*Store)(nil)
)

type Store struct {
	trie *trie.Trie
	id types.CommitID
}

func init() {

}

// LoadStore returns an IAVL Store as a CommitKVStore. Internally, it will load the
// store's version (id) from the provided DB. An error is returned if the version
// fails to load.
func LoadStore(db dbm.DB, id types.CommitID, lazyLoading bool, startVersion int64) (types.CommitKVStore, error) {
	return LoadStoreWithInitialVersion(db, id, lazyLoading, uint64(startVersion))
}

// LoadStore returns an IAVL Store as a CommitKVStore setting its initialVersion
// to the one given. Internally, it will load the store's version (id) from the
// provided DB. An error is returned if the version fails to load.
func LoadStoreWithInitialVersion(db dbm.DB, id types.CommitID, lazyLoading bool, initialVersion uint64) (types.CommitKVStore, error) {
	ethdb := &EthDbAdapter{db}

	trie, err := trie.New(ethcmn.BytesToHash(id.Hash), ethdb)
	if err != nil {
		return nil, err
	}
	store := &Store{
			trie,
			id,
	}
	return store, nil
}

// GetImmutable returns a reference to a new store backed by an immutable IAVL
// tree at a specific version (height) without any pruning options. This should
// be used for querying and iteration only. If the version does not exist or has
// been pruned, an empty immutable IAVL tree will be used.
// Any mutable operations executed will result in a panic.
func (st *Store) GetImmutable(version int64) (*Store, error) {
	return nil, nil
}

// Commit commits the current store state and returns a CommitID with the new
// version and hash.
func (st *Store) Commit() types.CommitID {
	hash, err := st.trie.Commit(nil)
	if err != nil {
		panic(err)
	}
	st.id.Version ++
	st.id.Hash = hash.Bytes()
	return types.CommitID{
		Version: st.id.Version,
		Hash:    st.id.Hash,
	}
}

// Implements Committer.
func (st *Store) LastCommitID() types.CommitID {
	return types.CommitID{
		Version: st.id.Version,
		Hash:    st.id.Hash,
	}
}

// SetPruning panics as pruning options should be provided at initialization
// since IAVl accepts pruning options directly.
func (st *Store) SetPruning(_ types.PruningOptions) {
	panic("cannot set pruning options on an initialized IAVL store")
}

// VersionExists returns whether or not a given version is stored.
func (st *Store) VersionExists(version int64) bool {
	return false
}

// Implements Store.
func (st *Store) GetStoreType() types.StoreType {
	return types.StoreTypeTrie
}

// Implements Store.
func (st *Store) CacheWrap() types.CacheWrap {
	return cachekv.NewStore(st)
}

// CacheWrapWithTrace implements the Store interface.
func (st *Store) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	return cachekv.NewStore(tracekv.NewStore(st, w, tc))
}

// Implements types.KVStore.
func (st *Store) Set(key, value []byte) {
	types.AssertValidValue(value)
	st.trie.TryUpdate(key, value)
}

// Implements types.KVStore.
func (st *Store) Get(key []byte) []byte {
	value, err := st.trie.TryGet(key)
	if err != nil {
		panic(err)
	}
	return value
}

// Implements types.KVStore.
func (st *Store) Has(key []byte) (exists bool) {
	_, err := st.trie.TryGet(key)
	if err != nil {
		return false
	}
	return true
}

// Implements types.KVStore.
func (st *Store) Delete(key []byte) {
	err := st.trie.TryDelete(key)
	if err != nil {
		panic(err)
	}
}

// DeleteVersions deletes a series of versions from the MutableTree. An error
// is returned if any single version is invalid or the delete fails. All writes
// happen in a single batch with a single commit.
func (st *Store) DeleteVersions(versions ...int64) error {
	return nil
}

// Implements types.KVStore.
func (st *Store) Iterator(start, end []byte) types.Iterator {
	return nil
}

// Implements types.KVStore.
func (st *Store) ReverseIterator(start, end []byte) types.Iterator {
	return nil
}

// Query implements ABCI interface, allows queries
//
// by default we will return from (latest height -1),
// as we will have merkle proofs immediately (header height = data height + 1)
// If latest-1 is not present, use latest (which must be present)
// if you care to have the latest data to see a tx results, you must
// explicitly set the height you want to see
func (st *Store) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	return abci.ResponseQuery{}
}

