package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagChannelID = "channel-id"
)

var (
	fsChannelID                = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsChannelID.String(FlagChannelID, "", "The channel ID to be subscribed")
}
