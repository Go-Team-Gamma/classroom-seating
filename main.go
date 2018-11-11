package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
	})

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("http.listenAndServe: ", err)
	}
}
