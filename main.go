package main

import (
	"context"
	"github.com/hyknerf/ethereum-parser/observer"
	"github.com/hyknerf/ethereum-parser/parser"
	router2 "github.com/hyknerf/ethereum-parser/router"
	"github.com/hyknerf/ethereum-parser/store"
	"net/http"
	"sync"
	"time"
)

const (
	blockObserverInterval = 2 * time.Second
)

func main() {
	ctx := context.Background()
	newBlockStream := make(chan int64)
	txHashStream := make(chan string)

	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: 30 * time.Second,
	}

	memStore := store.NewMemStore()

	// register Uniswap address so it is not empty
	_ = memStore.AddObservedAddress("0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD")

	defaultParser := parser.NewDefaultParser(memStore)
	blockNumObs := observer.NewBlockNumberObserver(ctx, newBlockStream, memStore, httpClient, blockObserverInterval)
	blockObs := observer.NewNewBlockObserver(ctx, newBlockStream, txHashStream, memStore, httpClient)
	txObs := observer.NewTransactionObserver(ctx, txHashStream, memStore, httpClient)

	wg := &sync.WaitGroup{}
	blockNumObs.Work(wg)
	blockObs.Work(wg)
	txObs.Work(wg)

	router := router2.NewRouter(defaultParser)

	http.HandleFunc("GET /block", router.BlockHandler)
	http.HandleFunc("POST /subscribe-address", router.SubscribeAddressHandler)
	http.HandleFunc("GET /transactions", router.TransactionHandler)
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}

	wg.Wait()
}
