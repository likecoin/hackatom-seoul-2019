package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/likecoin/hackatom-seoul-2019/chain/x/subscription"
	"github.com/likecoin/hackatom-seoul-2019/chain/x/subscription/types"
)

func GetCmdQuerySubscription(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "subscription [subscriber-addr] [channel-id]",
		Short: "Query a subscription",
		Long: strings.TrimSpace(`Query details about a subscription:

$ likecli query subscription subscription cosmos16s47cyy5w6ja07w42s3yxe7p37pdvcrr39sc8e 1
`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			subscriber, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			channelID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryStore(subscription.GetSubscriptionKey(subscriber, channelID), storeName)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return fmt.Errorf("No subscription data with address %s and channel ID %d found",
					subscriber, channelID)
			}

			sub := types.Subscription{}
			cdc.MustUnmarshalBinaryBare(res, &sub)
			return cliCtx.PrintOutput(sub)
		},
	}
}