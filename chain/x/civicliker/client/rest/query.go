package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/gorilla/mux"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(
		"/civicliker/like-count/{liker}/{url}",
		likeCountHandlerFn(cliCtx, cdc),
	).Methods("GET")
	r.HandleFunc(
		"/civicliker/like-history/{liker}",
		likeHistoryHandlerFn(cliCtx, cdc),
	).Methods("GET")
}

func likeCountHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryLikeCount(cliCtx, cdc, "custom/civicliker/like-count")
}

func likeHistoryHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryLikeHistory(cliCtx, cdc, "custom/civicliker/like-history")
}
