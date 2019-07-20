package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = MsgLike{}

type MsgLike struct {
	Liker sdk.AccAddress `json:"liker"`
	Likee sdk.AccAddress `json:"likee"`
	Url   string         `json:"url"`
	Count uint64         `json:"count"`
}

func (msg MsgLike) Route() string { return RouterKey }

func (msg MsgLike) Type() string { return "like" }

func (msg MsgLike) ValidateBasic() sdk.Error {
	if msg.Liker.Empty() {
		return sdk.ErrInvalidAddress(msg.Liker.String())
	}
	if msg.Likee.Empty() {
		return sdk.ErrInvalidAddress(msg.Likee.String())
	}
	if msg.Count <= 0 {
		return sdk.ErrUnknownRequest("invalid like count")
	}
	return nil
}

func (msg MsgLike) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgLike) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Liker}
}
