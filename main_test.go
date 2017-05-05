package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCrawling(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			pageContent := `
                <html>
                    <a href="/page2">Page 2</a>
                <html>`

			fmt.Fprintln(w, pageContent)
		} else if r.URL.Path == "/page2" {
			pageContent := `
                <html>
                    <a href="/page3">Page 3</a>
                <html>`

			fmt.Fprintln(w, pageContent)
		} else if r.URL.Path == "/page3" {
			pageContent := "<html></html>"

			fmt.Fprintln(w, pageContent)
		}
	}))

	startPage, _ := url.Parse(localServer.URL + "/")

	urls := Scan(*startPage, false)

	if len(urls) != 3 {
		t.Error("Expected 3 links, got", len(urls))
	}

	var stringUrls []string
	for _, v := range urls {
		stringUrls = append(stringUrls, v.String())
	}

	expected := []string{localServer.URL + "/", localServer.URL + "/page2", localServer.URL + "/page3"}

	for i, v := range expected {
		if v != stringUrls[i] {
			t.Error("Unexpected result, got", v, stringUrls[i])
		}
	}
}
