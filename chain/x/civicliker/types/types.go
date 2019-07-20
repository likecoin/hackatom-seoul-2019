package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type LikeRecord struct {
	Liker           sdk.AccAddress `json:"liker"`
	Likee           sdk.AccAddress `json:"likee"`
	Url             string         `json:"url"`
	Count           uint64         `json:"count"`
	CoinDistributed sdk.Coin       `json:"coin_distributed"`
	Time            time.Time      `json:"time"`
}

type LikerState struct {
	Liker                 sdk.AccAddress
	PrevPeriodLikeCount   uint64
	CurrPeriodLikeCount   uint64
	SubscriptionChannelID uint64
	PeriodEndingBlock     int64
}
