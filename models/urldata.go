package models

import "time"

// URLData defines the properties of URL's tracking data
type URLData struct {
	ShortURL      string     `json:"short_url"`
	IP            string     `json:"ip"`
	UserAgent     string     `json:"user_agent"`
	VisitLocation string     `json:"visit_location"`
	VisitTime     *time.Time `json:"visit_time"`
}
