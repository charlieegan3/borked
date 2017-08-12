package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"
)

//BuildHandler configures borked handlers with timeouts and concurrency settings
func BuildHandler(concurrency int, timeout time.Duration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		rootURLString := r.URL.Query()["root"]
		if rootURLString == nil {
			http.Error(w, "no root url", http.StatusBadRequest)
			return
		}

		rootURL, err := url.Parse(rootURLString[0])
		if err != nil {
			http.Error(w, "invalid root url", http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
		}

		var rawUrls []string
		err = json.Unmarshal(body, &rawUrls)
		if err != nil {
			http.Error(w, "failed to parse URL list", http.StatusBadRequest)
		}

		var urls []url.URL
		for _, v := range rawUrls {
			parsedURL, err := url.Parse(v)
			if err == nil {
				urls = append(urls, *parsedURL)
			}
		}

		completed, incomplete := Scan(*rootURL, urls, concurrency, timeout)

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
