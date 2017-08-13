package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// LoadPage requests the document at the provided URL
// returns status code (and body if HTML)
func LoadPage(link UnstartedURL, host string, result chan URLResult, unstarted *unstartedURLs) {
	req, err := http.NewRequest("GET", link.URL.String(), nil)

	if err != nil {
		result <- URLResult{link.URL, link.Source, 0, err.Error()}
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:55.0) Gecko/20100101 Firefox/55.0")

	var client http.Client
	resp, err := client.Do(req)

	if err != nil {
		result <- URLResult{link.URL, link.Source, 0, err.Error()}
		return
	}

	contentType := resp.Header.Get("Content-Type")

	defer resp.Body.Close()

	var body string
	if link.URL.Host == host && strings.Contains(contentType, "html") {
		rawBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			result <- URLResult{link.URL, link.Source, 0, err.Error()}
			return
		}
		body = string(rawBody)
	} else {
		body = ""
	}

	for _, v := range ExtractLinks(body, link.URL) {
		unstarted.append(UnstartedURL{URL: v, Source: link.URL})
	}

	result <- URLResult{link.URL, link.Source, resp.StatusCode, ""}
}
