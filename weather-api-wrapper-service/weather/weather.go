package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrInvalidCity         = errors.New("invalid city or location not found")
	ErrUpstreamUnavailable = errors.New("weather provider unavailable")
)

type Client struct {
	apiKey      string
	baseURL     string
	httpClient  *http.Client
	redisClient *redis.Client
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
		Humidity   float64 `json:"humidity"`
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
		apiKey:  apiKey,
		baseURL: "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		redisClient: redis.NewClient(
			&redis.Options{
				Addr: addr,
			},
		),
	}
}

// Close releases the underlying Redis connection. Call it once on shutdown.
func (c *Client) Close() error {
	return c.redisClient.Close()
}

func (c *Client) GetWeather(ctx context.Context, city string) (Response, error) {
	var vc vcResponse
	var response Response
	cacheKey := strings.ToLower(strings.TrimSpace(city))
	cityParam := url.QueryEscape(cacheKey)
	cached, err := c.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &response); err != nil {
			return Response{}, err
		}
	} else {
		if !errors.Is(err, redis.Nil) {
			log.Printf("redis lookup failed for %s, falling back to upstream API: %v", cacheKey, err)
		}
		reqUrl := fmt.Sprintf("%s/%s?key=%s&unitGroup=metric", c.baseURL, cityParam, c.apiKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
		if err != nil {
			return Response{}, err
		}
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return Response{}, fmt.Errorf("%w: %v", ErrUpstreamUnavailable, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			data, _ := io.ReadAll(resp.Body)
			switch {
			case resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound:
				return Response{}, fmt.Errorf("%w: %s", ErrInvalidCity, data)
			case resp.StatusCode >= http.StatusInternalServerError:
				return Response{}, fmt.Errorf("%w: status %d: %s", ErrUpstreamUnavailable, resp.StatusCode, data)
			default:
				return Response{}, fmt.Errorf("weather api returned status %d, err: %s", resp.StatusCode, data)
			}
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
		if err := c.redisClient.Set(ctx, cacheKey, dataResp, 12*time.Hour).Err(); err != nil {
			log.Printf("failed to cache weather for %s: %v\n", cacheKey, err)
		}
	}
	return response, nil
}
