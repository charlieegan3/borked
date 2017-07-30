package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
)

//ScanEndpoint is a handler for scan requests from the API gateway
func ScanEndpoint(w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(r.URL.Query()["url"][0])
	if err != nil {
		http.Error(w, "my own error message", http.StatusBadRequest)
		return
	}

	result := Scan(*url)
	sort.Sort(ByURL(result))
	jsonResult, _ := json.Marshal(result)
	w.Write(jsonResult)
}
