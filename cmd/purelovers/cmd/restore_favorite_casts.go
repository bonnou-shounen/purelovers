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

type RestoreFavoriteCasts struct{}

func (r *RestoreFavoriteCasts) Run() error {
	ctx := context.Background()

	newCasts := r.readCasts(os.Stdin)

	c, err := util.NewLoggedClient(ctx)
	if err != nil {
		return fmt.Errorf("on NewLoggedClient(): %w", err)
	}

	curCasts, err := c.GetFavoriteCasts(ctx)
	if err != nil {
		return fmt.Errorf("on GetFavoriteCasts(): %w", err)
	}

	delCasts, addCasts := util.ListDiff(curCasts, newCasts,
		func(a, b *purelovers.Cast) bool { return a.ID == b.ID && a.Shop.ID == b.Shop.ID },
	)

	if err := c.DeleteFavoriteCasts(ctx, delCasts); err != nil {
		return fmt.Errorf("on DeleteFavoriteCasts(): %w", err)
	}

	if err := c.AddFavoriteCasts(ctx, addCasts); err != nil {
		return fmt.Errorf("on AddFavoriteCasts(): %w", err)
	}

	return nil
}

func (r *RestoreFavoriteCasts) readCasts(reader io.Reader) []*purelovers.Cast {
	var casts []*purelovers.Cast

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		fields = append(fields, "", "", "", "")

		castID, _ := strconv.Atoi(fields[0])
		shopID, _ := strconv.Atoi(fields[1])
		castName := fields[2]
		shopName := fields[3]

		if castID == 0 || shopID == 0 {
			continue
		}

		casts = append(casts,
			&purelovers.Cast{
				ID:   castID,
				Name: castName,
				Shop: &purelovers.Shop{
					ID:   shopID,
					Name: shopName,
				},
			},
		)
	}

	return casts
}
