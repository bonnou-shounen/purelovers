package cmd

import (
	"fmt"
	"os"

	"github.com/bonnou-shounen/purelovers/cmd/purelovers/util"
)

type DumpFavoriteShops struct{}

func (d *DumpFavoriteShops) Run() error {
	c, err := util.NewLoggedClient()
	if err != nil {
		return err
	}

	shops, err := c.GetFavoriteShops()
	if err != nil {
		return err
	}

	for _, shop := range shops {
		fmt.Fprintf(os.Stdout, "%d\t%s\n", shop.ID, shop.Name)
	}

	return nil
}
