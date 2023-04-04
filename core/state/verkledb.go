package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gballet/go-verkle"
	"sync"
)

// VerkleDB implements state.Database for a verkle tree
type VerkleDB struct {
	cachingDB
	addrToPoint PointCache
}

// OpenTrie opens the main account trie.
func (db *VerkleDB) OpenTrie(root common.Hash) (Trie, error) {
	return nil, nil
}

type PointCache struct {
	cache map[string]*verkle.Point
	lock  sync.RWMutex
}

func NewPointCache() *PointCache {
	return &PointCache{
		cache: make(map[string]*verkle.Point),
	}
}
