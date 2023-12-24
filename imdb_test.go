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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != wantPath {
			t.Errorf("want: '%s', got: '%s'", wantPath, r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("q") != wantExpr {
			t.Errorf("want: '%s', got: '%s'", wantExpr, query.Get("q"))
		}

		reached = true
	}))
	defer server.Close()

	c := Client{
		BaseURL: server.URL,
	}

	_, err := c.SearchTitle(wantExpr)
	if err != nil {
		t.Error(err)
	}

	if !reached {
		t.Error("server not reached")
	}
}
