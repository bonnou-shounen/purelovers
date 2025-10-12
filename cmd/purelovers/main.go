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
		kong.Vars{"version": "0.1.3"},
		kong.ShortUsageOnError(),
	)

	if arg.Login != "" {
		os.Setenv("PURELOVERS_LOGIN", arg.Login)
	}

	if arg.Password != "" {
		os.Setenv("PURELOVERS_PASSWORD", arg.Password)
	}

	ctx.FatalIfErrorf(ctx.Run(&arg.Option))
}
