package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"weather-api-service/weather"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("missing .env file")
	}
	apiKey, ok := os.LookupEnv("WEATHER_API_KEY")
	if !ok {
		log.Fatal("WEATHER_API_KEY is requried")
	}
	c := weather.NewClient(apiKey, "localhost:6379")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather/{city}", func(w http.ResponseWriter, r *http.Request) {
		city := r.PathValue("city")
		resp, err := c.GetWeather(r.Context(), city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("server is running on port: 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
