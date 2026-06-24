package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		v := make(map[string]any)
		v["method"] = r.Method
		v["remote addr"] = r.RemoteAddr
		v["trailer"] = r.Trailer

		json.NewEncoder(w).Encode(&v)
	})

	http.HandleFunc("GET /{username}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"username":"%s"}`, r.PathValue("username"))
	})

	http.HandleFunc("POST /body", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", r.Header.Get("Content-Type"))

		b, _ := io.ReadAll(r.Body)
		fmt.Fprintf(w, `%s`, b)
	})

	http.HandleFunc("GET /wait", func(w http.ResponseWriter, r *http.Request) {
		s, err := strconv.Atoi(r.URL.Query().Get("s"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dur := time.Duration(s) * time.Second
		time.Sleep(dur)
	})

	http.ListenAndServe(":8080", nil)
}
