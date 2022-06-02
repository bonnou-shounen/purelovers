package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bonnou-shounen/purelovers"
	"github.com/bonnou-shounen/purelovers/cmd/purelovers/util"
)

type RestoreFavoriteShops struct{}

func (r *RestoreFavoriteShops) Run() error {
	ctx := context.Background()

	newShops := r.readShops(os.Stdin)

	c, err := util.NewLoggedClient(ctx)
	if err != nil {
		return fmt.Errorf("on NewLoggedClient(): %w", err)
	}

	curShops, err := c.GetFavoriteShops(ctx)
	if err != nil {
		return fmt.Errorf("on GetFavoriteShops(): %w", err)
	}

	delShops, addShops := util.ListDiff(curShops, newShops,
		func(a, b *purelovers.Shop) bool { return a.ID == b.ID },
	)

	if err := c.DeleteFavoriteShops(ctx, delShops); err != nil {
		return fmt.Errorf("on DeleteFavoriteShops(): %w", err)
	}

	if err := c.AddFavoriteShops(ctx, addShops); err != nil {
		return fmt.Errorf("on AddFavoriteShops(): %w", err)
	}

	return nil
}

func (r *RestoreFavoriteShops) readShops(reader io.Reader) []*purelovers.Shop {
	var shops []*purelovers.Shop

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		fields = append(fields, "", "")

		shopID, _ := strconv.Atoi(fields[0])
		shopName := fields[1]

		if shopID == 0 {
			continue
		}

		shops = append(shops,
			&purelovers.Shop{
				ID:   shopID,
				Name: shopName,
			},
		)
	}

	return shops
}
