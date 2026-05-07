package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.Body.Close()
		_, err := w.Write([]byte("Hello world"))
		if err != nil {
			fmt.Println(err.Error())
		}
	})

	http.ListenAndServe(":8080", nil)
}
