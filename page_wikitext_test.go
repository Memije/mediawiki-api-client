package mediawiki

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const pageWikitextTestURL = "/wikitext"
const pageWikitextTestTitle = "test"
const pageWikitextTestRevision = 1
const pageWikitextTestContent = "hello world"

func createPageWikitextServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(pageWikitextTestURL+"/"+pageWikitextTestTitle, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{ "continue": { "rvcontinue": "20020106234815|248709", "continue": "||" }, "query": { "pages": [ { "pageid": 9228, "ns": 0, "title": "%s", "revisions": [ { "slots": { "main": { "contentmodel": "wikitext", "contentformat": "text/x-wiki", "content": "%s" } } } ] } ] } }`, pageWikitextTestTitle, pageWikitextTestContent)))
	})

	return router
}

func TestPageWikitext(t *testing.T) {
	srv := httptest.NewServer(createPageWikitextServer())
	defer srv.Close()

	client := NewClient(srv.URL)
	client.options.PageWikitextURL = pageWikitextTestURL + "/%s?test=1"

	wikitext, status, err := client.PageWikitext(context.Background(), pageWikitextTestTitle, pageWikitextTestRevision)

	if err != nil {
		t.Error(err)
	}

	if status != http.StatusOK {
		t.Fatal("page wikitext response error")
	}

	if string(wikitext) != pageWikitextTestContent {
		t.Error("page wikitext content is wrong")
	}
}
