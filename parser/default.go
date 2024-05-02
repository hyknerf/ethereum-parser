package types

import (
	"github.com/hyknerf/ethereum-parser/store"
	"log"
	"net/http"
	"time"
)

const (
	Host           = "https://cloudflare-eth.com"
	JsonRPCVersion = "2.0"

	MethodEthBlockNumber = "eth_blockNumber"
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
		log.Fatal("failed to add address to observed addresses", err)
		return false
	}
	return true
}

func (dp *DefaultParser) GetTransactions(address string) []Transaction {
	return nil
}
