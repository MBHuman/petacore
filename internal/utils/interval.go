package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParseInterval parses PostgreSQL interval strings like '1d', '1 day', '2 hours', '1 year 2 months', etc.
// Returns duration in microseconds
func ParseInterval(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty interval string")
	}

	// Pattern to match number + unit pairs
	// Supports: years, months, weeks, days, hours, minutes, seconds, milliseconds, microseconds
	// Also supports abbreviated forms: y, mon, w, d, h, m, s, ms, us
	re := regexp.MustCompile(`([+-]?\d+(?:\.\d+)?)\s*([a-zA-Z]+)`)
	matches := re.FindAllStringSubmatch(s, -1)

	if len(matches) == 0 {
		// Try parsing as simple duration without unit (assume seconds)
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid interval format: %s", s)
		}
		return int64(val * float64(time.Second.Microseconds())), nil
	}

	var totalMicros int64

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		numStr := match[1]
		unit := strings.ToLower(match[2])

		val, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number in interval: %s", numStr)
		}

		switch unit {
		case "y", "year", "years":
			// 1 year = 365.25 days (accounting for leap years)
			totalMicros += int64(val * 365.25 * 24 * float64(time.Hour.Microseconds()))
		case "mon", "month", "months":
			// 1 month = 30 days (approximation)
			totalMicros += int64(val * 30 * 24 * float64(time.Hour.Microseconds()))
		case "w", "week", "weeks":
			totalMicros += int64(val * 7 * 24 * float64(time.Hour.Microseconds()))
		case "d", "day", "days":
			totalMicros += int64(val * 24 * float64(time.Hour.Microseconds()))
		case "h", "hour", "hours":
			totalMicros += int64(val * float64(time.Hour.Microseconds()))
		case "m", "min", "minute", "minutes", "mins":
			totalMicros += int64(val * float64(time.Minute.Microseconds()))
		case "s", "sec", "second", "seconds", "secs":
			totalMicros += int64(val * float64(time.Second.Microseconds()))
		case "ms", "millisecond", "milliseconds":
			totalMicros += int64(val * float64(time.Millisecond.Microseconds()))
		case "us", "microsecond", "microseconds":
			totalMicros += int64(val)
		default:
			return 0, fmt.Errorf("unknown interval unit: %s", unit)
		}
	}

	return totalMicros, nil
}

// FormatInterval formats microseconds as a human-readable interval string
func FormatInterval(micros int64) string {
	if micros == 0 {
		return "00:00:00"
	}

	negative := micros < 0
	if negative {
		micros = -micros
	}

	// Convert to days, hours, minutes, seconds, microseconds
	days := micros / (24 * int64(time.Hour.Microseconds()))
	micros %= 24 * int64(time.Hour.Microseconds())

	hours := micros / int64(time.Hour.Microseconds())
	micros %= int64(time.Hour.Microseconds())

	minutes := micros / int64(time.Minute.Microseconds())
	micros %= int64(time.Minute.Microseconds())

	seconds := micros / int64(time.Second.Microseconds())
	micros %= int64(time.Second.Microseconds())

	var result string
	if negative {
		result = "-"
	}

	if days > 0 {
		result += fmt.Sprintf("%d day", days)
		if days != 1 {
			result += "s"
		}
		result += " "
	}

	result += fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

	if micros > 0 {
		result += fmt.Sprintf(".%06d", micros)
	}

	return result
}
