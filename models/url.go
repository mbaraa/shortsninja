package models

import "time"

// URL defines the properties of a certain URL
type URL struct {
	Short        string
	FullURL      string
	CreationDate *time.Time
	UserEmail    string
}
