package subscription

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/likecoin/hackatom-seoul-2019/chain/x/subscription/types"
)

type GenesisState struct {
	Channels []types.Channel `json:"channels"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Channels: []types.Channel{
			types.Channel{
				ID:           1,
				Price:        sdk.NewCoin("nanolike", sdk.NewInt(100000000000)), // 100 LIKE
				PeriodBlocks: 12,
			},
		},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, ch := range data.Channels {
		keeper.SetChannel(ctx, ch)
	}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	channels := []types.Channel{}
	keeper.IterateChannels(ctx, func(index int64, ch types.Channel) (stop bool) {
		channels = append(channels, ch)
		return false
	})
	return GenesisState{
		Channels: channels,
	}
}
