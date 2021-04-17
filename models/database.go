package models

// Database represents a database that can be used for the program
type Database interface {
	AddURL(*URL) error
	RemoveURL(*URL) error
	GetURL(shortURL string) (*URL, error)
	GetURLs(*User) ([]*URL, error)

	AddUser(*User) error
	RemoveUser(*User) error
	GetUser(*User) (*User, error)

	AddURLData(*URLData) error
	RemoveURLData(*URL) error
	GetURLData(*URL) ([]*URLData, error)

	AddSession(*Session) error
	RemoveSession(*Session) error
	GetSession(*Session) (*Session, error)

	GetUsers() ([]*User, error)
	GetAllURLs() ([]*URL, error)
	GetSessions() ([]*Session, error)
}
