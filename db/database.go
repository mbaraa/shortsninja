package db

import "github.com/baraa-almasri/shortsninja/models"

// Database represents a database that can be used for the program
type Database interface {
	AddURL(url *models.URL) error
	RemoveURL(url *models.URL) error
	GetURL(shortURL string) (string, error)
	GetURLs(user *models.User) ([]*models.URL, error)

	AddUser(user *models.User) error
	RemoveUser(user *models.User) error
	GetUser(user *models.User) (*models.User, error)
}

// TODO
// add admin user functions :)
