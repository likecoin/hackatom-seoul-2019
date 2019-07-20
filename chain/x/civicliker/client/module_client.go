package client

import (
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/likecoin/hackatom-seoul-2019/chain/x/civicliker/client/cli"
	"github.com/likecoin/hackatom-seoul-2019/chain/x/civicliker/types"
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
	civiclikerQueryCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the civicliker module",
	}
	civiclikerQueryCmd.AddCommand(client.GetCommands(
		cli.GetCmdQueryLikeCount(mc.storeKey, mc.cdc),
	)...)

	return civiclikerQueryCmd

}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	civiclikerTxCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Civicliker transaction subcommands",
	}

	civiclikerTxCmd.AddCommand(client.PostCommands(
		cli.GetCmdLike(mc.cdc),
	)...)

	return civiclikerTxCmd
}
