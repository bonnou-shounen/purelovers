package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/bonnou-shounen/purelovers/cmd/purelovers/util"
)

type DumpFavoriteCasts struct{}

func (d *DumpFavoriteCasts) Run() error {
	ctx := context.Background()

	c, err := util.NewLoggedClient(ctx)
	if err != nil {
		return fmt.Errorf("on NewLoggedClient(): %w", err)
	}

	casts, err := c.GetFavoriteCasts(ctx)
	if err != nil {
		return fmt.Errorf("on GetFavoriteCasts(): %w", err)
	}

	for _, cast := range casts {
		fmt.Fprintf(os.Stdout, "%d\t%d\t%s\t%s\n", cast.ID, cast.Shop.ID, cast.Name, cast.Shop.Name)
	}

	return nil
}
