package rest

import (
	"bytes"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/likecoin/hackatom-seoul-2019/chain/x/subscription/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, kb keys.Keybase) {
	r.HandleFunc(
		"/contentdb/upload",
		postUploadHandlerFn(cdc, kb, cliCtx),
	).Methods("POST")
}

type (
	UploadContentRequest struct {
		BaseReq rest.BaseReq   `json:"base_req"`
		Author  sdk.AccAddress `json:"author"` // in bech32
		Url     string         `json:"url"`
	}
)

func postUploadHandlerFn(cdc *codec.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UploadContentRequest

		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.MsgUploadContent{
			Author: req.Author,
			Url:    req.Url,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, req.Author) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "author must be sender")
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
