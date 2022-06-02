package purelovers

import (
	"context"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type Shop struct {
	ID   int
	Name string
}

func (c *Client) getShopName(ctx context.Context, shopID int) (string, error) {
	strURL := fmt.Sprint("https://purelovers.com/shop/", shopID, "/")

	resp, err := c.get(ctx, strURL, "")
	if err != nil {
		return "", fmt.Errorf(`on get("%s"): %w`, strURL, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("on NewDocumentFromReader(): %w", err)
	}

	shopName := doc.Find("h1.k_title").Text()

	return shopName, nil
}
