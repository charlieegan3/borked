package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"time"
)

//BuildHandler configures borked handlers with timeouts and concurrency settings
func BuildHandler(concurrency int, timeout time.Duration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		urls := r.URL.Query()["url"]
		if len(urls) == 0 {
			http.Error(w, "no url", http.StatusBadRequest)
			return
		}

		var startingURLs []url.URL
		for _, v := range urls {
			parsedURL, err := url.Parse(v)
			if err != nil {
				http.Error(w, "url parse failed", http.StatusBadRequest)
				return
			}
			startingURLs = append(startingURLs, *parsedURL)
		}

		completed, incomplete := Scan(startingURLs[0], startingURLs, concurrency, timeout)

		sort.Sort(ByURL(completed))

		responseData := struct {
			Completed  []URLResult    `json:"completed"`
			Incomplete []UnstartedURL `json:"incomplete"`
		}{
			completed,
			incomplete,
		}

		jsonResult, _ := json.Marshal(responseData)
		w.Write(jsonResult)
	}
}
