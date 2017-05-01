package main

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"
)

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
	list []url.URL
	mux  sync.Mutex
}

func (v *visitedList) Append(link url.URL) {
	v.mux.Lock()
	v.list = append(v.list, link)
	v.mux.Unlock()
}

func (v *visitedList) List() []url.URL {
	v.mux.Lock()
	defer v.mux.Unlock()
	return v.list
}

func (v *visitedList) Contains(url url.URL) bool {
	for _, u := range v.List() {
		if u == url {
			return true
		}
	}
	return false
}

func checkURL(url url.URL, root url.URL, visited *visitedList, wg *sync.WaitGroup, cc *connectionCounter) {
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

	visited.Append(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(strconv.Itoa(pageResult.StatusCode) + " - " + url.String())

	if url.Host != root.Host {
		return
	}

	links := ExtractLinks(pageResult.Body, url)
	for _, l := range links {
		wg.Add(1)
		go checkURL(l, root, visited, wg, cc)
	}
}

func main() {
	start, err := url.Parse("https://charlieegan3.com/")
	if err != nil {
		fmt.Println(err)
		return
	}
	var visited visitedList
	var cc connectionCounter
	var wg sync.WaitGroup

	wg.Add(1)
	go checkURL(*start, *start, &visited, &wg, &cc)

	wg.Wait()
}
