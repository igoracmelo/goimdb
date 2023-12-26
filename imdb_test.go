package imdb

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test(t *testing.T) {
	const wantExpr = "south park"
	const wantPath = "/find/"
	wantQuery := url.Values{}
	wantQuery.Set("q", wantExpr)
	wantQuery.Set("s", "tt")

	reached := false

	wantResult := SearchResult{
		ID:    "tt0121955",
		Title: "South Park",
		Image: "https://m.media-amazon.com/images/M/MV5BZjNhODYzZGItZWQ3Ny00ZjViLTkxMTUtM2EzN2RjYjU2OGZiXkEyXkFqcGdeQXVyMTI5MTc0OTIy._V1_.jpg",
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

		_, _ = w.Write([]byte(`
		<script id="__NEXT_DATA__" type="application/json">
		{
			"props": {
				"pageProps": {
					"titleResults": {
						"results": [
							{
								"id": "` + wantResult.ID + `",
								"imageType": "tvSeries",
								"titleNameText": "` + wantResult.Title + `",
								"titlePosterImageModel": {
									"caption": "Matt Stone and Trey Parker in South Park (1997)",
									"maxHeight": 900,
									"maxWidth": 600,
									"url": "` + wantResult.Image + `"
								},
								"titleReleaseText": "1997– ",
								"titleTypeText": "TV Series",
								"topCredits": [
									"Trey Parker",
									"Matt Stone"
								]
							}
						]
					}
				}
			}
		}
		</script>`))
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

	if data.Results[0] != wantResult {
		t.Fatalf("want:\n%+v\n\ngot:\n%v", wantResult, data.Results[0])
	}
}
