package civicliker

import (
	"crypto/sha256"

	subTypes "github.com/likecoin/hackatom-seoul-2019/chain/x/subscription/types"
)

func hashUrl(url string) []byte {
	bz := sha256.Sum256([]byte(url))
	return bz[:]
}

func interpolate(a, b, t float64) float64 {
	return a*t + b*(1-t)
}

func estimateLikesInPeriod(
	prevPeriodLikeCount, currPeriodLikeCount uint64,
	remainingBlocks, periodLength int64) float64 {
	passedBlocks := periodLength - remainingBlocks
	e := float64(currPeriodLikeCount) * float64(periodLength) / float64(passedBlocks)
	if prevPeriodLikeCount != 0 {
		e = interpolate(float64(prevPeriodLikeCount), e, float64(passedBlocks)/float64(periodLength))
	}
	return e
}

func calCoinsAmountPerLike(
	prevPeriodLikeCount, currPeriodLikeCount uint64,
	sub subTypes.Subscription, ch subTypes.Channel,
	remainingBlocks int64, blockLikeSum uint64, adjFactor float64,
) float64 {
	if blockLikeSum == 0 {
		return 0.0
	}
	periodLength := ch.PeriodBlocks
	e := estimateLikesInPeriod(prevPeriodLikeCount, currPeriodLikeCount, remainingBlocks, periodLength)
	adj := interpolate(adjFactor, 2*adjFactor-1, 1-float64(remainingBlocks)/float64(periodLength))
	coinsAmountPerLike := float64(ch.Price.Amount.Int64()) / e * adj
	expense := coinsAmountPerLike * float64(blockLikeSum)
	remainingCoinAmount := float64(sub.Remaining.Amount.Int64())
	if remainingCoinAmount < expense {
		coinsAmountPerLike = remainingCoinAmount / float64(blockLikeSum)
	}
	return coinsAmountPerLike
}
