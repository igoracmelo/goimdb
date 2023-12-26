package imdb

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"text/template"
)

const responseTmpl = `
<script id="__NEXT_DATA__" type="application/json">
{
	"props": {
		"pageProps": {
			"searchResults": {
				"titleResults": {
					"titleListItems": [
						{
							"endYear": {{ .EndYear }},
							"genres": ["Animation", "Comedy"],
							"titleId": "{{ .ID }}",
							"titleText": "{{ .Title }}",
							"plot": "{{ .Description }}",
							"primaryImage": {
								"caption": "Matt Stone and Trey Parker in South Park (1997)",
								"id": "rm2471417089",
								"height": 900,
								"url": "{{ .Image }}",
								"width": 600
							},
							"ratingSummary": {
								"aggregateRating": {{ .Rating }},
								"voteCount": {{ .VoteCount }}
							},
							"releaseYear": {{ .ReleaseYear }},
							"runtime": 1320,
							"titleType": {
								"canHaveEpisodes": true,
								"id": "{{ .Type }}",
								"text": "TV Series"
							},
							"topCast": [],
							"trailerId": "{{ .TrailerID }}"
						}
					]
				}
			}
		}
	}
}
</script>`

func Test(t *testing.T) {
	const wantExpr = "south park"
	const wantPath = "/search/title/"
	wantQuery := url.Values{}
	wantQuery.Set("title", wantExpr)

	reached := false

	wantResult := SearchTitleResult{
		ID:          "tt0121955",
		Type:        "tvSeries",
		TrailerID:   "vi1155383321",
		Title:       "South Park",
		Image:       "imageurl.com/something.png",
		Description: "Follows the misadventures of four irreverent grade-schoolers in the quiet, dysfunctional town of South Park, Colorado.",
		Genres:      []string{"Animation", "Comedy"},
		Rating:      8.7,
		VoteCount:   399728,
		ReleaseYear: 1997,
		EndYear:     0,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() { reached = true }()

		if r.URL.Path != wantPath {
			t.Fatalf("want: '%s', got: '%s'", wantPath, r.URL.Path)
		}

		query := r.URL.Query()
		for k := range wantQuery {
			if wantQuery.Get(k) != query.Get(k) {
				t.Fatalf("want: '%s', got: '%v'", wantQuery.Get(k), query.Get(k))
			}
		}

		err := template.Must(template.New("").Parse(responseTmpl)).Execute(w, wantResult)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	c := Client{
		BaseURL: server.URL,
	}

	data, err := c.SearchTitle(wantExpr)
	if err != nil {
		t.Fatal(err)
	}

	if !reached {
		t.Fatal("server not reached")
	}

	if len(data.Results) < 1 {
		t.Fatalf("want: len >= 1, got: len %d", len(data.Results))
	}

	if !equalJSON(wantResult, data.Results[0]) {
		t.Fatalf("want:\n%+v\n\ngot:\n%v", wantResult, data.Results[0])
	}
}

func equalJSON(a, b any) bool {
	ja, err := json.Marshal(a)
	if err != nil {
		return false
	}

	jb, err := json.Marshal(b)
	if err != nil {
		return false
	}

	return bytes.Equal(ja, jb)
}
