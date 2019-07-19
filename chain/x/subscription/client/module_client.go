package client

import (
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/likecoin/hackatom-seoul-2019/chain/x/subscription/client/cli"
	"github.com/likecoin/hackatom-seoul-2019/chain/x/subscription/types"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	subscriptionQueryCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the subscription module",
	}
	subscriptionQueryCmd.AddCommand(client.GetCommands(
		cli.GetCmdQuerySubscription(mc.storeKey, mc.cdc),
	)...)

	return subscriptionQueryCmd

}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	subscriptionTxCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Subscription transaction subcommands",
	}

	subscriptionTxCmd.AddCommand(client.PostCommands(
		cli.GetCmdSubscribe(mc.cdc),
	)...)

	return subscriptionTxCmd
}
