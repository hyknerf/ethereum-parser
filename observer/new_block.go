package observer

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hyknerf/ethereum-parser/store"
	"github.com/hyknerf/ethereum-parser/types"
	"log"
	"net/http"
	"sync"
)

type NewBlockObserver struct {
	ctx               context.Context
	newBlockStream    <-chan int64
	transactionStream chan<- string
	store             store.Storer
	httpClient        *http.Client
}

func NewNewBlockObserver(ctx context.Context, newBlockStream <-chan int64, transactionStream chan<- string, store store.Storer, httpClient *http.Client) *NewBlockObserver {
	return &NewBlockObserver{
		ctx:               ctx,
		newBlockStream:    newBlockStream,
		transactionStream: transactionStream,
		store:             store,
		httpClient:        httpClient,
	}
}

func (t *NewBlockObserver) Work(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-t.ctx.Done():
				return
			case block := <-t.newBlockStream:
				log.Println("fetching txs on new block:", fmt.Sprintf("0x%x", block))
				baseReq := types.BaseRequest{
					JsonRPC: types.JsonRPCVersion,
					Method:  types.MethodEthGetBlockByNumber,
					Params: []interface{}{
						fmt.Sprintf("0x%x", block),
						true,
					},
				}

				payloadByte, err := json.Marshal(baseReq)
				if err != nil {
					fmt.Println(err)
				}

				req, err := http.NewRequest(http.MethodPost, types.Host, bytes.NewBuffer(payloadByte))
				res, err := t.httpClient.Do(req)
				if err != nil {
					fmt.Println(err)
				}

				baseRes := &types.BaseResponse{}
				err = json.NewDecoder(bufio.NewReader(res.Body)).Decode(&baseRes)
				if err != nil {
					fmt.Println(err)
				}
				_ = res.Body.Close()

				if baseRes.Error != nil {
					log.Fatal(baseRes.Error.Message)
				}

				blockTxs := new(types.BlockByNumberResponse)
				err = json.Unmarshal(baseRes.Result, blockTxs)
				if err != nil {
					log.Println(err)
				}

				for _, tx := range blockTxs.Transactions {
					if t.store.IsObservedAddress(tx.To) || t.store.IsObservedAddress(tx.From) {
						log.Printf("observed address, adding txs")
						t.transactionStream <- tx.Hash
					}
				}
			}
		}
	}()
}
