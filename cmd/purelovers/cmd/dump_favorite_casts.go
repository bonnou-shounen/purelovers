package cmd

import (
	"fmt"
	"os"

	"github.com/bonnou-shounen/purelovers/cmd/purelovers/util"
)

type DumpFavoriteCasts struct{}

func (d *DumpFavoriteCasts) Run() error {
	c, err := util.NewLoggedClient()
	if err != nil {
		return err
	}

	casts, err := c.GetFavoriteCasts()
	if err != nil {
		return err
	}

	for _, cast := range casts {
		fmt.Fprintf(os.Stdout, "%d\t%d\t%s\n", cast.ShopID, cast.ID, cast.Name)
	}

	return nil
}
