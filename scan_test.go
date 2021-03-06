package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"testing"
	"time"
)

func TestCrawling(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			pageContent := `
                <html>
                    <a href="/page2">Page 2</a>
                <html>`

			fmt.Fprintln(w, pageContent)
		} else if r.URL.Path == "/page2" {
			pageContent := `
                <html>
                    <a href="/page3">Page 3</a>
                    <a href="/404">Borked</a>
                    <a href="http://nowhere.com">Borked</a>
                <html>`

			fmt.Fprintln(w, pageContent)
		} else if r.URL.Path == "/page3" {
			pageContent := "<html></html>"

			fmt.Fprintln(w, pageContent)
		} else if r.URL.Path == "/404" {
			http.NotFound(w, r)
		}
	}))

	startPage, _ := url.Parse(localServer.URL + "/")

	completed, _ := Scan(*startPage, []JSONUnstartedURL{JSONUnstartedURL{startPage.String(), ""}}, []url.URL{}, 1, time.Second)
	sort.Sort(ByURL(completed))

	if len(completed) != 5 {
		t.Error("Expected 5 links, got", len(completed))
	}

	var urls []string
	var statuses []int
	var messages []string

	for _, u := range completed {
		urls = append(urls, u.URL.String())
		statuses = append(statuses, u.StatusCode)
		messages = append(messages, u.Message)
	}

	expectedUrls := []string{
		localServer.URL + "/",
		localServer.URL + "/404",
		localServer.URL + "/page2",
		localServer.URL + "/page3",
		"http://nowhere.com",
	}
	for i, v := range urls {
		if v != expectedUrls[i] {
			t.Error("Unexpected URL: ", v)
		}
	}

	expectedStatuses := []int{200, 404, 200, 200, 0}
	for i, v := range statuses {
		if v != expectedStatuses[i] {
			t.Error("Unexpected Status: ", v)
		}
	}

	expectedMessages := []string{"", "", "", "", "Head http://nowhere.com: dial tcp: lookup nowhere.com: no such host"}
	for i, v := range messages {
		if v != expectedMessages[i] {
			t.Error("Unexpected Message: ", v)
		}
	}
}

func TestCrawlingTimeout(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			pageContent := `
                <html>
                    <a href="/page2">Page 2</a>
                <html>`

			fmt.Fprintln(w, pageContent)
		} else if r.URL.Path == "/page2" {
			time.Sleep(5 * time.Millisecond)
			pageContent := `
                <html>
                    <a href="/page3">Page 3</a>
                    <a href="/404">Borked</a>
                    <a href="http://nowhere.com">Borked</a>
                <html>`

			fmt.Fprintln(w, pageContent)
		}
	}))

	startPage, _ := url.Parse(localServer.URL + "/")

	completed, incomplete := Scan(*startPage, []JSONUnstartedURL{JSONUnstartedURL{startPage.String(), ""}}, []url.URL{}, 1, 3*time.Millisecond)

	if len(completed) != 1 {
		t.Error("Expected 1 completed link, got", len(completed))
	}

	if len(incomplete) != 1 {
		t.Error("Expected 1 incomplete link, got", len(incomplete))
	}
}

func TestCrawlingCyclic(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			pageContent := `
                <html>
                    <a href="/page2">Page 2</a>
                <html>`

			fmt.Fprintln(w, pageContent)
		} else if r.URL.Path == "/page2" {
			pageContent := `
                <html>
                    <a href="/">Home again</a>
                <html>`

			fmt.Fprintln(w, pageContent)
		}
	}))

	startPage, _ := url.Parse(localServer.URL + "/")

	completed, _ := Scan(*startPage, []JSONUnstartedURL{JSONUnstartedURL{startPage.String(), ""}}, []url.URL{}, 1, time.Second)
	sort.Sort(ByURL(completed))

	if len(completed) != 2 {
		t.Error("Expected 2 links, got", len(completed))
	}
}
