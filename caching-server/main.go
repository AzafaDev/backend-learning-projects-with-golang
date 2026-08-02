package main

import (
	"caching-server/internal/cache"
	"caching-server/internal/proxy"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// --- CLI flags -----------------------------------------------------
	port := flag.Int("port", 3000, "port on which the caching proxy server will run")
	origin := flag.String("origin", "", "URL of the server to which requests will be forwarded")
	clearCache := flag.Bool("clear-cache", false, "clear the cache and exit")
	flag.Parse()

	// `caching-proxy --clear-cache` just needs to wipe the cache and quit,
	// it doesn't start a server.
	if *clearCache {
		if err := cache.Clear(); err != nil {
			log.Fatalf("failed to clear cache: %v", err)
		}
		fmt.Println("Cache cleared.")
		os.Exit(0)
	}

	if *origin == "" {
		fmt.Println("Usage:")
		fmt.Println("  caching-proxy --port <number> --origin <url>")
		fmt.Println("  caching-proxy --clear-cache")
		os.Exit(1)
	}

	handler, err := proxy.New(*origin)
	if err != nil {
		log.Fatalf("failed to create proxy: %v", err)
	}

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("caching-proxy listening on %s, forwarding to %s", addr, *origin)
	log.Fatal(http.ListenAndServe(addr, handler))
}
