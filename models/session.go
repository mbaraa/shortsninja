package models

// Session defines session's properties(of a registered user)
type Session struct {
	IP        string
	UserEmail string
	UserAgent string
}
