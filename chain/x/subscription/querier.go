package subscription

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QuerySubscription = "subscription"
)

func NewQuerier(k Keeper, cdc *codec.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QuerySubscription:
			return querySubscription(ctx, cdc, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown subscription query endpoint")
		}
	}
}

func querySubscription(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params QuerySubscriptionParams

	err := cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	sub := k.GetSubscription(ctx, params.Subscriber, params.ChannelID)
	if sub.ChannelID != params.ChannelID {
		return []byte("null"), nil
	}

	res, err := codec.MarshalJSONIndent(cdc, sub)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

type QuerySubscriptionParams struct {
	Subscriber sdk.AccAddress `json:"subscriber"`
	ChannelID  uint64         `json:"channel_id"`
}
