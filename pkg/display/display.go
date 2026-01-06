package display

import (
	"fmt"
	"io"
	"strings"

	"github.com/gastown/weather-cli/pkg/fetcher"
)

// Print outputs weather data to the given writer
func Print(w io.Writer, weather *fetcher.Weather) {
	tempSymbol := "C"
	windUnit := "km/h"
	if weather.Units == "imperial" {
		tempSymbol = "F"
		windUnit = "mph"
	}

	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Weather for %s\n", weather.Location)
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", len(weather.Location)+12))
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Current Conditions:\n")
	fmt.Fprintf(w, "  Temperature: %.1f°%s\n", weather.Temperature, tempSymbol)
	fmt.Fprintf(w, "  Condition:   %s\n", weather.Condition)
	fmt.Fprintf(w, "  Humidity:    %d%%\n", weather.Humidity)
	fmt.Fprintf(w, "  Wind:        %.1f %s\n", weather.WindSpeed, windUnit)

	if len(weather.Forecasts) > 0 {
		fmt.Fprintf(w, "\nForecast:\n")
		for _, f := range weather.Forecasts {
			fmt.Fprintf(w, "  %s: %.1f°%s / %.1f°%s - %s\n",
				f.Date, f.TempMax, tempSymbol, f.TempMin, tempSymbol, f.Condition)
		}
	}
	fmt.Fprintf(w, "\n")
}
