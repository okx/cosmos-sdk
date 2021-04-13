package trie

import (
	"github.com/ethereum/go-ethereum/ethdb"
	dbm "github.com/tendermint/tm-db"
)

type EthDbAdapter struct {
	dbmDB dbm.DB
}


func (dbAdapter *EthDbAdapter) Has(key []byte) (bool, error) {
	return dbAdapter.dbmDB.Has(key)
}

// Get retrieves the given key if it's present in the key-value data store.
func (dbAdapter *EthDbAdapter) Get(key []byte) ([]byte, error) {
	return dbAdapter.dbmDB.Get(key)
}
func (dbAdapter *EthDbAdapter) HasAncient(kind string, number uint64) (bool, error) {
	return false, nil
}

// Ancient retrieves an ancient binary blob from the append-only immutable files.
func (dbAdapter *EthDbAdapter) Ancient(kind string, number uint64) ([]byte, error) {
	return nil, nil
}

// Ancients returns the ancient item numbers in the ancient store.
func (dbAdapter *EthDbAdapter) Ancients() (uint64, error) {
	return 0, nil
}

// AncientSize returns the ancient size of the specified category.
func (dbAdapter *EthDbAdapter) AncientSize(kind string) (uint64, error) {
	return 0, nil
}

// Put inserts the given value into the key-value data store.
func (dbAdapter *EthDbAdapter) Put(key []byte, value []byte) error {
	return dbAdapter.dbmDB.Set(key, value)
}

// Delete removes the key from the key-value data store.
func (dbAdapter *EthDbAdapter) Delete(key []byte) error {
	return dbAdapter.dbmDB.Delete(key)
}

func (dbAdapter *EthDbAdapter) AppendAncient(number uint64, hash, header, body, receipt, td []byte) error {
	return nil
}

// TruncateAncients discards all but the first n ancient data from the ancient store.
func (dbAdapter *EthDbAdapter) TruncateAncients(n uint64) error {
	return nil
}

// Sync flushes all in-memory ancient store data to disk.
func (dbAdapter *EthDbAdapter) Sync() error {
	return nil
}

func (dbAdapter *EthDbAdapter) NewBatch() ethdb.Batch {
	return nil
}

func (dbAdapter *EthDbAdapter) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	return nil
}

func (dbAdapter *EthDbAdapter) Stat(property string) (string, error) {
	return "", nil
}

func (dbAdapter *EthDbAdapter) Compact(start []byte, limit []byte) error {
	return  nil
}

func (dbAdapter *EthDbAdapter) Close() error {
	return dbAdapter.dbmDB.Close()
}