package router

import (
	"encoding/json"
	"fmt"
	"github.com/hyknerf/ethereum-parser/parser"
	"github.com/hyknerf/ethereum-parser/types"
	"net/http"
)

type Router struct {
	parser parser.Parser
}

func NewRouter(parser parser.Parser) *Router {
	return &Router{
		parser: parser,
	}
}

func (r *Router) BlockHandler(res http.ResponseWriter, req *http.Request) {
	_, _ = fmt.Fprintf(res, "block: %d", r.parser.GetCurrentBlock())
}

func (r *Router) SubscribeAddressHandler(res http.ResponseWriter, req *http.Request) {
	reqAddress := new(types.SubscribeAddressRequest)

	err := json.NewDecoder(req.Body).Decode(&reqAddress)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if !r.parser.Subscribe(reqAddress.Address) {
		http.Error(res, fmt.Sprintf("failed to subscribe to address %s", reqAddress.Address), http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusAccepted)
	_, _ = fmt.Fprintf(res, "address: %s", reqAddress.Address)
}

func (r *Router) TransactionHandler(res http.ResponseWriter, req *http.Request) {
	addr := req.URL.Query().Get("address")
	txs := r.parser.GetTransactions(addr)
	res.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(res).Encode(txs)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
