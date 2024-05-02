package http

import (
	"fmt"
	"github.com/hyknerf/ethereum-parser/store"
	"github.com/hyknerf/ethereum-parser/types"
	"net/http"
)

type Router struct {
	store  store.Storer
	parser types.Parser
}

func NewRouter(store store.Storer, parser types.Parser) *Router {
	return &Router{
		store:  store,
		parser: parser,
	}
}

func (r *Router) BlockHandler(w http.ResponseWriter, req *http.Request) {
	_, _ = fmt.Fprintf(w, "block: %d", r.store.GetLastBlockNumber())
}
