package parser

import "github.com/hyknerf/ethereum-parser/types"

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []*types.TransactionReceipt
}
