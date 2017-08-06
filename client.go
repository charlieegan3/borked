package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// LoadPage requests the document at the provided URL
// returns status code (and body if HTML)
func LoadPage(url url.URL, host string, result chan URLResult, unstarted *unstartedURLs) {
	req, err := http.NewRequest("GET", url.String(), nil)

	if err != nil {
		result <- URLResult{url, 0, err.Error()}
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:55.0) Gecko/20100101 Firefox/55.0")

	var client http.Client
	resp, err := client.Do(req)

	if err != nil {
		result <- URLResult{url, 0, err.Error()}
		return
	}

	contentType := resp.Header.Get("Content-Type")

	defer resp.Body.Close()

	var body string
	if url.Host == host && strings.Contains(contentType, "html") {
		rawBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			result <- URLResult{url, 0, err.Error()}
			return
		}
		body = string(rawBody)
	} else {
		body = ""
	}

	for _, v := range ExtractLinks(body, url) {
		unstarted.append(v)
	}

	result <- URLResult{url, resp.StatusCode, ""}
}
