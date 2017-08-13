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
		w.Header().Set("Access-Control-Allow-Origin", "https://borked.charlieegan3.com")
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
		if rootURL.Scheme == "" {
			rootURL, err = url.Parse("http://" + rootURL.String())
			if err != nil {
				http.Error(w, "invalid root url", http.StatusBadRequest)
				return
			}
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
		}

		var task struct {
			VisitedURLs    []string `json:"visited"`
			IncompleteURLs []string `json:"incomplete"`
		}
		err = json.Unmarshal(body, &task)
		if err != nil {
			http.Error(w, "failed to parse URL list", http.StatusBadRequest)
		}

		var incompleteURLs []url.URL
		for _, v := range task.IncompleteURLs {
			parsedURL, err := url.Parse(v)
			if err == nil {
				incompleteURLs = append(incompleteURLs, *parsedURL)
			}
		}
		if len(incompleteURLs) == 0 {
			incompleteURLs = append(incompleteURLs, *rootURL)
		}

		var visitedURLs []url.URL
		for _, v := range task.VisitedURLs {
			parsedURL, err := url.Parse(v)
			if err == nil {
				visitedURLs = append(visitedURLs, *parsedURL)
			}
		}

		completed, incomplete := Scan(*rootURL, incompleteURLs, visitedURLs, concurrency, timeout)

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
