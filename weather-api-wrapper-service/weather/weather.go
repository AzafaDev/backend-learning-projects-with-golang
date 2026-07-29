package weather

import (
	"net/http"
	"time"
)

type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

type Response struct {
	Temperature float64 `json:"temperature"`
	City        string  `json:"city"`
	Condition   string  `json:"condition"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "https://weather.visualcrossing.com/VisualCrossing/rest/services/timeline",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetWeather(city string) (Response, error) {
	return Response{
		Temperature: 10,
		City:        city,
		Condition:   "Good",
	}, nil
}
