// Package cache implements a simple disk-backed cache for HTTP responses.
//
// Design: each cached response is stored as one JSON file inside cacheDir.
// The filename is derived from a hash of the request (method + URL), so
// looking up a cache entry means: hash the incoming request -> check if
// that file exists -> read + decode it if so.
//
// Storing on disk (instead of just an in-memory map) means `--clear-cache`
// can work as its own standalone CLI invocation, without needing to talk
// to a running server process.
package cache

import (
	"net/http"
)

// cacheDir is where cached entries live on disk.
const cacheDir = ".cache-proxy"

// Entry is what gets stored for a single cached response.
// You need to persist enough here to fully reconstruct the response later:
// status code, headers, and body.
type Entry struct {
	StatusCode int         `json:"status_code"`
	Header     http.Header `json:"header"`
	Body       []byte      `json:"body"`
}

// Key derives a cache key for an incoming request.
//
// TODO: Build a deterministic string/hash from the parts of the request
// that make two requests "the same" for caching purposes — at minimum
// method + path + query string. (crypto/sha256 + hex is a common choice.)
//
// Careful: GET /products and POST /products should NOT hit the same
// cache entry.
func Key(r *http.Request) string {
	panic("TODO: implement Key")
}

// Get looks up a cache entry by key.
//
// TODO:
//  1. Build the file path for this key inside cacheDir.
//  2. If the file doesn't exist, return (nil, false) — not an error.
//  3. If it exists, read it and json.Unmarshal into an Entry.
func Get(key string) (*Entry, bool) {
	panic("TODO: implement Get")
}

// Set writes a cache entry to disk.
//
// TODO:
//  1. Make sure cacheDir exists (os.MkdirAll).
//  2. json.Marshal the Entry.
//  3. Write it to the file for this key (os.WriteFile).
func Set(key string, entry *Entry) error {
	panic("TODO: implement Set")
}

// Clear removes all cached entries.
//
// TODO: os.RemoveAll(cacheDir) is enough — the directory gets recreated
// lazily the next time Set is called.
func Clear() error {
	panic("TODO: implement Clear")
}
