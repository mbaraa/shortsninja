package models

import "time"

// Session defines session's properties(of a registered user)
type Session struct {
	UserEmail string
	Token     string
	ExpiresAt time.Time
	Alter     bool
}
