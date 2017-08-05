package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestClientTimeout(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			time.Sleep(5 * time.Millisecond)
		}
	}))

	startPage, _ := url.Parse(localServer.URL + "/")

	_, err := LoadPage(*startPage, string(localServer.URL), "5ms")
	if !strings.Contains(err.Error(), "Timeout") {
		t.Error("Timeout error expected")
	}
}
