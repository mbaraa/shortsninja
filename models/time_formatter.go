package models

import (
	"fmt"
	"math"
	"time"
)

// TimeDurationFormatter is a wrapping for getting duration since a time stamp
type TimeDurationFormatter struct{}

// NewTimeDurationFormatter returns a new TimeDurationFormatter instance
func NewTimeDurationFormatter() *TimeDurationFormatter {
	return new(TimeDurationFormatter)
}

// GetDurationSince returns a string with number of days, hours, and minutes since the given timestamp
func (tf *TimeDurationFormatter) GetDurationSince(timestamp int64) string {
	duration := time.Since(time.Unix(timestamp, 0))
	return tf.getDaysSince(duration) + tf.getHoursSince(duration) +
		tf.getMinutesSince(duration) + " ago"
}

func (tf *TimeDurationFormatter) getDaysSince(duration time.Duration) string {
	days := duration.Hours() / 24
	if days >= 1 {
		return fmt.Sprintf("%d days | ", int(days))
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
