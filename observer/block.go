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
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"
)

type BlockNumberObserver struct {
	ctx            context.Context
	newBlockStream chan<- int64
	store          store.Storer
	httpClient     *http.Client
	interval       time.Duration
}

func NewBlockNumberObserver(ctx context.Context, newBlockStream chan<- int64, store store.Storer, httpClient *http.Client, interval time.Duration) *BlockNumberObserver {
	return &BlockNumberObserver{
		ctx:            ctx,
		newBlockStream: newBlockStream,
		store:          store,
		httpClient:     httpClient,
		interval:       interval,
	}
}

func (o *BlockNumberObserver) Work(wg *sync.WaitGroup) {
	ticker := time.NewTicker(o.interval)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-o.ctx.Done():
				return
			case <-ticker.C:
				log.Println("fetching latest block number")
				baseReq := types.BaseRequest{
					JsonRPC: types.JsonRPCVersion,
					Method:  types.MethodEthBlockNumber,
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
					log.Fatal(baseRes.Error)
				}

				n := new(big.Int)
				n.SetString(strings.Trim(string(baseRes.Result), "\""), 0)

				if lastBLock := o.store.GetLastBlockNumber(); lastBLock > n.Int64() || lastBLock != n.Int64() {
					o.newBlockStream <- n.Int64()
					o.store.SaveLastBlockNumber(n.Int64())
				}
			}
		}
	}()
}
