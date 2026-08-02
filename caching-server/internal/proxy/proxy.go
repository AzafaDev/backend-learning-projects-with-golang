// Package proxy implements the HTTP handler that either serves a cached
// response or forwards the request to the origin server and caches the
// result.
package proxy

import (
	"caching-server/internal/cache"
	"net/http"
	"net/url"
)

// Handler is the caching reverse proxy.
type Handler struct {
	origin *url.URL
	client *http.Client
}

// New builds a Handler that forwards requests to originURL.
func New(originURL string) (*Handler, error) {
	origin, err := url.Parse(originURL)
	if err != nil {
		return nil, err
	}
	return &Handler{
		origin: origin,
		client: &http.Client{},
	}, nil
}

// ServeHTTP is where the actual caching-proxy behavior lives.
//
// TODO, roughly in this order:
//
//  1. Compute a cache key for the incoming request:
//     key := cache.Key(r)
//
//  2. Check the cache:
//     if entry, ok := cache.Get(key); ok { ... }
//     If found: copy entry.Header onto w.Header(), set w.Header().Set("X-Cache", "HIT"),
//     w.WriteHeader(entry.StatusCode), then w.Write(entry.Body). Return.
//
//  3. Cache miss: build a new request to the origin server.
//     - Same method, path, and query as the incoming request.
//     - Target host/scheme should be h.origin instead of the proxy's own address.
//     - Don't forget to copy the body for non-GET requests (r.Body).
//
//  4. Send it with h.client.Do(req), read the response body
//     (io.ReadAll(resp.Body), then resp.Body.Close()).
//
//  5. Store it in the cache:
//     cache.Set(key, &cache.Entry{StatusCode: resp.StatusCode, Header: resp.Header, Body: body})
//
//  6. Write the response back to the client: copy headers, set
//     w.Header().Set("X-Cache", "MISS"), w.WriteHeader(resp.StatusCode), w.Write(body).
//
// Things worth thinking about as you implement this:
//   - What should happen on a non-GET request — should it be cached at all?
//   - What happens if the origin server is unreachable?
//   - h.origin gives you scheme + host; r.URL gives you path + query of the
//     incoming request. You'll need to combine them for the outgoing request.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = cache.Entry{} // remove once you're using the cache package for real
	panic("TODO: implement ServeHTTP")
}
