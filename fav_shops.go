package purelovers

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/remeh/sizedwaitgroup"
)

func (c *Client) GetFavoriteShops() ([]*Shop, error) {
	var lastPage int

	casts, err := c.getFavoriteShopsOnPage(1, &lastPage)
	if err != nil {
		return nil, err
	}

	if lastPage <= 1 {
		return casts, nil
	}

	shopsOnPage := make([][]*Shop, lastPage+1)
	swg := sizedwaitgroup.New(3)

	for page := 2; page <= lastPage; page++ {
		swg.Add()

		go func(page int) {
			defer swg.Done()

			shopsOnPage[page], _ = c.getFavoriteShopsOnPage(page, nil)
		}(page)
	}
	swg.Wait()

	for page := 2; page <= lastPage; page++ {
		casts = append(casts, shopsOnPage[page]...)
	}

	return casts, nil
}

func (c *Client) getFavoriteShopsOnPage(page int, pLastPage *int) ([]*Shop, error) {
	resp, err := c.http.Get(fmt.Sprintf("https://www.purelovers.com/user/favorite-shop/index/page/%d/", page))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var shops []*Shop

	doc.Find("p.shopList-nameDate a").Each(func(j int, a *goquery.Selection) {
		href, _ := a.Attr("href")
		shopID := c.parseNumber(href, "/shop/", "/")
		name := a.Text()

		if name != "" && shopID != 0 {
			shops = append(shops, &Shop{ShopID: shopID, Name: name})
		}
	})

	if pLastPage != nil {
		href, _ := doc.Find("ul.page-move li:last-child a").Attr("href")
		*pLastPage = c.parseNumber(href, "/page/", "/")
	}

	return shops, nil
}

func (c *Client) AddFavoriteShop(shop *Shop) error {
	return c.ajax("https://www.purelovers.com/ajax/user/regist-favorite-shop/", shop.urlValues())
}

func (c *Client) DeleteFavoriteShop(shop *Shop) error {
	return c.ajax("https://www.purelovers.com/ajax/user-my-page/shop-delete/", shop.urlValues())
}

func (c *Client) AddFavoriteShops(shops []*Shop) {
	for i := len(shops) - 1; i >= 0; i-- {
		c.AddFavoriteShop(shops[i]) //nolint:errcheck
	}
}

func (c *Client) DeleteFavoriteShops(shops []*Shop) {
	swg := sizedwaitgroup.New(5)

	for _, shop := range shops {
		swg.Add()

		go func(shop *Shop) {
			defer swg.Done()

			c.DeleteFavoriteShop(shop) //nolint:errcheck
		}(shop)
	}

	swg.Wait()
}
