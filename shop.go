package purelovers

import (
	"fmt"
	"net/url"
)

type Shop struct {
	ID   int
	Name string
}

func (s *Shop) urlValues() url.Values {
	return url.Values{
		"shop_id": []string{fmt.Sprint(s.ID)},
	}
}
