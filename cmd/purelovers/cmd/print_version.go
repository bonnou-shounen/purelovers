package cmd

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

type PrintVersion struct{}

func (PrintVersion) Run(vars kong.Vars) error {
	fmt.Fprintln(os.Stdout, vars["version"])

	return nil
}
