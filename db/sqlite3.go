package db

import (
	"database/sql"
	"github.com/baraa-almasri/shortsninja/models"
	_ "github.com/mattn/go-sqlite3"
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
		mustInitSQLiteDB(db)
		instance = &SQLite{manager: db}
	}

	return instance
}

// mustInitSQLiteDB creates database's tables if possible, if
func mustInitSQLiteDB(db *sql.DB) {

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS 
    USER (
    	email VARCHAR(255) PRIMARY KEY , 
    	avatar_link VARCHAR(2000),
    	creation_date TIMESTAMP
	);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS 
    URL (
    	short VARCHAR(5) PRIMARY KEY, 
    	full_url VARCHAR(2000),
    	creation_date TIMESTAMP,
    	user_email VARCHAR(255)
	);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS 
	SESSION (
		token VARCHAR(32),
		IP VARCHAR(15)   
	);`)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS 
	URL_DATA (
	    short VARCHAR(5),
	    IP VARCHAR(15),
	    user_agent VARCHAR(255),
	    visit_location VARCHAR(50),
	    visit_time TIMESTAMP
	);`)
}

// AddURL add a new url entry to the database
func (s *SQLite) AddURL(url *models.URL) error {
	_, err := s.manager.Exec(
		`INSERT INTO URL (short, full_url, creation_date, user_email) VALUES (?, ? , CURRENT_TIMESTAMP, ?);`,
		url.Short, url.FullURL, url.UserEmail)
	if err != nil {
		return err
	}

	// the happily ever after
	return nil
}

// RemoveURL sets short URL's row's values to zero, to minimize handlers regeneration :)
func (s *SQLite) RemoveURL(url *models.URL) error {
	_, err := s.manager.Exec(`DELETE FROM URL WHERE short=?;`, url.Short)
	if err != nil {
		return err
	}

	// the happily ever after
	return nil
}

// GetURL returns the full URL of a short URL
func (s *SQLite) GetURL(shortURL string) (string, error) {
	rows, err := s.manager.Query(`SELECT full_url FROM URL WHERE short = ?;`, shortURL)
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

// GetURLs returns a map that has short URLs of the given user
func (s *SQLite) GetURLs(user *models.User) ([]*models.URL, error) {
	rows, err := s.manager.Query(
		`SELECT short, full_url, creation_date FROM URL WHERE user_email=?;`, user.Email)
	if err != nil {
		return nil, err
	}

	var urls []*models.URL
	var url *models.URL

	for rows.Next() {
		url = new(models.URL)

		err = rows.Scan(&url.Short, &url.FullURL, &url.CreationDate)
		if err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}

	return urls, nil
}

// AddUser adds a user to the database
func (s *SQLite) AddUser(user *models.User) error {
	_, err := s.manager.Exec(
		`INSERT INTO USER (email, avatar_link, creation_date) VALUES (?, ?, CURRENT_TIMESTAMP);`,
		user.Email, user.Avatar)
	if err != nil {
		return err
	}

	// the happily ever after
	return nil
}

// RemoveUser removes the given user from the database
func (s *SQLite) RemoveUser(user *models.User) error {
	_, err := s.manager.Exec(`DELETE FROM USER WHERE email=?;`, user.Email)
	if err != nil {
		return err
	}

	// the happily ever after
	return nil
}
func (s *SQLite) GetUser(user *models.User) (*models.User, error) {
	rows, err := s.manager.Query(`SELECT * FROM USER WHERE email=?;`, user.Email)
	if err != nil {
		return nil, err
	}

	u := new(models.User)

	rows.Next()
	err = rows.Scan(&u.Email, &u.Avatar, &u.CreationDate)
	if err != nil {
		return nil, err
	}

	// happily ever after
	return u, nil
}

// AddURLData adds data of a certain url, and returns an occurred error
func (s *SQLite) AddURLData(urlData *models.URLData) error {
	_, err := s.manager.Exec(
		`INSERT INTO URL_DATA (short, IP, user_agent, visit_location, visit_time) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)`,
		urlData.ShortURL, urlData.IP, urlData.UserAgent, urlData.VisitLocation)
	if err != nil {
		return err
	}

	// happily ever after
	return nil
}

// RemoveURLData removes all the data of a given URL, and returns an occurred error
func (s *SQLite) RemoveURLData(url *models.URL) error {
	_, err := s.manager.Exec(
		`DELETE FROM URL_DATA WHERE short=?)`, url.Short)
	if err != nil {
		return err
	}

	// happily ever after
	return nil
}

// GetURLData returns a slice of URLData of the given URL and an occurred error
func (s *SQLite) GetURLData(url *models.URL) ([]*models.URLData, error) {
	rows, err := s.manager.Query(`SELECT * FROM URL_DATA WHERE short=?`, url.Short)
	if err != nil {
		return nil, err
	}

	var urlDatas []*models.URLData
	var urlData *models.URLData

	for rows.Next() {
		urlData = new(models.URLData)

		err = rows.Scan(&urlData.ShortURL, &urlData.IP, &urlData.UserAgent, &urlData.VisitLocation, &urlData.VisitTime)
		if err != nil {
			return nil, err
		}
		urlDatas = append(urlDatas, urlData)
	}

	// happily ever after
	return urlDatas, nil
}
