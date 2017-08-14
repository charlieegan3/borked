package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"
)

// JSONUnstartedURL struct to parse incomplete URLs listed in the request
type JSONUnstartedURL struct {
	URL    string
	Source string
}

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
			return
		}

		var task struct {
			VisitedURLs    []string           `json:"visited"`
			IncompleteURLs []JSONUnstartedURL `json:"incomplete"`
		}
		err = json.Unmarshal(body, &task)
		if err != nil {
			http.Error(w, "failed to parse URL list", http.StatusBadRequest)
			return
		}

		if len(task.IncompleteURLs) == 0 {
			task.IncompleteURLs = append(task.IncompleteURLs, JSONUnstartedURL{rootURL.String(), rootURL.String()})
		}

		var visitedURLs []url.URL
		for _, v := range task.VisitedURLs {
			parsedURL, err := url.Parse(v)
			if err == nil {
				visitedURLs = append(visitedURLs, *parsedURL)
			}
		}

		completed, incomplete := Scan(*rootURL, task.IncompleteURLs, visitedURLs, concurrency, timeout)

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

		LogJob(*rootURL, len(completed), r.RemoteAddr, r.Header.Get("User-Agent"))
	}
}
