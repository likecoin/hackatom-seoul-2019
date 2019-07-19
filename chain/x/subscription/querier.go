package contentdb

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryContent = "content"
)

func NewQuerier(k Keeper, cdc *codec.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryContent:
			return queryContent(ctx, cdc, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown staking query endpoint")
		}
	}
}

func queryContent(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params QueryContentParams

	err := cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	content := k.GetContent(ctx, params.Url)
	if content.Author.Empty() {
		return []byte("null"), nil
	}

	res, err := codec.MarshalJSONIndent(cdc, content)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

type QueryContentParams struct {
	Url string
}
