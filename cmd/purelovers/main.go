package main

import (
	"os"

	"github.com/alecthomas/kong"

	"github.com/bonnou-shounen/purelovers/cmd/purelovers/cmd"
)

func main() {
	arg := cmd.Arg{}
	ctx := kong.Parse(
		&arg,
		kong.Name("purelovers"),
		kong.Vars{"version": "0.1.2"},
		kong.ShortUsageOnError(),
	)

	if arg.Option.Login != "" {
		os.Setenv("PURELOVERS_LOGIN", arg.Option.Login)
	}

	if arg.Option.Password != "" {
		os.Setenv("PURELOVERS_PASSWORD", arg.Option.Password)
	}

	ctx.FatalIfErrorf(ctx.Run(&arg.Option))
}
