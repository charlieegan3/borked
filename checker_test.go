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

	result := Scan(*startPage, "10s")
	sort.Sort(ByURL(result))

	if len(result) != 5 {
		t.Error("Expected 5 links, got", len(result))
	}

	var urls []string
	var statuses []int
	var messages []string

	for _, u := range result {
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

	expectedMessages := []string{"", "", "", "", "Get http://nowhere.com: dial tcp: lookup nowhere.com: no such host"}
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
                    <a href="/slow">Page 2</a>
                <html>`

			fmt.Fprintln(w, pageContent)
		} else if r.URL.Path == "/slow" {
			time.Sleep(5 * time.Millisecond) // this pushes the task past the soft timeout
			pageContent := `
                <html>
                    <a href="/saved_for_next_time">Page 3</a>
                <html>`

			fmt.Fprintln(w, pageContent)
		} else if r.URL.Path == "/saved_for_next_time" {
			pageContent := `
                <html>
                    <a href="/missed">Page 4</a>
                <html>`

			fmt.Fprintln(w, pageContent)
		} else if r.URL.Path == "/missed" {
			pageContent := `
                <html>
                    missed until the next attempt
                <html>`

			fmt.Fprintln(w, pageContent)
		}
	}))

	startPage, _ := url.Parse(localServer.URL + "/")

	result := Scan(*startPage, "3ms") // make sure after completing the second page the time has run out

	if len(result) != 3 {
		t.Error("Expected 3 links, got", len(result))
	}

	if result[len(result)-1].StatusCode != -1 {
		fmt.Println(result[len(result)-1].StatusCode)
		t.Error("Expected final page to be incomplete")
	}
}
