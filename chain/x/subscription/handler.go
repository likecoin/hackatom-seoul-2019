package subscription

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/hackatom-seoul-2019/chain/x/subscription/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		switch msg := msg.(type) {
		case types.MsgSubscribe:
			return handleMsgSubscribe(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("invalid message parse in subscription module").Result()
		}
	}
}

func handleMsgSubscribe(ctx sdk.Context, msg types.MsgSubscribe, keeper Keeper) sdk.Result {
	ch := keeper.GetChannel(ctx, msg.ChannelID)
	if ch.ID != msg.ChannelID {
		return sdk.ErrUnauthorized("Invalid channel ID").Result()
	}
	if keeper.HasSubscription(ctx, msg.Subscriber, msg.ChannelID) {
		return sdk.ErrUnauthorized("Already subscribed").Result()
	}
	_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Subscriber, sdk.NewCoins(ch.Price))
	if err != nil {
		return sdk.ErrInsufficientCoins("Not enough balance for subscription payment").Result()
	}
	sub := types.Subscription{
		Subscriber:       msg.Subscriber,
		ChannelID:        msg.ChannelID,
		Remaining:        ch.Price,
		NextPaymentBlock: ctx.BlockHeight() + ch.PeriodBlocks,
	}
	for _, hook := range keeper.paymentHooks {
		hook(ctx, &sub, ch)
	}
	keeper.SetSubscription(ctx, sub)
	tags := sdk.NewTags(
		"subscriber", msg.Subscriber.String(),
		"subscription_payment_channel", fmt.Sprintf("%d", ch.ID),
		"subscription_payment_value", ch.Price.String(),
	)
	return sdk.Result{
		Tags: tags,
	}
}

func BeginBlocker(ctx sdk.Context, keeper Keeper) {
	height := ctx.BlockHeight()
	keeper.IterateSubscriptions(ctx, func(index int64, sub types.Subscription) (stop bool) {
		if sub.NextPaymentBlock == height {
			ch := keeper.GetChannel(ctx, sub.ChannelID)
			_, _, err := keeper.coinKeeper.SubtractCoins(ctx, sub.Subscriber, sdk.NewCoins(ch.Price))
			if err == nil {
				sub.Remaining = sub.Remaining.Add(ch.Price)
			}
			sub.NextPaymentBlock = height + ch.PeriodBlocks
			for _, hook := range keeper.paymentHooks {
				hook(ctx, &sub, ch)
			}
			keeper.SetSubscription(ctx, sub)
		}
		return false
	})
}
