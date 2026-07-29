package article

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Article struct {
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

func SaveArticle(dir string, a Article) error {
	ab, err := json.MarshalIndent(a, "", " ")
	if err != nil {
		return err
	}

	fileName := filepath.Join(dir, a.Slug)
	if err := os.WriteFile(fileName+".json", ab, 0644); err != nil {
		return err
	}
	return nil
}

var ErrNotFound = errors.New("article not found")

func DeleteArticle(dir string, slug string) error {
	fileName := filepath.Join(dir, slug)
	_, err := os.Stat(fileName + ".json")
	if err != nil {
		return ErrNotFound
	}
	if err := os.Remove(fileName + ".json"); err != nil {
		return err
	}
	return nil
}

func Slugify(title string) string {
	s := strings.ToLower(title)

	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ' || r == '-':
			b.WriteRune('-')
		}
	}

	slug := b.String()

	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	return strings.Trim(slug, "-")
}
