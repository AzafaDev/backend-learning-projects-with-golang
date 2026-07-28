package article

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Article struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Slug        string    `json:"slug"`
	PublishedAt time.Time `json:"published_at"`
}

func LoadArticles(dir string) ([]Article, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var articles []Article
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		fileName := filepath.Join(dir, entry.Name())
		articleByte, err := os.ReadFile(fileName)
		if err != nil {
			continue
		}
		var article Article
		slug := strings.TrimSuffix(entry.Name(), ".json")
		article.Slug = slug
		if err := json.Unmarshal(articleByte, &article); err != nil {
			continue
		}
		articles = append(articles, article)
	}
	return articles, nil
}
