package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/bonnou-shounen/purelovers/cmd/purelovers/util"
)

type DumpFavoriteShops struct{}

func (d *DumpFavoriteShops) Run() error {
	ctx := context.Background()

	c, err := util.NewLoggedClient(ctx)
	if err != nil {
		return fmt.Errorf("on NewLoggedClient(): %w", err)
	}

	shops, err := c.GetFavoriteShops(ctx)
	if err != nil {
		return fmt.Errorf("on GetFavoriteShops(): %w", err)
	}

	for _, shop := range shops {
		fmt.Fprintf(os.Stdout, "%d\t%s\n", shop.ID, shop.Name)
	}

	return nil
}
