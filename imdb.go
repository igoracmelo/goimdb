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
	ID          string
	ResultType  string
	Image       string
	Title       string
	Description string
}

type SearchTitleData struct {
	Results []SearchTitleResult
}

type SearchTitleResult struct {
	ID   string
	Type string

	// Optional.
	TrailerID string

	Title string

	// Optional.
	Image string

	// Optional.
	Description string
	Genres      []string
	Rating      float64
	VoteCount   int
	ReleaseYear int

	// Optional.
	EndYear int
}

type Client struct {
	BaseURL string
	HTTP    http.Client
}

func NewClient() Client {
	return Client{
		BaseURL: "https://www.imdb.com",
	}
}

func (c Client) SearchTitle(expr string) (data SearchTitleData, err error) {
	vals := url.Values{}
	vals.Set("title", expr)

	var b []byte

	coll := colly.NewCollector()

	coll.OnHTML("script#__NEXT_DATA__", func(h *colly.HTMLElement) {
		b = []byte(h.Text)
	})

	err = coll.Visit(c.BaseURL + "/search/title/?" + vals.Encode())
	if err != nil {
		return
	}

	dto := struct {
		Props struct {
			PageProps struct {
				SearchResults struct {
					TitleResults struct {
						TitleListItems []struct {
							Creators     []any    `json:"creators"`
							Directors    []any    `json:"directors"`
							EndYear      int      `json:"endYear"`
							Genres       []string `json:"genres"`
							TitleID      string   `json:"titleId"`
							TitleText    string   `json:"titleText"`
							Plot         string   `json:"plot"`
							PrimaryImage struct {
								Caption string `json:"caption"`
								ID      string `json:"id"`
								Height  int    `json:"height"`
								URL     string `json:"url"`
								Width   int    `json:"width"`
							} `json:"primaryImage"`
							RatingSummary struct {
								AggregateRating float64 `json:"aggregateRating"`
								VoteCount       int     `json:"voteCount"`
							} `json:"ratingSummary"`
							ReleaseYear int `json:"releaseYear"`
							Runtime     int `json:"runtime"`
							TitleType   struct {
								CanHaveEpisodes bool   `json:"canHaveEpisodes"`
								ID              string `json:"id"`
								Text            string `json:"text"`
							} `json:"titleType"`
							TopCast   []any  `json:"topCast"`
							TrailerID string `json:"trailerId"`
						} `json:"titleListItems"`
					} `json:"titleResults"`
				} `json:"searchResults"`
			} `json:"pageProps"`
		} `json:"props"`
	}{}

	err = json.Unmarshal(b, &dto)
	if err != nil {
		return
	}

	for _, i := range dto.Props.PageProps.SearchResults.TitleResults.TitleListItems {
		data.Results = append(data.Results, SearchTitleResult{
			ID:          i.TitleID,
			Type:        i.TitleType.ID,
			Image:       i.PrimaryImage.URL,
			Title:       i.TitleText,
			Description: i.Plot,
			TrailerID:   i.TrailerID,
			Genres:      i.Genres,
			Rating:      i.RatingSummary.AggregateRating,
			VoteCount:   i.RatingSummary.VoteCount,
			ReleaseYear: i.ReleaseYear,
			EndYear:     i.EndYear,
		})
	}

	return
}
