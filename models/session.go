package models

// Session defines session's properties(of a registered user)
type Session struct {
	Token string `json:"token"`
	IP    string `json:"ip"`
}
