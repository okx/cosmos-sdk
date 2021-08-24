package rootmulti

import (
	"github.com/cosmos/cosmos-sdk/store/dbadapter"
	"github.com/cosmos/cosmos-sdk/store/types"
	"github.com/tendermint/iavl"
)

var commithash = []byte("FAKE_HASH")

//----------------------------------------
// commitDBStoreWrapper should only be used for simulation/debugging,
// as it doesn't compute any commit hash, and it cannot load older state.

// Wrapper type for dbm.Db with implementation of KVStore
type commitDBStoreAdapter struct {
	dbadapter.Store
}

func (cdsa commitDBStoreAdapter) Commit(*iavl.TreeDelta, []byte) (types.CommitID, iavl.TreeDelta) {
	return types.CommitID{
		Version: -1,
		Hash:    commithash,
	}, iavl.TreeDelta{}
}

func (cdsa commitDBStoreAdapter) LastCommitID() types.CommitID {
	return types.CommitID{
		Version: -1,
		Hash:    commithash,
	}
}

func (cdsa commitDBStoreAdapter) SetPruning(_ types.PruningOptions) {}
