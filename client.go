package purelovers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	http *http.Client
	ui   string
	uci  string
}

func NewClient() *Client {
	jar, _ := cookiejar.New(nil)

	return &Client{
		http: &http.Client{Jar: jar},
	}
}

func (c *Client) Login(id, password string) error {
	values := url.Values{
		"id":            []string{id},
		"password":      []string{password},
		"submit_button": []string{"ログイン"},
	}

	resp, err := c.http.PostForm("https://www.purelovers.com/user/login.html", values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	c.ui, _ = doc.Find(`input[name="ui"]`).Attr("value")
	c.uci, _ = doc.Find(`input[name="uci"]`).Attr("value")

	if c.ui == "" || c.uci == "" {
		return fmt.Errorf("login failed")
	}

	return nil
}

func (c *Client) Ajax(strURL string, values url.Values) error {
	values.Add("ui", c.ui)
	values.Add("uci", c.uci)

	req, err := http.NewRequest(http.MethodPost, strURL, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(bin) != "1" {
		return fmt.Errorf("ajax call returns: [%s]", bin)
	}

	return nil
}
