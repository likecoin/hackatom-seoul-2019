package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagLikee = "likee"
	FlagUrl   = "url"
	FlagCount = "count"
)

var (
	fsLikee = flag.NewFlagSet("", flag.ContinueOnError)
	fsUrl   = flag.NewFlagSet("", flag.ContinueOnError)
	fsCount = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsLikee.String(FlagLikee, "", "The author of the URL")
	fsUrl.String(FlagUrl, "", "The URL to be liked")
	fsCount.Uint64(FlagCount, 1, "Number of likes to give")
}
