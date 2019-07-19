package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/likecoin/hackatom-seoul-2019/chain/x/subscription/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetCmdSubscribe(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscribe",
		Short: "subscribe to a subscription channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			txBldr, msg, err := BuildSubscribeMsg(cliCtx, txBldr)
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, true)
		},
	}

	cmd.Flags().AddFlagSet(fsChannelID)
	cmd.MarkFlagRequired(FlagChannelID)

	return cmd
}

func BuildSubscribeMsg(cliCtx context.CLIContext, txBldr authtxb.TxBuilder) (authtxb.TxBuilder, sdk.Msg, error) {
	channelIDString := viper.GetString(FlagChannelID)
	channelID, err := strconv.ParseUint(channelIDString, 10, 64)
	if err != nil {
		return txBldr, nil, err
	}

	subscriber := cliCtx.GetFromAddress()
	msg := types.MsgSubscribe{
		Subscriber: subscriber,
		ChannelID:  channelID,
	}

	return txBldr, msg, nil
}
