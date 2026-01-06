package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Weather holds the weather data for a location
type Weather struct {
	Location    string
	Temperature float64
	Units       string
	Condition   string
	Humidity    int
	WindSpeed   float64
	Forecasts   []DayForecast
}

// DayForecast holds forecast data for a single day
type DayForecast struct {
	Date    string
	TempMax float64
	TempMin float64
	Condition string
}

// geoResponse is the response from Open-Meteo geocoding API
type geoResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

// weatherResponse is the response from Open-Meteo weather API
type weatherResponse struct {
	Current struct {
		Temperature    float64 `json:"temperature_2m"`
		RelativeHumidity int   `json:"relative_humidity_2m"`
		WindSpeed      float64 `json:"wind_speed_10m"`
		WeatherCode    int     `json:"weather_code"`
	} `json:"current"`
	Daily struct {
		Time        []string  `json:"time"`
		TempMax     []float64 `json:"temperature_2m_max"`
		TempMin     []float64 `json:"temperature_2m_min"`
		WeatherCode []int     `json:"weather_code"`
	} `json:"daily"`
}

// Fetch retrieves weather data for the given location
func Fetch(location string, days int, units string) (*Weather, error) {
	// Get coordinates for location
	lat, lon, resolvedName, err := geocode(location)
	if err != nil {
		return nil, fmt.Errorf("failed to find location %q: %w", location, err)
	}

	// Build weather API URL
	tempUnit := "celsius"
	windUnit := "kmh"
	if units == "imperial" {
		tempUnit = "fahrenheit"
		windUnit = "mph"
	}

	apiURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,relative_humidity_2m,wind_speed_10m,weather_code&daily=temperature_2m_max,temperature_2m_min,weather_code&temperature_unit=%s&wind_speed_unit=%s&forecast_days=%d",
		lat, lon, tempUnit, windUnit, days,
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var wr weatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&wr); err != nil {
		return nil, fmt.Errorf("failed to decode weather response: %w", err)
	}

	// Build weather result
	weather := &Weather{
		Location:    resolvedName,
		Temperature: wr.Current.Temperature,
		Units:       units,
		Condition:   weatherCodeToCondition(wr.Current.WeatherCode),
		Humidity:    wr.Current.RelativeHumidity,
		WindSpeed:   wr.Current.WindSpeed,
	}

	// Add forecasts
	for i := 0; i < len(wr.Daily.Time) && i < days; i++ {
		weather.Forecasts = append(weather.Forecasts, DayForecast{
			Date:      wr.Daily.Time[i],
			TempMax:   wr.Daily.TempMax[i],
			TempMin:   wr.Daily.TempMin[i],
			Condition: weatherCodeToCondition(wr.Daily.WeatherCode[i]),
		})
	}

	return weather, nil
}

func geocode(location string) (lat, lon float64, name string, err error) {
	geoURL := fmt.Sprintf(
		"https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1",
		url.QueryEscape(location),
	)

	resp, err := http.Get(geoURL)
	if err != nil {
		return 0, 0, "", err
	}
	defer resp.Body.Close()

	var gr geoResponse
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		return 0, 0, "", err
	}

	if len(gr.Results) == 0 {
		return 0, 0, "", fmt.Errorf("location not found")
	}

	return gr.Results[0].Latitude, gr.Results[0].Longitude, gr.Results[0].Name, nil
}

func weatherCodeToCondition(code int) string {
	switch {
	case code == 0:
		return "Clear sky"
	case code <= 3:
		return "Partly cloudy"
	case code <= 49:
		return "Foggy"
	case code <= 59:
		return "Drizzle"
	case code <= 69:
		return "Rain"
	case code <= 79:
		return "Snow"
	case code <= 84:
		return "Rain showers"
	case code <= 86:
		return "Snow showers"
	case code <= 99:
		return "Thunderstorm"
	default:
		return "Unknown"
	}
}
