package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type QueryResSubscriptionPool struct {
	Value sdk.Coins `json:"value"`
}

// implement fmt.Stringer
func (r QueryResSubscriptionPool) String() string {
	return r.Value.String()
}
