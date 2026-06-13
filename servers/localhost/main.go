package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `Method: %s`, r.Method)
	})

	http.HandleFunc("/{username}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"username":"%s"}`, r.PathValue("username"))
	})

	http.HandleFunc("/body", func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		fmt.Fprintf(w, `%s`, b)
	})

	http.ListenAndServe(":8080", nil)
}
