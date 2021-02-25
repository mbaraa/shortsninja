package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"time"
)

// SQLite represents a sqlite database for the program
type SQLite struct {
	manager *sql.DB
}

// hmm...
var instance *SQLite = nil

// NewSQLiteDB returns a singleton SQLite instance
func NewSQLiteDB() *SQLite {
	if instance == nil {
		db, err := sql.Open("sqlite3", "./ninja.db")
		if err != nil {
			panic(err)
		}
		initDB(db)
		instance = &SQLite{manager: db}
	}

	return instance
}

// initDB creates database's tables
func initDB(db *sql.DB) {
	createURLsTable, _ := db.Prepare(
		`CREATE TABLE IF NOT EXISTS urls (
    		user_id INT, 
    		short_url VARCHAR(4) UNIQUE, 
    		url VARCHAR(2000), 
    		creation_date TIMESTAMP NOT NULL
        );`)
	_, _ = createURLsTable.Exec()

	createUsersTable, _ := db.Prepare(
		`CREATE TABLE IF NOT EXISTS users (
    		id INT PRIMARY KEY,
    		name VARCHAR(200),
    		email VARCHAR(255) NOT NULL,
    		password VARCHAR(255) NOT NULL, 
    		creation_date TIMESTAMP NOT NULL 
    	);`)
	_, _ = createUsersTable.Exec()
}

// AddURL add a new url entry to the database
func (s *SQLite) AddURL(userID int, shortURL, url string, creationDate time.Time) error {
	stmt, err := s.manager.Prepare(`INSERT INTO urls (user_id, short_url, url, creation_date) VALUES (?, ?, ?, ?);`)
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(userID, shortURL, url, creationDate)
	if err != nil {
		return err
	}

	// the happily ever after
	return nil
}

// RemoveURL sets short URL's row's values to zero, to minimize handlers regeneration :)
func (s *SQLite) RemoveURL(shortURL string) error {
	stmt, err := s.manager.Prepare(`UPDATE urls SET user_id=?, url=?, creation_date=? WHERE short_url = ?;`)
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(0, "", time.Time{}, shortURL)
	if err != nil {
		return err
	}

	// the happily ever after
	return nil
}

// UpdateURL updates a short URL's row's values
func (s *SQLite) UpdateURL(userID int, shortURL, url string, creationDate time.Time) error {
	stmt, err := s.manager.Prepare(`UPDATE urls SET user_id = ?, url = ?, creation_date = ? WHERE short_url = ?;`)
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(userID, url, creationDate, shortURL)
	if err != nil {
		return err
	}

	// the happily ever after
	return nil
}

// GetURLs returns a map that has short URLs with the corresponding shortened URL
func (s *SQLite) GetURLs() (map[string]string, error) {
	rows, err := s.manager.Query(`SELECT url, short_url FROM urls`)
	if err != nil {
		return nil, err
	}

	urls := make(map[string]string)
	var url, shortURL string
	for rows.Next() {
		err = rows.Scan(&url, &shortURL)
		if err != nil {
			return nil, err
		}
		urls[shortURL] = url
	}

	return urls, nil
}

// GetURLs returns a map that has all the data of a specific user
func (s *SQLite) GetURLsByUserID(userID int) (map[string][]string, error) {
	rows, err := s.manager.Query(`SELECT short_url, url, creation_date FROM urls WHERE user_id = ?;`, userID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var creationDate time.Time
	urlRow := make([]string, 3)
	urls := make(map[string][]string)

	for rows.Next() {
		rows.Scan(&urlRow[0], &urlRow[1], &creationDate)
		urlRow[2] = creationDate.Format("Mon Jan/2/2006 15:04 MST")
		urls[urlRow[0]] = urlRow
	}

	// the happily ever after
	return urls, nil
}

// GetURL returns the full URL of a short URL
func (s *SQLite) GetURL(shortURL string) (string, error) {
	rows, err := s.manager.Query(`SELECT url FROM urls WHERE short_url = ?;`, shortURL)
	defer rows.Close()
	if err != nil {
		return "", err
	}

	var url string
	rows.Next()
	err = rows.Scan(&url)
	if err != nil {
		return "", err
	}

	// the happily ever after
	return url, nil
}

// GetSURLs returns all the short URLs that matches the given one with their given data
func (s *SQLite) GetSURLs(url string) (map[int][]string, error) {
	rows, err := s.manager.Query(`SELECT user_id, short_url, creation_date FROM urls WHERE url = ?;`, url)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var userID int
	var creationDate time.Time
	urlRow := make([]string, 2)
	urls := make(map[int][]string)

	for rows.Next() {
		err = rows.Scan(&userID, &urlRow[1], &creationDate)
		if err != nil {
			return nil, err
		}

		urlRow[0] = strconv.Itoa(userID)
		urlRow[1] = creationDate.Format("Mon Jan/2/2006 15:04 MST")
		urls[userID] = urlRow
	}

	// the happily ever after
	return urls, nil
}

// TODO
// make these fuckers :)
func (s *SQLite) AddUser(id int, name, password, email string) error {
	return nil
}
func (s *SQLite) RemoveUser(id int) error {
	return nil
}
func (s *SQLite) GetUser(id int) error {
	return nil
}
