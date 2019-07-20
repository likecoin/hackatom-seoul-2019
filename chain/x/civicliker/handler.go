package civicliker

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/hackatom-seoul-2019/chain/x/civicliker/types"

	subTypes "github.com/likecoin/hackatom-seoul-2019/chain/x/subscription/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		switch msg := msg.(type) {
		case types.MsgLike:
			return handleMsgLike(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("invalid message parse in subscription module").Result()
		}
	}
}

func handleMsgLike(ctx sdk.Context, msg types.MsgLike, keeper Keeper) sdk.Result {
	s := keeper.GetLikerState(ctx, msg.Liker)
	if !msg.Liker.Equals(s.Liker) {
		return sdk.ErrUnauthorized("Unsubscribed user cannot like").Result()
	}
	r := types.LikeRecord{
		Liker: msg.Liker,
		Likee: msg.Likee,
		Url:   msg.Url,
		Count: msg.Count,
		Time:  ctx.BlockHeader().Time,
	}
	likerStr := msg.Liker.String()
	cache := keeper.GetBlockLikeRecordsCache()
	cache[likerStr] = append(cache[likerStr], r)
	return sdk.Result{}
}

func EndBlocker(ctx sdk.Context, keeper Keeper) {
	cache := keeper.GetBlockLikeRecordsCache()
	likeGlobalCount := keeper.GetLikeGlobalCount(ctx)
	adjFactor := float64(0.9) // TODO: no hard coding
	for _, records := range cache {
		liker := records[0].Liker
		s := keeper.GetLikerState(ctx, liker)
		sub := keeper.subscriptionKeeper.GetSubscription(ctx, liker, s.SubscriptionChannelID)
		ch := keeper.subscriptionKeeper.GetChannel(ctx, sub.ChannelID)
		likerBlockLikeSum := uint64(0)
		for _, r := range records {
			likerBlockLikeSum += r.Count
		}
		s.CurrPeriodLikeCount += likerBlockLikeSum
		remainingBlocks := s.PeriodEndingBlock - ctx.BlockHeight() - 1
		coinsAmountPerLike := calCoinsAmountPerLike(
			s.PrevPeriodLikeCount, s.CurrPeriodLikeCount,
			sub, ch, remainingBlocks, likerBlockLikeSum, adjFactor,
		)
		for _, r := range records {
			coinsAmount := float64(r.Count) * coinsAmountPerLike
			coin := sdk.NewCoin(ch.Price.Denom, sdk.NewInt(int64(coinsAmount)))
			sub.Remaining = sub.Remaining.Sub(coin)
			r.CoinDistributed = coin
			keeper.coinKeeper.AddCoins(ctx, r.Likee, sdk.NewCoins(coin))
			likeGlobalCount++
			keeper.SetLikeGlobalCount(ctx, likeGlobalCount)
			keeper.SetLikeRecord(ctx, r)
			count := keeper.GetLikeRecordCount(ctx, liker, r.Url)
			count += r.Count
			keeper.SetLikeRecordCount(ctx, liker, r.Url, count)
		}
		keeper.subscriptionKeeper.SetSubscription(ctx, sub)
		keeper.SetLikerState(ctx, s)
	}
	for likerStr := range cache {
		delete(cache, likerStr)
	}
}

func GetPaymentHook(keeper Keeper) func(sdk.Context, *subTypes.Subscription, subTypes.Channel) {
	return func(ctx sdk.Context, sub *subTypes.Subscription, ch subTypes.Channel) {
		// TODO: select channel, or only hook specific channel
		s := keeper.GetLikerState(ctx, sub.Subscriber)
		s.Liker = sub.Subscriber
		s.PrevPeriodLikeCount = s.CurrPeriodLikeCount
		s.CurrPeriodLikeCount = 0
		s.SubscriptionChannelID = ch.ID
		s.PeriodEndingBlock = sub.NextPaymentBlock
		keeper.SetLikerState(ctx, s)
	}
}
