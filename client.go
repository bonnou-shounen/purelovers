package purelovers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	http *http.Client
}

func NewClient() *Client {
	jar, _ := cookiejar.New(nil)

	return &Client{
		http: &http.Client{Jar: jar},
	}
}

func (c *Client) Login(ctx context.Context, id, password string) error {
	values := url.Values{
		"mail_address":  []string{id},
		"password":      []string{password},
		"submit_button": []string{"ログイン"},
	}

	strURL := "https://purelovers.com/user/login"

	resp, err := c.post(ctx, strURL, values.Encode())
	if err != nil {
		return fmt.Errorf(`on post("%s"): %w`, strURL, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("on NewDocumentFromReader(): %w", err)
	}

	title := doc.Find(`title`).Text()
	if !strings.Contains(title, "マイページ") {
		return fmt.Errorf("login failed: [%s]", title)
	}

	return nil
}

func (c *Client) get(ctx context.Context, strURL string, query string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprint(strURL, "?", query), nil)
	if err != nil {
		return nil, fmt.Errorf("on NewRequest(): %w", err)
	}

	req.Header.Set("x-requested-with", "XMLHttpRequest")

	resp, err := c.http.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("on http.Do(): %w", err)
	}

	return resp, nil
}

func (c *Client) post(ctx context.Context, strURL string, form string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, strURL, strings.NewReader(form))
	if err != nil {
		return nil, fmt.Errorf("on NewRequest(): %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("on http.Do(): %w", err)
	}

	return resp, nil
}

func (c *Client) parseNumber(str, prefix, suffix string) int {
	if i := strings.Index(str, prefix); i >= 0 {
		str = str[i+len(prefix):]
		if j := strings.Index(str, suffix); j >= 0 {
			str = str[:j]
		}
	}

	num, _ := strconv.Atoi(str)

	return num
}
