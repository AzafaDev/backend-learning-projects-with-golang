package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"personal-blog/article"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		articles, err := article.LoadArticles("data")
		if err != nil {
			fmt.Fprintln(w, "error:", err)
			return
		}
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			fmt.Fprintln(w, "error:", err)
			return
		}
		if err := tmpl.Execute(w, articles); err != nil {
			fmt.Fprintln(w, "error:", err)
			return
		}
	})
	mux.HandleFunc("GET /articles/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		articles, err := article.LoadArticles("data")
		if err != nil {
			fmt.Fprintln(w, "error:", err)
			return
		}
		found := false
		for _, v := range articles {
			if slug == v.Slug {
				tmpl, err := template.ParseFiles("templates/article.html")
				if err != nil {
					fmt.Fprintln(w, "error:", err)
					return
				}
				if err := tmpl.Execute(w, v); err != nil {
					fmt.Fprintln(w, "error:", err)
					return
				}
				found = true
				break
			} else {
				continue
			}
		}
		if found == false {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "not found 404")
		}
	})
	mux.HandleFunc("GET /admin", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "GET /admin")
	})

	fmt.Println("server is running on port: 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("error:", err)
	}

}
