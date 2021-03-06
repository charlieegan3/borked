package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestScanEndpoint(t *testing.T) {
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

	lsURL := localServer.URL
	req, err := http.NewRequest("POST", fmt.Sprintf("/?root=%s", lsURL), bytes.NewBuffer([]byte(`{ "visited": [], "incomplete": []}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(BuildHandler(10, 10*time.Second))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: \ngot\n%v\nwant\n%v",
			status, http.StatusOK)
	}

	expected := fmt.Sprintf(`{"completed":[{"url":"%v","source":"%v","status_code":200,"message":""},{"url":"%v/404","source":"%v/page2","status_code":404,"message":""},{"url":"%v/page2","source":"%v","status_code":200,"message":""},{"url":"%v/page3","source":"%v/page2","status_code":200,"message":""},{"url":"http://nowhere.com","source":"%v/page2","status_code":0,"message":"Head http://nowhere.com: dial tcp: lookup nowhere.com: no such host"}],"incomplete":[]}`, lsURL, lsURL, lsURL, lsURL, lsURL, lsURL, lsURL, lsURL, lsURL)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: \ngot\n%v\nwant\n%v",
			rr.Body.String(), expected)
	}
}

func TestScanEndpointIncomplete(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			pageContent := `
            <html>
            <a href="/page2">Page 2</a>
            <html>`

			fmt.Fprintln(w, pageContent)
		} else if r.URL.Path == "/page2" {
			time.Sleep(time.Second)
			pageContent := `
            <html>
            <a href="/page3">Page 3</a>
            <html>`

			fmt.Fprintln(w, pageContent)
		}
	}))

	lsURL := localServer.URL
	req, err := http.NewRequest("POST", fmt.Sprintf("/?root=%s", lsURL), bytes.NewBuffer([]byte(fmt.Sprintf(`{ "visited": [], "incomplete": []}`))))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(BuildHandler(1, 5*time.Millisecond))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := fmt.Sprintf(
		`{"completed":[{"url":"%v","source":"%v","status_code":200,"message":""}],"incomplete":[{"url":"%v/page2","source":"%v"}]}`, lsURL, lsURL, lsURL, lsURL)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestScanEndpointMultipleStartingUrls(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			pageContent := `
            <html>
            <html>`

			fmt.Fprintln(w, pageContent)
		} else if r.URL.Path == "/page2" {
			pageContent := `
            <html>
            <a href="/page3">Page 3</a>
            <html>`

			fmt.Fprintln(w, pageContent)
		} else if r.URL.Path == "/page3" {
			pageContent := `
            <html>
            <a href="/">Home</a>
            <html>`

			fmt.Fprintln(w, pageContent)
		}
	}))

	lsURL := localServer.URL
	req, err := http.NewRequest("POST", fmt.Sprintf("/?root=%s", lsURL), bytes.NewBuffer([]byte(fmt.Sprintf(`{"visited":[],"incomplete":[{"url":"%v/page2","source":""},{"url":"%v/page3","source":""}]}`, lsURL, lsURL))))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(BuildHandler(2, 2*time.Millisecond))

	handler.ServeHTTP(rr, req)

	expected := fmt.Sprintf(
		`{"completed":[{"url":"%v/","source":"%v/page3","status_code":200,"message":""},{"url":"%v/page2","source":"","status_code":200,"message":""},{"url":"%v/page3","source":"","status_code":200,"message":""}],"incomplete":[]}`, lsURL, lsURL, lsURL, lsURL)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: \ngot \n%v\nwant\n%v", rr.Body.String(), expected)
	}
}

func TestScanEndpointWithWorkInProgress(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/page1" {
			pageContent := `
            <html>
            <a href="/page2">Page 2</a>
            <html>`

			fmt.Fprintln(w, pageContent)
		} else if r.URL.Path == "/page2" {
			pageContent := "<html>Never visited<html>"
			fmt.Fprintln(w, pageContent)
		}
	}))

	lsURL := localServer.URL
	req, err := http.NewRequest("POST", fmt.Sprintf("/?root=%s/page1", lsURL), bytes.NewBuffer([]byte(fmt.Sprintf(`{ "visited": ["%v/page2"], "incomplete": [{"url":"%v/page1","source":""}]}`, lsURL, lsURL))))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(BuildHandler(2, 5*time.Millisecond))

	handler.ServeHTTP(rr, req)

	expected := fmt.Sprintf(`{"completed":[{"url":"%v/page1","source":"","status_code":200,"message":""}],"incomplete":[]}`, lsURL)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestScanEndpointNoUrl(t *testing.T) {
	req, err := http.NewRequest("POST", fmt.Sprintf("/?no_params"), bytes.NewBuffer([]byte("")))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(BuildHandler(10, 10*time.Second))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
