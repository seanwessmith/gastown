// Package weather provides ASCII art weather display functionality.
package weather

import (
	"fmt"
	"strings"
)

// Condition represents a weather condition type.
type Condition string

const (
	Sunny  Condition = "sunny"
	Rainy  Condition = "rainy"
	Cloudy Condition = "cloudy"
	Snowy  Condition = "snowy"
)

// WeatherData holds weather information for rendering.
type WeatherData struct {
	Condition   Condition
	Temperature float64
	Humidity    int
	WindSpeed   float64
	Location    string
}

// ANSI color codes for terminal output.
const (
	colorReset  = "\033[0m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[97m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

// ASCII art for sunny weather.
var sunnyArt = []string{
	"       \\   |   /       ",
	"         .-'-.         ",
	"    --- (     ) ---    ",
	"         `-.-'         ",
	"       /   |   \\       ",
}

// ASCII art for rainy weather.
var rainyArt = []string{
	"       .-~~~-.         ",
	"      (       )        ",
	"    (          )       ",
	"     `-.___.-'         ",
	"      ' ' ' ' '        ",
	"     ' ' ' ' '         ",
}

// ASCII art for cloudy weather.
var cloudyArt = []string{
	"                       ",
	"       .-~~~-.         ",
	"      (       )        ",
	"    (          )       ",
	"     `-.___.-'         ",
	"                       ",
}

// ASCII art for snowy weather.
var snowyArt = []string{
	"       .-~~~-.         ",
	"      (       )        ",
	"    (          )       ",
	"     `-.___.-'         ",
	"      *  *  *  *       ",
	"     *  *  *  *        ",
}

// getArt returns the ASCII art and color for a given condition.
func getArt(c Condition) ([]string, string) {
	switch c {
	case Sunny:
		return sunnyArt, colorYellow
	case Rainy:
		return rainyArt, colorBlue
	case Cloudy:
		return cloudyArt, colorGray
	case Snowy:
		return snowyArt, colorCyan
	default:
		return cloudyArt, colorGray
	}
}

// conditionLabel returns a human-readable label for the condition.
func conditionLabel(c Condition) string {
	switch c {
	case Sunny:
		return "Sunny"
	case Rainy:
		return "Rainy"
	case Cloudy:
		return "Cloudy"
	case Snowy:
		return "Snowy"
	default:
		return "Unknown"
	}
}

// Render prints formatted colorful weather output to the terminal.
func Render(weather WeatherData) {
	art, artColor := getArt(weather.Condition)

	// Header
	fmt.Println()
	fmt.Printf("%s%s", colorBold, colorWhite)
	if weather.Location != "" {
		fmt.Printf("  Weather for %s", weather.Location)
	} else {
		fmt.Print("  Current Weather")
	}
	fmt.Printf("%s\n", colorReset)
	fmt.Println(strings.Repeat("-", 30))

	// ASCII art
	fmt.Print(artColor)
	for _, line := range art {
		fmt.Println(line)
	}
	fmt.Print(colorReset)

	// Weather details
	fmt.Println(strings.Repeat("-", 30))
	fmt.Printf("  %sCondition:%s  %s\n", colorBold, colorReset, conditionLabel(weather.Condition))
	fmt.Printf("  %sTemp:%s       %.1f°C\n", colorBold, colorReset, weather.Temperature)
	fmt.Printf("  %sHumidity:%s   %d%%\n", colorBold, colorReset, weather.Humidity)
	fmt.Printf("  %sWind:%s       %.1f km/h\n", colorBold, colorReset, weather.WindSpeed)
	fmt.Println()
}

// RenderCompact prints a compact single-line weather summary.
func RenderCompact(weather WeatherData) {
	_, artColor := getArt(weather.Condition)
	icon := getIcon(weather.Condition)
	fmt.Printf("%s%s%s %.1f°C | %s\n", artColor, icon, colorReset, weather.Temperature, weather.Location)
}

// getIcon returns a simple icon character for the condition.
func getIcon(c Condition) string {
	switch c {
	case Sunny:
		return "☀"
	case Rainy:
		return "🌧"
	case Cloudy:
		return "☁"
	case Snowy:
		return "❄"
	default:
		return "?"
	}
}
