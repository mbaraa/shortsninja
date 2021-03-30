package models

// Database represents a database that can be used for the program
type Database interface {
	AddURL(url *URL) error
	RemoveURL(url *URL) error
	GetURL(shortURL string) (string, error)
	GetURLs(user *User) ([]*URL, error)

	AddUser(user *User) error
	RemoveUser(user *User) error
	GetUser(user *User) (*User, error)

	AddURLData(urlData *URLData) error
	RemoveURLData(url *URL) error
	GetURLData(url *URL) ([]*URLData, error)

	AddSession(sess *Session) error
	RemoveSession(sess *Session) error
	GetSession(token string) (*Session, error)
}

// TODO
// add admin user functions :)
