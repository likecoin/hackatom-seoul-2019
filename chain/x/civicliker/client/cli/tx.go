package cli

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/likecoin/hackatom-seoul-2019/chain/x/civicliker/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetCmdLike(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "like",
		Short: "like a URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			txBldr, msg, err := BuildLikeMsg(cliCtx, txBldr)
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, true)
		},
	}

	cmd.Flags().AddFlagSet(fsLikee)
	cmd.Flags().AddFlagSet(fsUrl)
	cmd.Flags().AddFlagSet(fsCount)
	cmd.MarkFlagRequired(FlagLikee)
	cmd.MarkFlagRequired(FlagUrl)

	return cmd
}

func BuildLikeMsg(cliCtx context.CLIContext, txBldr authtxb.TxBuilder) (authtxb.TxBuilder, sdk.Msg, error) {
	likeeStr := viper.GetString(FlagLikee)
	likee, err := sdk.AccAddressFromBech32(likeeStr)
	if err != nil {
		return txBldr, nil, err
	}

	url := viper.GetString(FlagUrl)
	count := viper.GetInt64(FlagCount)
	if count <= 0 {
		return txBldr, nil, errors.New("Invalid like count")
	}

	liker := cliCtx.GetFromAddress()
	msg := types.MsgLike{
		Liker: liker,
		Likee: likee,
		Url:   url,
		Count: uint64(count),
	}

	return txBldr, msg, nil
}
