package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `Method: %s`, r.Method)
	})

	http.HandleFunc("/{username}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"username":"%s"}`, r.PathValue("username"))
	})

	http.ListenAndServe(":8080", nil)
}
