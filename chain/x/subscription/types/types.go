package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Channel struct {
	ID           uint64    `json:"id"`
	Price        sdk.Coins `json:"price"`
	PeriodBlocks int64     `json:"period_blocks"`
}

func (ch Channel) String() string {
	return fmt.Sprintf(`ID: %d
Price: %s
Period: %d blocks`, ch.ID, ch.Price, ch.PeriodBlocks)
}

type Subscription struct {
	Subscriber       sdk.AccAddress `json:"subscriber"`
	ChannelID        uint64         `json:"channel_id"`
	Remaining        sdk.Coins      `json:"remaining"`
	NextPaymentBlock int64          `json:"next_payment_block"`
}

func (sub Subscription) String() string {
	return fmt.Sprintf(`Subscriber: %s
ChannelID: %d
Remaining: %s
NextPaymentBlock: %d`, sub.Subscriber, sub.ChannelID, sub.Remaining, sub.NextPaymentBlock)
}
