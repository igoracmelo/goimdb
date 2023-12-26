package imdb

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gocolly/colly/v2"
)

type SearchData struct {
	SearchType   string
	Expression   string
	Results      []SearchResult
	ErrorMessage string
}

type SearchResult struct {
	ID string
	// ResultType string
	Image string
	Title string
	// Description string
}

type Client struct {
	BaseURL string
	HTTP    http.Client
}

func (c Client) SearchTitle(expr string) (data SearchData, err error) {
	vals := url.Values{}
	vals.Set("q", expr)
	vals.Set("s", "tt") // search type 'title' I guess

	var b []byte

	coll := colly.NewCollector()

	coll.OnHTML("script#__NEXT_DATA__", func(h *colly.HTMLElement) {
		b = []byte(h.Text)
	})

	err = coll.Visit(c.BaseURL + "/find/?" + vals.Encode())
	if err != nil {
		return
	}

	dto := struct {
		Props struct {
			PageProps struct {
				TitleResults struct {
					Results []struct {
						ID                    string `json:"id"`
						ImageType             string `json:"imageType"`
						TitleNameText         string `json:"titleNameText"`
						TitlePosterImageModel struct {
							Caption   string `json:"caption"`
							MaxHeight int    `json:"maxHeight"`
							MaxWidth  int    `json:"maxWidth"`
							URL       string `json:"url"`
						} `json:"titlePosterImageModel"`
						TitleReleaseText string   `json:"titleReleaseText"`
						TitleTypeText    string   `json:"titleTypeText"`
						TopCredits       []string `json:"topCredits"`
					} `json:"results"`
				} `json:"titleResults"`
			} `json:"pageProps"`
		} `json:"props"`
	}{}

	err = json.Unmarshal(b, &dto)
	if err != nil {
		return
	}

	for _, r := range dto.Props.PageProps.TitleResults.Results {
		data.Results = append(data.Results, SearchResult{
			ID:    r.ID,
			Image: r.TitlePosterImageModel.URL,
			Title: r.TitleNameText,
		})
	}

	return
}
