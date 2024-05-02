package observer

import (
	"context"
	"fmt"
	"sync"
)

type TransactionObserver struct {
	ctx            context.Context
	newBlockStream <-chan string
}

func NewTransactionObserver(ctx context.Context, newBlockStream <-chan string) *TransactionObserver {
	return &TransactionObserver{
		ctx:            ctx,
		newBlockStream: newBlockStream,
	}
}

func (t *TransactionObserver) Work(wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-t.ctx.Done():
				return
			case block := <-t.newBlockStream:
				fmt.Printf("New Block: %s\n", block)
			}
		}
	}()
}
