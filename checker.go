package main

import (
	"net/url"
	"sync"
	"time"
)

// URLResult represents a checked page and the result of making that check
type URLResult struct {
	URL        url.URL
	StatusCode int
	Message    string
}

// ByURL can be used to sort a list by the string URL value
// used in tests to get a consistent order of URL results
type ByURL []URLResult

func (a ByURL) Len() int           { return len(a) }
func (a ByURL) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByURL) Less(i, j int) bool { return a[i].URL.String() < a[j].URL.String() }

// Scan for broken links starting from a given page
func Scan(url url.URL) []URLResult {
	var result scannedURLs

	var cc connectionCounter
	var wg sync.WaitGroup

	wg.Add(1)
	go checkURL(url, url, url, &result, &wg, &cc)
	wg.Wait()

	return result.URLs
}

type connectionCounter struct {
	count int
	mux   sync.Mutex
}

func (c *connectionCounter) Inc() {
	c.mux.Lock()
	c.count++
	c.mux.Unlock()
}

func (c *connectionCounter) Dec() {
	c.mux.Lock()
	c.count--
	c.mux.Unlock()
}

func (c *connectionCounter) Count() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.count
}

func (c *connectionCounter) Wait() {
	for {
		if c.Count() < 500 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// ScanResult reprents a list of URLs checked and the result for each
type scannedURLs struct {
	URLs []URLResult
	mux  sync.Mutex
}

func (v *scannedURLs) append(link URLResult) {
	v.mux.Lock()
	v.URLs = append(v.URLs, link)
	v.mux.Unlock()
}

func (v *scannedURLs) list() []URLResult {
	v.mux.Lock()
	defer v.mux.Unlock()
	return v.URLs
}

func (v *scannedURLs) contains(url url.URL) bool {
	for _, u := range v.list() {
		if u.URL == url {
			return true
		}
	}
	return false
}

func checkURL(url url.URL, source url.URL, root url.URL, results *scannedURLs,
	wg *sync.WaitGroup, cc *connectionCounter) {

	defer wg.Done()

	if results.contains(url) {
		return
	}

	cc.Wait()

	cc.Inc()
	pageResult, err := LoadPage(url, root.Host)
	cc.Dec()

	if results.contains(url) {
		return
	}

	if err != nil {
		results.append(URLResult{url, pageResult.StatusCode, err.Error()})
	} else {
		results.append(URLResult{url, pageResult.StatusCode, ""})
	}

	if url.Host != root.Host {
		return
	}

	links := ExtractLinks(pageResult.Body, url)
	for _, l := range links {
		wg.Add(1)
		go checkURL(l, url, root, results, wg, cc)
	}
}
