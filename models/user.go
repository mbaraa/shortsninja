package models

import "time"

// User defines the properties of a user
type User struct {
	Email        string
	Avatar       string
	CreationDate *time.Time
}
