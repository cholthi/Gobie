package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
)

func main() {
	if err := cgi.Serve(http.HandlerFunc(handle)); err != nil {
		fmt.Println(err)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	query := r.URL.Query()
	header.Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, r.Method)
	fmt.Fprintln(w, r.URL.String())

	for k := range query {
		fmt.Fprintln(w, k+":", query.Get(k))
	}
}
