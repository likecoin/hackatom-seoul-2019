package cli

import (
	"encoding/binary"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/likecoin/hackatom-seoul-2019/chain/x/civicliker"
)

func GetCmdQueryLikeCount(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "like-count [liker] [url]",
		Short: "Query like count",
		Long: strings.TrimSpace(`Query the number of like for a liker on a url:

$ likecli query civicliker like-count cosmos16s47cyy5w6ja07w42s3yxe7p37pdvcrr39sc8e "https://like.co"
`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			liker, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			url := args[1]
			res, err := cliCtx.QueryStore(civicliker.GetLikeRecordCountKey(liker, url), storeName)
			if err != nil {
				return err
			}

			count := uint64(0)
			if len(res) > 0 {
				count = binary.BigEndian.Uint64(res)
			}
			output := sdk.NewInt(int64(count))

			return cliCtx.PrintOutput(output)
		},
	}
}
