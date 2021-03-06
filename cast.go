package purelovers

import (
	"fmt"
	"net/url"
)

type Cast struct {
	ID   int
	Name string
	Shop *Shop
}

func (c *Cast) urlValues() url.Values {
	return url.Values{
		"girl_id": []string{fmt.Sprint(c.ID)},
		"shop_id": []string{fmt.Sprint(c.Shop.ID)},
	}
}
