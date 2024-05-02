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

type TransactionObserver struct {
	ctx               context.Context
	transactionStream chan string
	store             store.Storer
	httpClient        *http.Client
}

func NewTransactionObserver(ctx context.Context, transactionStream chan string, store store.Storer, httpClient *http.Client) *TransactionObserver {
	return &TransactionObserver{
		ctx:               ctx,
		transactionStream: transactionStream,
		store:             store,
		httpClient:        httpClient,
	}
}

func (o *TransactionObserver) Work(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-o.ctx.Done():
				return
			case txHash := <-o.transactionStream:
				log.Println("fetching txs receipt:", txHash)
				baseReq := types.BaseRequest{
					JsonRPC: types.JsonRPCVersion,
					Method:  types.MethodEthGetTransactionReceipt,
					Params: []interface{}{
						txHash,
					},
				}

				payloadByte, err := json.Marshal(baseReq)
				if err != nil {
					fmt.Println(err)
				}

				req, err := http.NewRequest(http.MethodPost, types.Host, bytes.NewBuffer(payloadByte))
				res, err := o.httpClient.Do(req)
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

				txRcpt := new(types.TransactionReceipt)
				err = json.Unmarshal(baseRes.Result, txRcpt)
				if err != nil {
					log.Println(err)
				}

				// TODO: Receipt is not available until txs is final, need to re-queue and check again
				o.store.AddTransaction(txRcpt)
			}
		}
	}()
}
