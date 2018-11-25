package main

import (
	"net/http"
	"time"

	"github.com/gobuffalo/packr"
)

func main() {
	box := packr.NewBox("./web")

	http.Handle("/process", http.HandlerFunc(BuildHandler(10, 10*time.Second)))
	http.Handle("/", http.FileServer(box))
	http.ListenAndServe(":3000", nil)
}
