package twitch

import (
	"fmt"
	"regexp"
	"strconv"
)

// ParseDurationToSeconds converts Twitch duration format (e.g., "1h2m3s", "45m12s", "30s") to seconds
func ParseDurationToSeconds(duration string) int {
	if duration == "" {
		return 0
	}

	// Regex to parse hours, minutes, and seconds
	re := regexp.MustCompile(`(?:(\d+)h)?(?:(\d+)m)?(?:(\d+)s)?`)
	matches := re.FindStringSubmatch(duration)

	if len(matches) < 4 {
		return 0
	}

	var totalSeconds int

	// Parse hours
	if matches[1] != "" {
		if hours, err := strconv.Atoi(matches[1]); err == nil {
			totalSeconds += hours * 3600
		}
	}

	// Parse minutes
	if matches[2] != "" {
		if minutes, err := strconv.Atoi(matches[2]); err == nil {
			totalSeconds += minutes * 60
		}
	}

	// Parse seconds
	if matches[3] != "" {
		if seconds, err := strconv.Atoi(matches[3]); err == nil {
			totalSeconds += seconds
		}
	}

	return totalSeconds
}

// FormatSecondsToHMS converts seconds to HH:MM:SS format for display
func FormatSecondsToHMS(seconds int) string {
	if seconds <= 0 {
		return "00:00:00"
	}

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}

// FormatSecondsToCompact converts seconds to compact format (e.g., "1h 23m", "45m 30s", "2m 15s")
func FormatSecondsToCompact(seconds int) string {
	if seconds <= 0 {
		return "0s"
	}

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	if hours > 0 {
		if minutes > 0 {
			return fmt.Sprintf("%dh %dm", hours, minutes)
		}
		return fmt.Sprintf("%dh", hours)
	} else if minutes > 0 {
		if secs > 0 {
			return fmt.Sprintf("%dm %ds", minutes, secs)
		}
		return fmt.Sprintf("%dm", minutes)
	} else {
		return fmt.Sprintf("%ds", secs)
	}
}
