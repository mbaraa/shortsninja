package models

import "time"

// URL defines the properties of a certain URL
type URL struct {
	Short        string     `json:"short"`
	FullURL      string     `json:"full_url"`
	CreationDate *time.Time `json:"creation_date"`
	UserEmail    string     `json:"user_email"`
}
