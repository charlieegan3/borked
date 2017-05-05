package main

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// Scan for broken links starting from a given page
func Scan(url url.URL, showSuccessful bool) []url.URL {
	var visited visitedList
	var cc connectionCounter
	var wg sync.WaitGroup

	wg.Add(1)
	go checkURL(url, url, url, &visited, &wg, &cc, showSuccessful)
	wg.Wait()

	return visited.URLs
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

type visitedList struct {
	URLs []url.URL
	mux  sync.Mutex
}

func (v *visitedList) Append(link url.URL) {
	v.mux.Lock()
	v.URLs = append(v.URLs, link)
	v.mux.Unlock()
}

func (v *visitedList) List() []url.URL {
	v.mux.Lock()
	defer v.mux.Unlock()
	return v.URLs
}

func (v *visitedList) Contains(url url.URL) bool {
	for _, u := range v.List() {
		if u == url {
			return true
		}
	}
	return false
}

func checkURL(url url.URL, source url.URL, root url.URL, visited *visitedList,
	wg *sync.WaitGroup, cc *connectionCounter, showSuccessful bool) {

	defer wg.Done()

	if visited.Contains(url) {
		return
	}

	for {
		if cc.Count() < 500 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	cc.Inc()
	pageResult, err := LoadPage(url, root.Host)
	cc.Dec()

	if visited.Contains(url) {
		return
	}
	visited.Append(url)

	if err != nil {
		fmt.Println(source.String() + "\n  " + err.Error())
		return
	}

	if showSuccessful == true || pageResult.StatusCode != 200 {
		fmt.Println(source.String() + "\n  " + strconv.Itoa(pageResult.StatusCode) + " - " + url.String())
	}

	if url.Host != root.Host {
		return
	}

	links := ExtractLinks(pageResult.Body, url)
	for _, l := range links {
		wg.Add(1)
		go checkURL(l, url, root, visited, wg, cc, showSuccessful)
	}
}
