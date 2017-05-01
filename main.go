package main

import (
	"fmt"
	"net/url"
)

func main() {
	startingPage, err := url.Parse("https://charlieegan3.com/")
	if err != nil {
		fmt.Println(err)
		return
	}

	pageResult, err := LoadPage(*startingPage)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(pageResult.StatusCode)

	links := ExtractLinks(pageResult.Body, *startingPage)
	for _, l := range links {
		fmt.Println(l.String())
	}
}
