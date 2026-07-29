package main

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"personal-blog/article"
	"time"
)

const adminUsername = "Akmal"
const adminPassword = "AkmalGantengBanget"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		articles, err := article.LoadArticles("data")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, articles); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	mux.HandleFunc("GET /articles/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		articles, err := article.LoadArticles("data")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		found := false
		for _, v := range articles {
			if slug == v.Slug {
				tmpl, err := template.ParseFiles("templates/article.html")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if err := tmpl.Execute(w, v); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
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
	mux.HandleFunc("GET /admin", requireAuth(func(w http.ResponseWriter, r *http.Request) {
		articles, err := article.LoadArticles("data")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		templ, err := template.ParseFiles("templates/admin.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := templ.Execute(w, articles); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))

	mux.HandleFunc("GET /admin/articles/new", requireAuth(func(w http.ResponseWriter, r *http.Request) {
		templ, err := template.ParseFiles("templates/new_article.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := templ.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))

	mux.HandleFunc("POST /admin/articles/new", requireAuth(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		title := r.FormValue("title")
		content := r.FormValue("content")
		publishedAt, err := time.Parse("2006-01-02", r.FormValue("published_at"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		slug := article.Slugify(title)
		a := article.Article{
			Title:       title,
			Content:     content,
			PublishedAt: publishedAt,
			Slug:        slug,
		}
		entries, err := os.ReadDir("data")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			} else if e.Name() == a.Slug+".json" {
				http.Error(w, "error: slug is already exists", http.StatusConflict)
				return
			} else {
				continue
			}
		}
		if err := article.SaveArticle("data", a); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}))

	mux.HandleFunc("GET /admin/articles/{slug}/edit", requireAuth(func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		templ, err := template.ParseFiles("templates/edit_article.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		articles, err := article.LoadArticles("data")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		found := false
		var art article.Article
		for _, a := range articles {
			if a.Slug == slug {
				found = true
				art = a
				break
			} else {
				continue
			}
		}
		if !found {
			http.Error(w, "error: not found", http.StatusNotFound)
			return
		}
		if err := templ.Execute(w, art); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))

	mux.HandleFunc("POST /admin/articles/{slug}/edit", requireAuth(func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		title := r.FormValue("title")
		content := r.FormValue("content")
		publishedAt, err := time.Parse("2006-01-02", r.FormValue("published_at"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		a := article.Article{
			Title:       title,
			Content:     content,
			PublishedAt: publishedAt,
			Slug:        slug,
		}
		articles, err := article.LoadArticles("data")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		found := false
		for _, a := range articles {
			if a.Slug == slug {
				found = true
				break
			} else {
				continue
			}
		}
		if !found {
			http.Error(w, "error: not found", http.StatusNotFound)
			return
		}
		if err := article.SaveArticle("data", a); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}))

	mux.HandleFunc("POST /admin/articles/{slug}/delete", requireAuth(func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		if err := article.DeleteArticle("data", slug); err != nil {
			if errors.Is(err, article.ErrNotFound) {
				http.Error(w, "article not found", http.StatusNotFound)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}))

	fmt.Println("server is running on port: 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("error:", err)
	}

}

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		same := subtle.ConstantTimeCompare([]byte(password), []byte(adminPassword))
		if !ok || username != adminUsername || same != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "unauthorized")
			return
		}
		next(w, r)
	}
}
