package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

const cacheDir = ".cache-proxy"

type Entry struct {
	StatusCode int         `json:"status_code"`
	Header     http.Header `json:"header"`
	Body       []byte      `json:"body"`
}

func Key(r *http.Request) string {
	raw := r.Method + ":" + r.URL.Path + ":" + r.URL.RawQuery
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func Get(key string) (*Entry, bool) {
	path := filepath.Join(cacheDir, key+".json")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}

	return &entry, true
}

func Set(key string, entry *Entry) error {
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	path := filepath.Join(cacheDir, key+".json")

	return os.WriteFile(path, data, 0644)
}

func Clear() error {
	if err := os.RemoveAll(cacheDir); err != nil {
		return err
	}
	return os.MkdirAll(cacheDir, 0755)
}
