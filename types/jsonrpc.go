package types

import "encoding/json"

const (
	Host           = "https://cloudflare-eth.com"
	JsonRPCVersion = "2.0"

	MethodEthBlockNumber           = "eth_blockNumber"
	MethodEthGetBlockByNumber      = "eth_getBlockByNumber"
	MethodEthGetTransactionReceipt = "eth_getTransactionReceipt"
)

type BaseRequest struct {
	JsonRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

//func (br *BaseRequest) MarshalJSON() ([]byte, error) {
//	arr := []interface{}{br.JsonRPC, br.Method, br.ID}
//	return json.Marshal(arr)
//}

type BaseError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

type BaseResponse struct {
	ID      uint8           `json:"id"`
	JsonRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *BaseError      `json:"error"`
}

type BlockByNumberResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type TransactionReceipt struct {
	Hash    string `json:"transactionHash"`
	From    string `json:"from"`
	To      string `json:"to"`
	GasUsed string `json:"gasUsed"`
}

type Transaction struct {
	Hash string `json:"hash"`
	From string `json:"from"`
	To   string `json:"to"`
}

type ResultBlockByNumber struct {
	Transactions []Transaction `json:"transactions"`
}
