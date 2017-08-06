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
		urlParams := r.URL.Query()["url"]
		if len(urlParams) == 0 {
			http.Error(w, "no url", http.StatusBadRequest)
			return
		}

		rootURL, err := url.Parse(urlParams[0])
		if err != nil {
			http.Error(w, "url parse failed", http.StatusBadRequest)
			return
		}

		completed, incomplete := Scan(*rootURL, []url.URL{*rootURL}, concurrency, timeout)

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
