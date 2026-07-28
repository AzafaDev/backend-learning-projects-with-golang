package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "GET /")
	})
	mux.HandleFunc("GET /articles/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		fmt.Fprintf(w, "GET /articles/%s\n", slug)
	})
	mux.HandleFunc("GET /admin", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "GET /admin")
	})

	fmt.Println("server is running on port: 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("error:", err)
	}

}
