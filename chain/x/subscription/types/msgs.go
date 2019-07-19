package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = MsgSubscribe{}

type MsgSubscribe struct {
	Subscriber sdk.AccAddress `json:"subscriber"`
	ChannelID  uint64         `json:"channel_id"`
}

func (msg MsgSubscribe) Route() string { return RouterKey }

func (msg MsgSubscribe) Type() string { return "subscribe" }

func (msg MsgSubscribe) ValidateBasic() sdk.Error {
	if msg.Subscriber.Empty() {
		return sdk.ErrInvalidAddress(msg.Subscriber.String())
	}
	return nil
}

func (msg MsgSubscribe) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSubscribe) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Subscriber}
}
