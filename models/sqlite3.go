package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
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
		IP VARCHAR(45),
		user_email VARCHAR(255)
	);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS 
	URL_DATA (
	    short VARCHAR(5),
	    IP VARCHAR(45),
	    user_agent VARCHAR(255),
	    visit_location VARCHAR(50),
	    visit_time TIMESTAMP
	);`)
	if err != nil {
		panic(err)
	}
}

// AddURL add a new url entry to the database
func (s *SQLite) AddURL(url *URL) error {
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
func (s *SQLite) RemoveURL(url *URL) error {
	_, err := s.manager.Exec(`DELETE FROM URL WHERE short=?;`, url.Short)
	if err != nil {
		return err
	}

	// the happily ever after
	return nil
}

// GetURL returns the full URL of a short URL
func (s *SQLite) GetURL(shortURL string) (*URL, error) {
	rows, err := s.manager.Query(`SELECT short, full_url, creation_date, user_email FROM URL WHERE short = ?;`, shortURL)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	url := new(URL)
	rows.Next()
	var timeStamp time.Time
	err = rows.Scan(&url.Short, &url.FullURL, &timeStamp, &url.UserEmail)
	url.Created = (&TimeDurationFormatter{}).GetDurationSince(timeStamp.Unix())
	if err != nil {
		return nil, err
	}

	// the happily ever after
	return url, nil
}

// GetURLs returns a map that has short URLs of the given user
func (s *SQLite) GetURLs(user *User) ([]*URL, error) {
	rows, err := s.manager.Query(
		`SELECT short, full_url, creation_date FROM URL WHERE user_email=?;`, user.Email)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var urls []*URL
	var url *URL
	var timeStamp time.Time

	for rows.Next() {
		url = new(URL)

		err = rows.Scan(&url.Short, &url.FullURL, &timeStamp)
		if err != nil {
			return nil, err
		}

		url.Created = (&TimeDurationFormatter{}).GetDurationSince(timeStamp.Unix())
		url.Visits = 0
		urls = append(urls, url)
	}

	return urls, nil
}

// AddUser adds a user to the database
func (s *SQLite) AddUser(user *User) error {
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
func (s *SQLite) RemoveUser(user *User) error {
	_, err := s.manager.Exec(`DELETE FROM USER WHERE email=?;`, user.Email)
	if err != nil {
		return err
	}

	// the happily ever after
	return nil
}

// GetUser returns an existing user from the database
func (s *SQLite) GetUser(user *User) (*User, error) {
	rows, err := s.manager.Query(`SELECT * FROM USER WHERE email=?;`, user.Email)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	u := new(User)

	rows.Next()
	err = rows.Scan(&u.Email, &u.Avatar, &u.CreationDate)
	if err != nil {
		return nil, err
	}

	// happily ever after
	return u, nil
}

// AddURLData adds data of a certain url, and returns an occurred error
func (s *SQLite) AddURLData(urlData *URLData) error {
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
func (s *SQLite) RemoveURLData(url *URL) error {
	_, err := s.manager.Exec(
		`DELETE FROM URL_DATA WHERE short=?)`, url.Short)
	if err != nil {
		return err
	}

	// happily ever after
	return nil
}

// GetURLData returns a slice of URLData of the given URL and an occurred error
func (s *SQLite) GetURLData(url *URL) ([]*URLData, error) {
	rows, err := s.manager.Query(`SELECT * FROM URL_DATA WHERE short=?`, url.Short)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var urlDatas []*URLData
	var urlData *URLData

	for rows.Next() {
		urlData = new(URLData)

		err = rows.Scan(&urlData.ShortURL, &urlData.IP, &urlData.UserAgent, &urlData.VisitLocation, &urlData.VisitTime)
		if err != nil {
			return nil, err
		}
		urlDatas = append(urlDatas, urlData)
	}

	// happily ever after
	return urlDatas, nil
}

// AddSession adds a new session to the database
func (s *SQLite) AddSession(sess *Session) error {
	_, err := s.manager.Exec(`INSERT INTO SESSION (token, IP, user_email) VALUES (?, ? ,?)`,
		sess.Token, sess.IP, sess.UserEmail)
	if err != nil {
		return err
	}

	// happily ever after
	return nil
}

// RemoveSession a specific session from the database
func (s *SQLite) RemoveSession(sess *Session) error {
	_, err := s.manager.Exec(`DELETE FROM SESSION WHERE token=?`, sess.Token)
	if err != nil {
		return err
	}

	// happily ever after
	return nil
}

// GetSession returns a specific session from the database
func (s *SQLite) GetSession(token string) (*Session, error) {
	rows, err := s.manager.Query(`SELECT * FROM SESSION WHERE token=?`, token)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	session := new(Session)
	rows.Next()
	err = rows.Scan(&session.Token, &session.IP, &session.UserEmail)
	if err != nil {
		return nil, err
	}

	// happily ever after
	return session, nil
}
