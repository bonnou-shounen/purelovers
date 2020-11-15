package purelovers

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/remeh/sizedwaitgroup"
)

func (c *Client) GetFavoriteCasts() ([]*Cast, error) {
	var lastPage int

	casts, err := c.getFavoriteCastsOnPage(1, &lastPage)
	if err != nil {
		return nil, err
	}

	if lastPage <= 1 {
		return casts, nil
	}

	castsOnPage := make([][]*Cast, lastPage+1)
	swg := sizedwaitgroup.New(3)

	for page := 2; page <= lastPage; page++ {
		swg.Add()

		go func(page int) {
			defer swg.Done()

			castsOnPage[page], _ = c.getFavoriteCastsOnPage(page, nil)
		}(page)
	}
	swg.Wait()

	for page := 2; page <= lastPage; page++ {
		casts = append(casts, castsOnPage[page]...)
	}

	return casts, nil
}

func (c *Client) getFavoriteCastsOnPage(page int, pLastPage *int) ([]*Cast, error) {
	resp, err := c.http.Get(fmt.Sprintf("https://www.purelovers.com/user/favorite-girl/index/page/%d/", page))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var casts []*Cast

	doc.Find("ul.myGirlList li").Each(func(_ int, li *goquery.Selection) {
		a := li.Find("div.girlList-girlDate a").First()
		href, _ := a.Attr("href")
		shopID := c.parseNumber(href, "/shop/", "/")
		castID := c.parseNumber(href, "/girl/", "/")
		castName := c.parseCastName(a.Text())
		shopName := li.Find("div.girlList-shopDate a").First().Text()

		if castID == 0 || castName == "" || shopID == 0 || shopName == "" {
			return
		}

		casts = append(casts,
			&Cast{
				ID:   castID,
				Name: castName,
				Shop: &Shop{
					ID:   shopID,
					Name: shopName,
				},
			},
		)
	})

	if pLastPage != nil {
		href, _ := doc.Find("ul.page-move li:last-child a").Attr("href")
		*pLastPage = c.parseNumber(href, "/page/", "/")
	}

	return casts, nil
}

func (c *Client) AddFavoriteCast(cast *Cast) error {
	return c.ajax("https://www.purelovers.com/ajax/user/regist-favorite-girl/", cast.urlValues())
}

func (c *Client) DeleteFavoriteCast(cast *Cast) error {
	return c.ajax("https://www.purelovers.com/ajax/user-my-page/girl-delete/", cast.urlValues())
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

func (c *Client) parseCastName(str string) string {
	if i := strings.Index(str, "\u00A0"); i >= 0 {
		str = str[:i]
	}

	return str
}
