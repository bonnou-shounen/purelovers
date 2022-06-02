package purelovers

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jesse0michael/errgroup"
)

func (c *Client) GetFavoriteCasts(ctx context.Context) ([]*Cast, error) {
	var lastPage int

	casts, err := c.getFavoriteCastsOnPage(ctx, 1, &lastPage)
	if err != nil {
		return nil, fmt.Errorf("on getFavoriteCastsOnPage(1): %w", err)
	}

	if lastPage <= 1 {
		return casts, nil
	}

	castsOnPage := make([][]*Cast, lastPage+1)
	eg, egCtx := errgroup.WithContext(ctx, 3)

	for page := 2; page <= lastPage; page++ {
		page := page

		eg.Go(func() error {
			casts, err := c.getFavoriteCastsOnPage(egCtx, page, nil)
			if err != nil {
				return fmt.Errorf("on getFavoriteCastsOnPage(%d): %w", page, err)
			}

			castsOnPage[page] = casts

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("on goroutine: %w", err)
	}

	for page := 2; page <= lastPage; page++ {
		casts = append(casts, castsOnPage[page]...)
	}

	if err := c.getShopNames(ctx, casts); err != nil {
		return nil, fmt.Errorf("on getShopNames(): %w", err)
	}

	return casts, nil
}

func (c *Client) getFavoriteCastsOnPage(ctx context.Context, page int, pLastPage *int) ([]*Cast, error) {
	strURL := fmt.Sprintf("https://purelovers.com/user/favorite-girl/pg%d/", page)

	resp, err := c.get(ctx, strURL, "")
	if err != nil {
		return nil, fmt.Errorf(`on get("%s"): %w`, strURL, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("on NewDocumentFromReader(): %w", err)
	}

	var casts []*Cast

	doc.Find("div.k_row-grid--small a.k_box--scale").Each(func(_ int, a *goquery.Selection) {
		href, _ := a.Attr("href")
		shopID := c.parseNumber(href, "/shop/", "/")
		castID := c.parseNumber(href, "/girl/", "/")
		castName, _ := a.Find("img").Attr("alt")

		if castID == 0 || castName == "" || shopID == 0 {
			return
		}

		casts = append(casts,
			&Cast{
				ID:   castID,
				Name: castName,
				Shop: &Shop{
					ID: shopID,
				},
			},
		)
	})

	if pLastPage != nil {
		text, _ := doc.Find("div.k_searchResult p").Html()
		*pLastPage = (c.parseNumber(text, "", "件") + 19) / 20
	}

	return casts, nil
}

func (c *Client) AddFavoriteCast(ctx context.Context, cast *Cast) error {
	time.Sleep(1 * time.Second)

	return c.modFavoriteCast(ctx, cast, "set-follow")
}

func (c *Client) DeleteFavoriteCast(ctx context.Context, cast *Cast) error {
	return c.modFavoriteCast(ctx, cast, "delete-follow")
}

func (c *Client) modFavoriteCast(ctx context.Context, cast *Cast, operation string) error {
	strURL := fmt.Sprintf(
		"https://purelovers.com/shop/%d/girl/%d/%s/?_=%d",
		cast.Shop.ID,
		cast.ID,
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

func (c *Client) AddFavoriteCasts(ctx context.Context, casts []*Cast) error {
	for i := len(casts) - 1; i >= 0; i-- {
		cast := casts[i]

		if err := c.AddFavoriteCast(ctx, cast); err != nil {
			return fmt.Errorf("on AddFavoriteCast(%d=%s): %w", cast.ID, cast.Name, err)
		}
	}

	return nil
}

func (c *Client) DeleteFavoriteCasts(ctx context.Context, casts []*Cast) error {
	eg, egCtx := errgroup.WithContext(ctx, 5)

	for _, cast := range casts {
		cast := cast

		eg.Go(func() error {
			if err := c.DeleteFavoriteCast(egCtx, cast); err != nil {
				return fmt.Errorf("on DeleteFavoriteCast(%d=%s): %w", cast.ID, cast.Name, err)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("on goroutine: %w", err)
	}

	return nil
}

func (c *Client) getShopNames(ctx context.Context, casts []*Cast) error {
	shopNameOf := map[int]string{}
	for _, cast := range casts {
		shopNameOf[cast.Shop.ID] = ""
	}

	eg, egCtx := errgroup.WithContext(ctx, 3)

	for shopID := range shopNameOf {
		shopID := shopID

		eg.Go(func() error {
			shopName, err := c.getShopName(egCtx, shopID)
			if err != nil {
				return fmt.Errorf("on getShopName(%d): %w", shopID, err)
			}

			shopNameOf[shopID] = shopName

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("on goroutine: %w", err)
	}

	for _, cast := range casts {
		cast.Shop.Name = shopNameOf[cast.Shop.ID]
	}

	return nil
}
