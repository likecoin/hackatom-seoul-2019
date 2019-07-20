package civicliker

import (
	"fmt"
	"sort"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/likecoin/hackatom-seoul-2019/chain/x/civicliker/types"
)

const (
	QueryLikeCount   = "like-count"
	QueryLikeHistory = "like-history"
)

func NewQuerier(k Keeper, cdc *codec.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryLikeCount:
			return queryLikeCount(ctx, cdc, req, k)
		case QueryLikeHistory:
			return queryLikeHistory(ctx, cdc, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown civicliker query endpoint")
		}
	}
}

func queryLikeCount(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params QueryLikeCountParams

	err := cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	count := k.GetLikeRecordCount(ctx, params.Liker, params.Url)
	res := []byte(fmt.Sprintf("%d", count))
	return res, nil
}

type QueryLikeCountParams struct {
	Liker sdk.AccAddress `json:"liker"`
	Url   string         `json:"url"`
}

type LikeHistory []types.LikeRecord

func (h LikeHistory) Len() int {
	return len(h)
}

func (h LikeHistory) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h LikeHistory) Less(i, j int) bool {
	return h[i].Time.Before(h[j].Time)
}

func queryLikeHistory(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params QueryLikeHistoryParams

	err := cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	history := LikeHistory{}
	k.IterateLikerLikeHistory(ctx, params.Liker, func(index int64, r types.LikeRecord) (stop bool) {
		history = append(history, r)
		return false
	})
	sort.Sort(history)
	res, err := codec.MarshalJSONIndent(cdc, history)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}

type QueryLikeHistoryParams struct {
	Liker sdk.AccAddress `json:"liker"`
	Url   string         `json:"url"`
}
