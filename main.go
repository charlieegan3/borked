package main

import (
	"fmt"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide a starting page as an argument")
		return
	}

	start, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	Scan(*start)
}
