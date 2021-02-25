package db

import "time"

// Database represents a database that can be used for the program
type Database interface {
	// AddURL adds a new url entry to the database
	AddURL(userID int, shortURL, url string, creationDate time.Time) error

	// RemoveURL sets short URL's row's values to zero, to minimize handlers regeneration :)
	RemoveURL(shortURL string) error

	// UpdateURL updates a short URL's row's values
	UpdateURL(userID int, shortURL, url string, creationDate time.Time) error

	// GetURLs returns a map that has short URLs with the corresponding shortened URL
	GetURLs() (map[string]string, error)

	// GetURLsByUserID returns a map that has all the data of a specific user
	GetURLsByUserID(userID int) (map[string][]string, error)

	// GetURL returns the full URL of a short URL
	GetURL(shortURL string) (string, error)

	// GetSURLs returns all the short URLs that matches the given one with their given data
	GetSURLs(url string) (map[int][]string, error)

	AddUser(id int, name, password, email string) error
	RemoveUser(id int) error
	GetUser(id int) error
}
