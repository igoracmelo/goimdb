package imdb

import (
	"net/http"
	"net/url"
)

type SearchData struct {
	SearchType   string
	Expression   string
	Results      []SearchResult
	ErrorMessage string
}

type SearchResult struct {
	ID          string
	ResultType  string
	Image       string
	Title       string
	Description string
}

type Client struct {
	BaseURL string
	HTTP    http.Client
}

func NewClient() Client {
	return Client{}
}

func (c Client) SearchTitle(expr string) (data SearchData, err error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/find/", nil)
	if err != nil {
		return
	}
	vals := url.Values{}
	vals.Set("q", expr)
	req.URL.RawQuery = vals.Encode()

	_, err = c.HTTP.Do(req)
	return SearchData{}, err
}
