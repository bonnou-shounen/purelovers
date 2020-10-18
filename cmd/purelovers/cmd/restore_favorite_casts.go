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

type RestoreFavoriteCasts struct{}

func (r *RestoreFavoriteCasts) Run() error {
	newCasts := r.readCasts(os.Stdin)

	c, err := util.NewLoggedClient()
	if err != nil {
		return err
	}

	curCasts, err := c.GetFavoriteCasts()
	if err != nil {
		return err
	}

	delCasts, addCasts := r.castsDiff(curCasts, newCasts)
	c.DeleteFavoriteCasts(delCasts)
	c.AddFavoriteCasts(addCasts)

	return nil
}

func (r *RestoreFavoriteCasts) readCasts(reader io.Reader) []*purelovers.Cast {
	casts := make([]*purelovers.Cast, 0)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}

		shopID, _ := strconv.Atoi(fields[0])
		castID, _ := strconv.Atoi(fields[1])

		if !(shopID > 0 && castID > 0) {
			continue
		}

		casts = append(casts, &purelovers.Cast{ShopID: shopID, CastID: castID})
	}

	return casts
}

//nolint:lll
func (r *RestoreFavoriteCasts) castsDiff(curCasts, newCasts []*purelovers.Cast) (delCasts, addCasts []*purelovers.Cast) {
	ic := len(curCasts) - 1
	in := len(newCasts) - 1

	for ic >= 0 && in >= 0 {
		curCast := curCasts[ic]
		newCast := newCasts[in]

		if curCast.ShopID == newCast.ShopID && curCast.CastID == newCast.CastID {
			ic--
			in--

			continue
		}

		delCasts = append(delCasts, curCast)
		ic--
	}

	if ic >= 0 {
		delCasts = append(delCasts, curCasts[:ic+1]...)
	}

	return delCasts, newCasts[:in+1]
}
