package store

import "github.com/hyknerf/ethereum-parser/types"

type Storer interface {
	SaveLastBlockNumber(block int64)
	GetLastBlockNumber() int64
	AddObservedAddress(address string) error
	IsObservedAddress(address string) bool
	AddTransaction(transaction *types.TransactionReceipt)
	GetTransactions(address string) []*types.TransactionReceipt
}
