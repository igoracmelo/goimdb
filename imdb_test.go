package imdb

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) {
	const wantExpr = "south park"
	const wantPath = "/find/"

	reached := false

	wantResult := SearchResult{
		ID:    "tt0121955",
		Title: "South Park",
		Image: "https://m.media-amazon.com/images/M/MV5BZjNhODYzZGItZWQ3Ny00ZjViLTkxMTUtM2EzN2RjYjU2OGZiXkEyXkFqcGdeQXVyMTI5MTc0OTIy._V1_.jpg",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached = true

		if r.URL.Path != wantPath {
			t.Fatalf("want: '%s', got: '%s'", wantPath, r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("q") != wantExpr {
			t.Fatalf("want: '%s', got: '%s'", wantExpr, query.Get("q"))
		}

		_, _ = w.Write([]byte(`{
			"props": {
				"pageProps": {
					"titleResults": {
						"results": [
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
				}
			}
		}`))
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

	if len(data.Results) != 1 {
		t.Fatalf("want: len 1, got: len %d", len(data.Results))
	}

	if data.Results[0] != wantResult {
		t.Fatalf("want:\n%+v\n\ngot:\n%v", wantResult, data.Results[0])
	}
}
