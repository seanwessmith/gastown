// Package weather provides functionality to fetch weather data from wttr.in API.
package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Weather contains the current weather conditions for a location.
type Weather struct {
	Location   string  `json:"location"`
	TempC      int     `json:"temp_c"`
	TempF      int     `json:"temp_f"`
	Conditions string  `json:"conditions"`
	Humidity   int     `json:"humidity"`
	WindKmph   int     `json:"wind_kmph"`
	WindMph    int     `json:"wind_mph"`
	WindDir    string  `json:"wind_dir"`
}

// Fetcher fetches weather data from the wttr.in API.
type Fetcher struct {
	client  *http.Client
	baseURL string
}

// NewFetcher creates a new weather fetcher with default settings.
func NewFetcher() *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://wttr.in",
	}
}

// Fetch retrieves weather data for the specified location.
// Location can be a city name, airport code, or coordinates.
func (f *Fetcher) Fetch(ctx context.Context, location string) (*Weather, error) {
	if location == "" {
		return nil, fmt.Errorf("location cannot be empty")
	}

	reqURL := fmt.Sprintf("%s/%s?format=j1", f.baseURL, url.PathEscape(location))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var wttrResp wttrResponse
	if err := json.NewDecoder(resp.Body).Decode(&wttrResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return parseResponse(&wttrResp, location)
}

// wttrResponse represents the JSON response from wttr.in with format=j1.
type wttrResponse struct {
	CurrentCondition []currentCondition `json:"current_condition"`
	NearestArea      []nearestArea      `json:"nearest_area"`
}

type currentCondition struct {
	TempC       string        `json:"temp_C"`
	TempF       string        `json:"temp_F"`
	Humidity    string        `json:"humidity"`
	WindspeedKmph string      `json:"windspeedKmph"`
	WindspeedMiles string     `json:"windspeedMiles"`
	WindDir16Point string     `json:"winddir16Point"`
	WeatherDesc []weatherDesc `json:"weatherDesc"`
}

type weatherDesc struct {
	Value string `json:"value"`
}

type nearestArea struct {
	AreaName []areaName `json:"areaName"`
}

type areaName struct {
	Value string `json:"value"`
}

func parseResponse(resp *wttrResponse, fallbackLocation string) (*Weather, error) {
	if len(resp.CurrentCondition) == 0 {
		return nil, fmt.Errorf("no current conditions in response")
	}

	cc := resp.CurrentCondition[0]

	location := fallbackLocation
	if len(resp.NearestArea) > 0 && len(resp.NearestArea[0].AreaName) > 0 {
		location = resp.NearestArea[0].AreaName[0].Value
	}

	conditions := ""
	if len(cc.WeatherDesc) > 0 {
		conditions = cc.WeatherDesc[0].Value
	}

	return &Weather{
		Location:   location,
		TempC:      parseInt(cc.TempC),
		TempF:      parseInt(cc.TempF),
		Conditions: conditions,
		Humidity:   parseInt(cc.Humidity),
		WindKmph:   parseInt(cc.WindspeedKmph),
		WindMph:    parseInt(cc.WindspeedMiles),
		WindDir:    cc.WindDir16Point,
	}, nil
}

func parseInt(s string) int {
	var v int
	fmt.Sscanf(s, "%d", &v)
	return v
}
