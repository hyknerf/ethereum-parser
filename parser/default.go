package parser

import (
	"github.com/hyknerf/ethereum-parser/store"
	"github.com/hyknerf/ethereum-parser/types"
	"log"
	"net/http"
	"time"
)

type DefaultParser struct {
	httpClient *http.Client
	store      store.Storer
}

func NewDefaultParser(store store.Storer) *DefaultParser {
	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: 30 * time.Second,
	}
	return &DefaultParser{
		httpClient: httpClient,
		store:      store,
	}
}

func (dp *DefaultParser) GetCurrentBlock() int {
	return int(dp.store.GetLastBlockNumber())
}

func (dp *DefaultParser) Subscribe(address string) bool {
	if err := dp.store.AddObservedAddress(address); err != nil {
		return false
	}

	log.Println("observed address:", address)
	return true
}

func (dp *DefaultParser) GetTransactions(address string) []*types.TransactionReceipt {
	res := dp.store.GetTransactions(address)
	return res
}
