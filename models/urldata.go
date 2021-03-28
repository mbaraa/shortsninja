package models

import "time"

// URLData defines the properties of URL's tracking data
type URLData struct {
	ShortURL      string
	IP            string
	UserAgent     string
	VisitLocation string
	VisitTime     *time.Time
}
