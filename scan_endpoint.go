package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
)

//ScanEndpoint is a handler for scan requests from the API gateway
func ScanEndpoint(w http.ResponseWriter, r *http.Request) {
	urlParams := r.URL.Query()["url"]
	if len(urlParams) == 0 {
		http.Error(w, "no url", http.StatusBadRequest)
		return
	}

	url, err := url.Parse(urlParams[0])
	if err != nil {
		http.Error(w, "url parse failed", http.StatusBadRequest)
		return
	}

	result := Scan(*url, "10s")
	sort.Sort(ByURL(result))
	jsonResult, _ := json.Marshal(result)
	w.Write(jsonResult)
}
