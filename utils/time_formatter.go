package utils

import (
	"fmt"
	"math"
	"time"
)

// TimeDurationFormatter is a wrapping for getting duration since a time stamp
type TimeDurationFormatter struct{}

// GetDurationSince returns a string with number of days, hours, and minutes since the given timestamp
func (tf *TimeDurationFormatter) GetDurationSince(timestamp int64) string {
	duration := time.Since(time.Unix(timestamp, 0))

	days := tf.getDaysSince(duration)
	hours := tf.getHoursSince(duration)
	minutes := tf.getMinutesSince(duration)
	durationStr := ""

	if days != "" {
		durationStr = days
	} else {
		durationStr += hours + minutes
	}

	return durationStr + " ago"
}

func (tf *TimeDurationFormatter) getDaysSince(duration time.Duration) string {
	days := duration.Hours() / 24
	if days >= 1 {
		return fmt.Sprintf("%d days", int(days))
	}
	return ""
}

func (tf *TimeDurationFormatter) getHoursSince(duration time.Duration) string {
	hours := math.Mod(duration.Hours(), 24)
	if hours >= 1 && hours < 24 {
		return fmt.Sprintf("%d hours | ", int(hours))
	}
	return ""
}

func (tf *TimeDurationFormatter) getMinutesSince(duration time.Duration) string {
	minutes := math.Mod(duration.Hours()*60, 60)
	if minutes >= 0 && minutes < 60 {
		return fmt.Sprintf("%d minutes", int(minutes))
	}
	return ""
}
