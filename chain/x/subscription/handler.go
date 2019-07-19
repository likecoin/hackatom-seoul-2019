package contentdb

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/likecoin/hackatom-seoul-2019/chain/x/subscription/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		switch msg := msg.(type) {
		case types.MsgUploadContent:
			return handleMsgUploadContent(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("invalid message parse in contentdb module").Result()
		}
	}
}

func handleMsgUploadContent(ctx sdk.Context, msg types.MsgUploadContent, keeper Keeper) sdk.Result {
	if keeper.HasContent(ctx, msg.Url) {
		return sdk.ErrUnauthorized("Content is already set by others").Result()
	}
	keeper.SetContent(ctx, types.Content{
		Author: msg.Author,
		Url:    msg.Url,
	})
	tags := sdk.NewTags(
		"author", msg.Author.String(),
		"url", msg.Url,
	)
	return sdk.Result{
		Tags: tags,
	}
}
