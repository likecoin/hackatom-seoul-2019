package types

import (
	"net/url"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = MsgUploadContent{}

type MsgUploadContent struct {
	Author sdk.AccAddress `json:"author"`
	Url    string         `json:"url"`
}

func (msg MsgUploadContent) Route() string { return RouterKey }

func (msg MsgUploadContent) Type() string { return "upload_content" }

func (msg MsgUploadContent) ValidateBasic() sdk.Error {
	if msg.Author.Empty() {
		return sdk.ErrInvalidAddress(msg.Author.String())
	}
	_, err := url.Parse(msg.Url)
	if err != nil {
		return sdk.ErrUnknownRequest("Invalid URL")
	}
	return nil
}

func (msg MsgUploadContent) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgUploadContent) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Author}
}
