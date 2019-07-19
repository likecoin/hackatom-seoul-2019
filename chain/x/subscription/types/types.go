package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Content struct {
	Author sdk.AccAddress `json:"author"`
	Url    string         `json:"url"`
}

func (content Content) String() string {
	return fmt.Sprintf(`Author: %s
Url: %s`, content.Author, content.Url)
}
