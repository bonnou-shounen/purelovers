package purelovers

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jesse0michael/errgroup"
)

func (c *Client) GetFavoriteShops(ctx context.Context) ([]*Shop, error) {
	var lastPage int

	shops, err := c.getFavoriteShopsOnPage(ctx, 1, &lastPage)
	if err != nil {
		return nil, fmt.Errorf("on getFavoriteShopsOnPage(1): %w", err)
	}

	if lastPage <= 1 {
		return shops, nil
	}

	shopsOnPage := make([][]*Shop, lastPage+1)
	eg, egCtx := errgroup.WithContext(ctx, 3)

	for page := 2; page <= lastPage; page++ {
		page := page

		eg.Go(func() error {
			shops, err := c.getFavoriteShopsOnPage(egCtx, page, nil)
			if err != nil {
				return fmt.Errorf("on getFavoriteShopsOnPage(%d): %w", page, err)
			}

			shopsOnPage[page] = shops

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("on goroutine: %w", err)
	}

	for page := 2; page <= lastPage; page++ {
		shops = append(shops, shopsOnPage[page]...)
	}

	return shops, nil
}

func (c *Client) getFavoriteShopsOnPage(ctx context.Context, page int, pLastPage *int) ([]*Shop, error) {
	strURL := fmt.Sprintf("https://purelovers.com/user/favorite-shop/pg%d/", page)

	resp, err := c.get(ctx, strURL, "")
	if err != nil {
		return nil, fmt.Errorf(`on get("%s"): %w`, strURL, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("on NewDocumentFromReader(): %w", err)
	}

	var shops []*Shop

	doc.Find("div.k_row-grid--small a.k_box").Each(func(_ int, a *goquery.Selection) {
		href, _ := a.Attr("href")
		shopID := c.parseNumber(href, "/shop/", "/")
		shopName := a.Find("span.k_text--xlarge").Text()

		if shopID == 0 || shopName == "" {
			return
		}

		shops = append(shops,
			&Shop{
				ID:   shopID,
				Name: shopName,
			},
		)
	})

	if pLastPage != nil {
		text, _ := doc.Find("div.k_searchResult p").Html()
		*pLastPage = (c.parseNumber(text, "", "件") + 19) / 20
	}

	return shops, nil
}

func (c *Client) AddFavoriteShop(ctx context.Context, shop *Shop) error {
	time.Sleep(1 * time.Second)

	return c.modFavoriteShop(ctx, shop, "set-follow")
}

func (c *Client) DeleteFavoriteShop(ctx context.Context, shop *Shop) error {
	return c.modFavoriteShop(ctx, shop, "delete-follow")
}

func (c *Client) modFavoriteShop(ctx context.Context, shop *Shop, operation string) error {
	strURL := fmt.Sprintf(
		"https://purelovers.com/shop/%d/%s/?_=%d",
		shop.ID,
		operation,
		time.Now().UnixMilli(),
	)

	resp, err := c.get(ctx, strURL, "")
	if err != nil {
		return fmt.Errorf(`on get("%s"): %w`, strURL, err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("on ReadAll(): %w", err)
	}

	if string(b) != "true" {
		return fmt.Errorf("on ReadAll(): invalid response [%s]", b)
	}

	return nil
}

func (c *Client) AddFavoriteShops(ctx context.Context, shops []*Shop) error {
	for i := len(shops) - 1; i >= 0; i-- {
		shop := shops[i]

		if err := c.AddFavoriteShop(ctx, shop); err != nil {
			return fmt.Errorf("on AddFavoriteShop(%d=%s): %w", shop.ID, shop.Name, err)
		}
	}

	return nil
}

func (c *Client) DeleteFavoriteShops(ctx context.Context, shops []*Shop) error {
	eg, egCtx := errgroup.WithContext(ctx, 5)

	for _, shop := range shops {
		shop := shop

		eg.Go(func() error {
			if err := c.DeleteFavoriteShop(egCtx, shop); err != nil {
				return fmt.Errorf("on DeleteFavoriteShop(%d=%s): %w", shop.ID, shop.Name, err)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("on goroutine: %w", err)
	}

	return nil
}
