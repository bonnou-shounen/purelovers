package cmd

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bonnou-shounen/purelovers"
	"github.com/bonnou-shounen/purelovers/cmd/purelovers/util"
)

type RestoreFavoriteShops struct{}

func (r *RestoreFavoriteShops) Run() error {
	newShops := r.readShops(os.Stdin)

	c, err := util.NewLoggedClient()
	if err != nil {
		return err
	}

	curShops, err := c.GetFavoriteShops()
	if err != nil {
		return err
	}

	delShops, addShops := r.shopsDiff(curShops, newShops)
	c.DeleteFavoriteShops(delShops)
	c.AddFavoriteShops(addShops)

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

//nolint:lll
func (r *RestoreFavoriteShops) shopsDiff(curShops, newShops []*purelovers.Shop) (delShops, addShops []*purelovers.Shop) {
	ic := len(curShops) - 1
	in := len(newShops) - 1

	for ic >= 0 && in >= 0 {
		curShop := curShops[ic]
		newShop := newShops[in]

		if curShop.ID == newShop.ID {
			ic--
			in--

			continue
		}

		delShops = append(delShops, curShop)
		ic--
	}

	if ic >= 0 {
		delShops = append(delShops, curShops[:ic+1]...)
	}

	return delShops, newShops[:in+1]
}
