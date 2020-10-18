package purelovers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/remeh/sizedwaitgroup"
)

func (c *Client) GetFavoriteCasts() ([]*Cast, error) {
	casts, lastPage, err := c.getFavoriteCasts(1)
	if err != nil {
		return nil, err
	}

	if lastPage <= 1 {
		return casts, nil
	}

	castsList := make([][]*Cast, lastPage+1)
	swg := sizedwaitgroup.New(3)

	for page := 2; page <= lastPage; page++ {
		swg.Add()

		go func(p int) {
			defer swg.Done()

			castsList[p], _, _ = c.getFavoriteCasts(p)
		}(page)
	}
	swg.Wait()

	for page := 2; page <= lastPage; page++ {
		casts = append(casts, castsList[page]...)
	}

	return casts, nil
}

func (c *Client) getFavoriteCasts(page int) ([]*Cast, int, error) {
	strURL := "https://www.purelovers.com/user/favorite-girl/"
	if page > 1 {
		strURL += fmt.Sprintf("index/page/%d/", page)
	}

	resp, err := c.http.Get(strURL)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	var casts []*Cast

	doc.Find("div.girlList-girlDate a").Each(func(j int, a *goquery.Selection) {
		href, _ := a.Attr("href")
		shopID := c.getNum(href, "/shop/", "/")
		castID := c.getNum(href, "/girl/", "/")

		name := a.Text()
		if i := strings.Index(name, "\u00A0"); i >= 0 {
			name = name[:i]
		}

		if name != "" && shopID != 0 && castID != 0 {
			casts = append(casts, &Cast{ShopID: shopID, CastID: castID, Name: name})
		}
	})

	if page > 1 {
		return casts, 0, nil
	}

	href, _ := doc.Find("ul.page-move li:last-child a").Attr("href")
	lastPage := c.getNum(href, "/page/", "/")

	return casts, lastPage, nil
}

func (c *Client) getNum(str, prefix, suffix string) int {
	if i := strings.Index(str, prefix); i >= 0 {
		str = str[i+len(prefix):]
		if j := strings.Index(str, suffix); j >= 0 {
			str = str[:j]
		}
	}

	num, _ := strconv.Atoi(str)

	return num
}

func (c *Client) AddFavoriteCast(cast *Cast) error {
	return c.Ajax("https://www.purelovers.com/ajax/user/regist-favorite-girl/", cast.urlValues())
}

func (c *Client) DeleteFavoriteCast(cast *Cast) error {
	return c.Ajax("https://www.purelovers.com/ajax/user-my-page/girl-delete/", cast.urlValues())
}

func (c *Client) AddFavoriteCasts(casts []*Cast) {
	for i := len(casts) - 1; i >= 0; i-- {
		c.AddFavoriteCast(casts[i]) //nolint:errcheck
	}
}

func (c *Client) DeleteFavoriteCasts(casts []*Cast) {
	swg := sizedwaitgroup.New(5)

	for _, cast := range casts {
		swg.Add()

		go func(cast *Cast) {
			defer swg.Done()

			c.DeleteFavoriteCast(cast) //nolint:errcheck
		}(cast)
	}

	swg.Wait()
}
