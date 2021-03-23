package models

import "time"

// User defines the properties of a user
type User struct {
	Email        string     `json:"email"`
	Avatar       string     `json:"avatar"`
	CreationDate *time.Time `json:"creation_date"`
}
