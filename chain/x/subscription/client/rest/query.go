package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/gorilla/mux"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// Get all delegations from a delegator
	r.HandleFunc(
		"/contentdb/content/{url}",
		contentHandlerFn(cliCtx, cdc),
	).Methods("GET")
}

// HTTP request handler to query a delegator delegations
func contentHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryContent(cliCtx, cdc, "custom/contentdb/content")
}
