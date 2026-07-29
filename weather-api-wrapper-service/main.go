package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"weather-api-service/weather"

	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("missing .env file")
	}
	apiKey, ok := os.LookupEnv("WEATHER_API_KEY")
	if !ok {
		log.Fatal("WEATHER_API_KEY is required")
	}
	redisAddr, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		log.Fatal("REDIS_ADDR is required")
	}
	c := weather.NewClient(apiKey, redisAddr)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather/{city}", func(w http.ResponseWriter, r *http.Request) {
		city := r.PathValue("city")
		if strings.TrimSpace(city) == "" {
			log.Println("client does not fill the city")
			http.Error(w, "city is required", http.StatusBadRequest)
			return
		}
		resp, err := c.GetWeather(r.Context(), city)
		if err != nil {
			log.Printf("GetWeather failed for city=%q: %v", city, err)
			switch {
			case errors.Is(err, weather.ErrInvalidCity):
				http.Error(w, "city not found", http.StatusNotFound)
			case errors.Is(err, weather.ErrUpstreamUnavailable):
				http.Error(w, "weather service temporarily unavailable", http.StatusBadGateway)
			default:
				http.Error(w, "failed to fetch weather data", http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("failed to encode response for city=%q: %v", city, err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("server is running on port: 8080")
	if err := http.ListenAndServe(":8080", rateLimitMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	var mu sync.Mutex
	limiters := make(map[string]*clientLimiter)

	go func() {
		ticker := time.NewTicker(time.Minute * 1)
		for range ticker.C {
			mu.Lock()
			for ip, cl := range limiters {
				if time.Since(cl.lastSeen) > time.Minute*5 {
					delete(limiters, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("error ip: %v", err)
			http.Error(w, "failed to get ip", http.StatusInternalServerError)
			return
		}

		mu.Lock()
		defer mu.Unlock()
		limiter, exists := limiters[ip]
		if !exists {
			limiter = &clientLimiter{
				limiter:  rate.NewLimiter(1, 10),
				lastSeen: time.Now(),
			}
			limiters[ip] = limiter
		} else {
			limiter.lastSeen = time.Now()
		}

		if !limiter.limiter.Allow() {
			log.Printf("client has too many request")
			http.Error(w, "too many request", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
