package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
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

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetWeather(city string) (Response, error) {
	cityParam := url.QueryEscape(city)
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
	var vc vcResponse
	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&vc); err != nil {
		return Response{}, err
	}
	response = Response{
		Temperature: vc.CurrentConditions.Temp,
		City:        vc.Address,
		Condition:   vc.CurrentConditions.Conditions,
	}
	return response, nil
}
