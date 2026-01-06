package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/gastown/weather-cli/pkg/display"
	"github.com/gastown/weather-cli/pkg/fetcher"
)

var (
	location string
	days     int
	units    string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "weather",
		Short: "Get weather information for a location",
		Long:  `A CLI tool to fetch and display weather information for any city.`,
		RunE:  run,
	}

	rootCmd.Flags().StringVarP(&location, "location", "l", "", "City name (required)")
	rootCmd.Flags().IntVarP(&days, "days", "d", 1, "Forecast days (1-3)")
	rootCmd.Flags().StringVarP(&units, "units", "u", "metric", "Units: metric or imperial")

	rootCmd.MarkFlagRequired("location")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Validate days
	if days < 1 || days > 3 {
		return fmt.Errorf("days must be between 1 and 3")
	}

	// Validate units
	if units != "metric" && units != "imperial" {
		return fmt.Errorf("units must be 'metric' or 'imperial'")
	}

	// Fetch weather
	weather, err := fetcher.Fetch(location, days, units)
	if err != nil {
		return fmt.Errorf("failed to get weather: %w", err)
	}

	// Display results
	display.Print(os.Stdout, weather)

	return nil
}
