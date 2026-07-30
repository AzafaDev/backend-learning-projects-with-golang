package main

import (
	"blogging-platform-api/internal/config"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
	fmt.Println("server is running on port:", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, mux)
}
