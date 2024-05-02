package store

import (
	"fmt"
	"log"
	"sync"
)

type MemStore struct {
	lastBlock int64
	addresses map[string]struct{}
	lock      sync.RWMutex
}

func NewMemStore() *MemStore {
	return &MemStore{}
}

func (mem *MemStore) SaveLastBlockNumber(block int64) {
	mem.lock.Lock()
	defer mem.lock.Unlock()

	log.Printf("saving last block number to mem store: %d", block)

	mem.lastBlock = block
}

func (mem *MemStore) GetLastBlockNumber() int64 {
	mem.lock.RLock()
	defer mem.lock.RUnlock()

	log.Printf("getting last block number from mem store: %d", mem.lastBlock)

	return mem.lastBlock
}

func (mem *MemStore) AddObservedAddress(address string) error {
	mem.lock.Lock()
	defer mem.lock.Unlock()

	if _, exists := mem.addresses[address]; exists {
		log.Printf("address %s already exists in mem store", address)
		return fmt.Errorf("address %s already exists", address)
	}

	log.Printf("adding address %s to mem store", address)

	mem.addresses[address] = struct{}{}
	return nil
}
