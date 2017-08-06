package main

import (
	"context"
	"encoding/json"
	"errors"
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

//MarshalJSON converts a URLResult into a json string
func (u *URLResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		URL        string `json:"url"`
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
	}{
		URL:        u.URL.String(),
		StatusCode: u.StatusCode,
		Message:    u.Message,
	})
}

// ByURL can be used to sort a list by the string URL value
// used in tests to get a consistent order of URL results
type ByURL []URLResult

func (a ByURL) Len() int           { return len(a) }
func (a ByURL) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByURL) Less(i, j int) bool { return a[i].URL.String() < a[j].URL.String() }

//UnstartedURL is a URL yet to be scanned
type UnstartedURL struct {
	URL url.URL
}

//MarshalJSON converts a URLResult into a json string
func (u *UnstartedURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		URL string `json:"url"`
	}{
		URL: u.URL.String(),
	})
}

type unstartedURLs struct {
	Elements []UnstartedURL
	mux      sync.Mutex
}

func (v *unstartedURLs) append(item UnstartedURL) {
	if !v.contains(item.URL) {
		v.mux.Lock()
		v.Elements = append(v.Elements, item)
		v.mux.Unlock()
	}
}

func (v *unstartedURLs) pop() (UnstartedURL, error) {
	v.mux.Lock()
	defer v.mux.Unlock()

	length := len(v.Elements)
	if length == 0 {
		return UnstartedURL{}, errors.New("Empty")
	}

	var element UnstartedURL
	element, v.Elements = v.Elements[length-1], v.Elements[:length-1]
	return element, nil
}

func (v *unstartedURLs) contains(item url.URL) bool {
	found := false
	for _, e := range v.Elements {
		if e.URL == item {
			return true
		}
	}
	return found
}

type completedURLs struct {
	Elements []URLResult
	mux      sync.Mutex
}

func (v *completedURLs) append(item URLResult) {
	if !v.contains(item.URL) {
		v.mux.Lock()
		v.Elements = append(v.Elements, item)
		v.mux.Unlock()
	}
}

func (v *completedURLs) contains(item url.URL) bool {
	found := false
	for _, e := range v.Elements {
		if e.URL == item {
			return true
		}
	}
	return found
}

type idleCounter struct {
	Count int
	Total int
	mux   sync.Mutex
}

func (c *idleCounter) Inc() {
	c.mux.Lock()
	if c.Count < c.Total {
		c.Count++
	}
	c.mux.Unlock()
}

func (c *idleCounter) Dec() {
	c.mux.Lock()
	if c.Count > 0 {
		c.Count--
	}
	c.mux.Unlock()
}

func (c *idleCounter) All() bool {
	return c.Count == c.Total
}

// Scan for broken links starting from a given page
func Scan(root url.URL, urls []url.URL, concurrency int, timeout time.Duration) ([]URLResult, []UnstartedURL) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var unstarted unstartedURLs
	for _, v := range urls {
		unstarted.append(UnstartedURL{URL: v})
	}
	var completed completedURLs
	idle := idleCounter{Total: concurrency}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go work(ctx, &wg, &idle, root.Host, &unstarted, &completed)
	}
	wg.Wait()

	return completed.Elements, unstarted.Elements
}

func work(ctx context.Context, wg *sync.WaitGroup, idle *idleCounter, host string, unstarted *unstartedURLs, completed *completedURLs) error {
	defer wg.Done()

	for {
		unstartedURL, err := unstarted.pop()
		if !completed.contains(unstartedURL.URL) {
			if err == nil {
				idle.Dec()
				result := make(chan URLResult)
				go LoadPage(unstartedURL.URL, host, result, unstarted)
				select {
				case urlResult := <-result:
					completed.append(urlResult)
				case <-ctx.Done():
					unstarted.append(unstartedURL)
					return ctx.Err()
				}
			} else {
				idle.Inc()
				if idle.All() {
					return nil
				}
			}
		}
	}
}
