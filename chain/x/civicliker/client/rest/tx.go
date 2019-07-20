package rest

import (
	"bytes"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/likecoin/hackatom-seoul-2019/chain/x/civicliker/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(
		"/civicliker/like",
		postLikeHandlerFn(cdc, cliCtx),
	).Methods("POST")
}

type (
	LikeRequest struct {
		BaseReq rest.BaseReq   `json:"base_req"`
		Liker   sdk.AccAddress `json:"liker"` // in bech32
		Likee   sdk.AccAddress `json:"likee"` // in bech32
		Url     string         `json:"url"`
		Count   uint64         `json:"count"`
	}
)

func postLikeHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LikeRequest

		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.MsgLike{
			Liker: req.Liker,
			Likee: req.Likee,
			Url:   req.Url,
			Count: req.Count,
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

		if !bytes.Equal(fromAddr, req.Liker) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "liker must be transaction sender")
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
