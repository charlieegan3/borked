package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"time"
)

//ScanEndpoint is a handler for scan requests from the API gateway
func ScanEndpoint(w http.ResponseWriter, r *http.Request) {
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

	completed, _ := Scan(*rootURL, []url.URL{*rootURL}, 10, 10*time.Second)
	sort.Sort(ByURL(completed))
	jsonResult, _ := json.Marshal(completed)
	w.Write(jsonResult)
}
