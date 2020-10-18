package purelovers

import (
	"fmt"
	"net/url"
)

type Cast struct {
	ShopID int
	CastID int
	Name   string
}

func (c *Cast) urlValues() url.Values {
	return url.Values{
		"shop_id": []string{fmt.Sprint(c.ShopID)},
		"girl_id": []string{fmt.Sprint(c.CastID)},
	}
}
