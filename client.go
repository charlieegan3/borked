package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// PageResult a structure for storing a loaded page
type PageResult struct {
	StatusCode int
	Body       string
}

// LoadPage requests the document at the provided URL
// returns status code (and body if HTML)
func LoadPage(url url.URL, host string, timeOut string) (PageResult, error) {
	var pr PageResult

	req, err := http.NewRequest("GET", url.String(), nil)

	if err != nil {
		return pr, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:55.0) Gecko/20100101 Firefox/55.0")

	timeOutDuration, err := time.ParseDuration(timeOut)
	if err != nil {
		return pr, err
	}

	client := http.Client{
		Timeout: timeOutDuration,
	}
	resp, err := client.Do(req)

	if err != nil {
		return pr, err
	}

	contentType := resp.Header.Get("Content-Type")

	defer resp.Body.Close()

	var body string
	if url.Host == host && strings.Contains(contentType, "html") {
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
