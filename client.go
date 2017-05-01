package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

// PageResult a structure for storing a loaded page
type PageResult struct {
	StatusCode int
	Body       string
}

// LoadPage requests the document at the provided URL
// returns status code (and body if HTML)
func LoadPage(url url.URL) (PageResult, error) {
	var pr PageResult

	resp, err := http.Get(url.String())

	if err != nil {
		return pr, err
	}

	contentType := resp.Header.Get("Content-Type")

	defer resp.Body.Close()

	var body string
	if contentType == "text/html" {
		rawBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return pr, err
		}
		body = string(rawBody)
	} else {
		body = ""
	}

	return PageResult{resp.StatusCode, string(body)}, nil
}
