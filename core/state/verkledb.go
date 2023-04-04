package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/trie/utils"
	"github.com/gballet/go-verkle"
)

// VerkleDB implements state.Database for a verkle tree
type VerkleDB struct {
	cachingDB
	addrToPoint *utils.PointCache
}

// OpenTrie opens the main account trie.
func (db *VerkleDB) OpenTrie(root common.Hash) (Trie, error) {
	if root == (common.Hash{}) || root == types.EmptyRootHash {
		return trie.NewVerkleTrie(verkle.New(), db.cachingDB.triedb, db.addrToPoint), nil
	}
	payload, err := db.DiskDB().Get(root[:])
	if err != nil {
		return nil, err
	}

	r, err := verkle.ParseNode(payload, 0, root[:])
	if err != nil {
		panic(err)
	}
	return trie.NewVerkleTrie(r, db.cachingDB.triedb, db.addrToPoint), err
}

// OpenStorageTrie opens the storage trie of an account.
func (db *VerkleDB) OpenStorageTrie(stateRoot, addrHash, root common.Hash) (Trie, error) {
	// alternatively, return accTrie
	panic("should not be called")
}

// CopyTrie returns an independent copy of the given trie.
func (db *VerkleDB) CopyTrie(tr Trie) Trie {
	t, ok := tr.(*trie.VerkleTrie)
	if ok {
		return t.Copy(db.cachingDB.triedb)
	}

	panic("invalid tree type != VerkleTrie")
}

// ContractCode retrieves a particular contract's code.
//func (db *VerkleDB) ContractCode(addrHash, codeHash common.Hash) ([]byte, error) {
//	return db.cachingDB.ContractCode(addrHash, codeHash)
//}

// ContractCodeSize retrieves a particular contracts code's size.
// func (db *VerkleDB) ContractCodeSize(addrHash, codeHash common.Hash) (int, error) {
//	return db.cachingDB.ContractCodeSize(addrHash, codeHash)
//}

// DiskDB retrieves the low level trie database used for data storage.
// func (db *VerkleDB) DiskDB() ethdb.KeyValueStore {
//	return db.cachingDB.disk
//}

// TrieDB retrieves the low level trie database used for data storage.
// func (db *VerkleDB) TrieDB() *trie.Database {
//	return db.cachingDB.triedb
//}
