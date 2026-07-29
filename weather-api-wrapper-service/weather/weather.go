package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	APIKey      string
	BaseURL     string
	HTTPClient  *http.Client
	RedisClient *redis.Client
}

type vcResponse struct {
	ResolvedAddress string `json:"resolvedAddress"`
	Address         string `json:"address"`
	TimeZone        string `json:"timezone"`
	Days            []struct {
		DateTime   string  `json:"datetime"`
		Temp       float64 `json:"temp"`
		TempMax    float64 `json:"tempmax"`
		TempMin    float64 `json:"tempmin"`
		Humadity   float64 `json:"humidity"`
		Conditions string  `json:"conditions"`
		Icon       string  `json:"icon"`
	}
	CurrentConditions struct {
		Temp       float64 `json:"temp"`
		Conditions string  `json:"conditions"`
		Humidity   float64 `json:"humidity"`
		WindSpeed  float64 `json:"windspeed"`
	}
}

type Response struct {
	Temperature float64 `json:"temperature"`
	City        string  `json:"city"`
	Condition   string  `json:"condition"`
}

func NewClient(apiKey, addr string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		RedisClient: redis.NewClient(
			&redis.Options{
				Addr: addr,
			},
		),
	}
}

func (c *Client) GetWeather(ctx context.Context, city string) (Response, error) {
	var vc vcResponse
	var response Response
	cacheKey := strings.ToLower(strings.TrimSpace(city))
	cityParam := url.QueryEscape(cacheKey)
	cached, err := c.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &response); err != nil {
			return Response{}, err
		}
	} else {
		reqUrl := fmt.Sprintf("%s/%s?key=%s&unitGroup=metric", c.BaseURL, cityParam, c.APIKey)
		resp, err := c.HTTPClient.Get(reqUrl)
		if err != nil {
			return Response{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			data, _ := io.ReadAll(resp.Body)
			return Response{}, fmt.Errorf("weather api returned status %d, err: %s", resp.StatusCode, data)
		}
		if err := json.NewDecoder(resp.Body).Decode(&vc); err != nil {
			return Response{}, err
		}
		response = Response{
			Temperature: vc.CurrentConditions.Temp,
			City:        vc.Address,
			Condition:   vc.CurrentConditions.Conditions,
		}
		dataResp, err := json.Marshal(response)
		if err != nil {
			return Response{}, err
		}
		if err := c.RedisClient.Set(ctx, cacheKey, dataResp, 12*time.Hour).Err(); err != nil {
			log.Printf("failed to cache weather for %s: %v\n", cacheKey, err)
		}
	}
	return response, nil
}
