package store

import (
	"fmt"
	"github.com/hyknerf/ethereum-parser/types"
	"log"
	"strings"
	"sync"
)

type MemStore struct {
	lastBlock    int64
	addresses    map[string]struct{}
	transactions map[string][]*types.TransactionReceipt
	lock         sync.RWMutex
}

func NewMemStore() *MemStore {
	return &MemStore{
		lastBlock:    0,
		addresses:    make(map[string]struct{}),
		transactions: make(map[string][]*types.TransactionReceipt),
		lock:         sync.RWMutex{},
	}
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

	address = strings.ToLower(address)
	if _, exists := mem.addresses[address]; exists {
		log.Printf("address %s already exists in mem store", address)
		return fmt.Errorf("address %s already exists", address)
	}

	log.Printf("adding address %s to mem store", address)

	mem.addresses[address] = struct{}{}
	log.Println(mem.addresses)
	return nil
}

func (mem *MemStore) IsObservedAddress(address string) bool {
	mem.lock.RLock()
	defer mem.lock.RUnlock()

	address = strings.ToLower(address)
	_, observed := mem.addresses[address]

	return observed
}

func (mem *MemStore) AddTransaction(transaction *types.TransactionReceipt) {
	mem.lock.Lock()
	defer mem.lock.Unlock()

	if transaction.Hash == "" {
		log.Printf("transaction receipt hash is empty")
		return
	}

	addressTo := strings.ToLower(transaction.To)
	addressFrom := strings.ToLower(transaction.From)
	mem.transactions[addressTo] = append(mem.transactions[addressTo], transaction)
	mem.transactions[addressFrom] = append(mem.transactions[addressFrom], transaction)
	log.Printf("added new transaction %s to address [to:%s, from:%s]", transaction.Hash, addressTo, addressFrom)
}

func (mem *MemStore) GetTransactions(address string) []*types.TransactionReceipt {
	mem.lock.RLock()
	defer mem.lock.RUnlock()

	address = strings.ToLower(address)
	log.Printf("getting transactions for address %s", address)
	log.Println(mem.transactions[address])

	return mem.transactions[address]
}
