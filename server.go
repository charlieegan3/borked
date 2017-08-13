package main

import (
	"net"
	"net/http"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	http.Serve(ln, http.HandlerFunc(BuildHandler(10, 10*time.Second)))
}
