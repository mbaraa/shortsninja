package models

import (
	"database/sql"
	"github.com/baraa-almasri/shortsninja/utils"
	_ "github.com/mattn/go-sqlite3"
	"strings"
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

// mustInitSQLiteDB creates database's tables if possible, if not the program crashes
func mustInitSQLiteDB(db *sql.DB) {

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS 
    USER (
    	email VARCHAR(255) PRIMARY KEY , 
    	avatar_link VARCHAR(2000),
    	created_at TIMESTAMP
	);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS 
    URL (
    	short VARCHAR(5) PRIMARY KEY, 
    	full_url VARCHAR(2000),
    	created_at TIMESTAMP,
    	user_email VARCHAR(255)
	);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS 
	SESSION (
		user_email VARCHAR(255),
		token VARCHAR(40),
		expires_at TIMESTAMP
	);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS 
	URL_DATA (
	    short VARCHAR(5),
	    IP VARCHAR(45),
	    user_agent VARCHAR(4000),
	    visit_location VARCHAR(50),
	    visited_at TIMESTAMP
	);`)
	if err != nil {
		panic(err)
	}
}

// AddURL add a new url entry to the database
func (s *SQLite) AddURL(url *URL) error {
	_, err := s.manager.Exec(
		`INSERT INTO URL (short, full_url, created_at, user_email) VALUES (?, ? , CURRENT_TIMESTAMP, ?);`,
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
	rows, err := s.manager.Query(`SELECT short, full_url, created_at, user_email FROM URL WHERE short = ?;`, shortURL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	url := new(URL)
	rows.Next()
	var timeStamp time.Time
	err = rows.Scan(&url.Short, &url.FullURL, &timeStamp, &url.UserEmail)
	url.Created = (new(utils.TimeDurationFormatter)).GetDurationSince(timeStamp.Unix())
	if err != nil {
		return nil, err
	}

	// the happily ever after
	return url, nil
}

// GetURLs returns a map that has short URLs of the given user
func (s *SQLite) GetURLs(user *User) ([]*URL, error) {
	rows, err := s.manager.Query(
		`SELECT short, full_url, created_at FROM URL WHERE user_email=?;`, user.Email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	alter := false
	var urls []*URL
	var url *URL
	var timeStamp time.Time

	for rows.Next() {
		url = new(URL)

		url.Alter = alter
		err = rows.Scan(&url.Short, &url.FullURL, &timeStamp)
		if err != nil {
			return nil, err
		}

		url.Created = (new(utils.TimeDurationFormatter)).GetDurationSince(timeStamp.Unix())
		url.Visits = 0
		urls = append(urls, url)
		alter = !alter
	}

	return urls, nil
}

// AddUser adds a user to the database
func (s *SQLite) AddUser(user *User) error {
	_, err := s.manager.Exec(
		`INSERT INTO USER (email, avatar_link, created_at) VALUES (?, ?, CURRENT_TIMESTAMP);`,
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
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u := new(User)
	timeStamp := new(time.Time)

	rows.Next()
	err = rows.Scan(&u.Email, &u.Avatar, &timeStamp)
	if err != nil {
		return nil, err
	}
	u.Created = (new(utils.TimeDurationFormatter)).GetDurationSince(timeStamp.Unix())

	// happily ever after
	return u, nil
}

// AddURLData adds data of a certain url, and returns an occurred error
func (s *SQLite) AddURLData(urlData *URLData) error {
	_, err := s.manager.Exec(
		`INSERT INTO URL_DATA (short, IP, user_agent, visit_location, visited_at) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)`,
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
		`DELETE FROM URL_DATA WHERE short=?`, url.Short)
	if err != nil {
		return err
	}

	// happily ever after
	return nil
}

// GetURLData returns a slice of URLData of the given URL and an occurred error
func (s *SQLite) GetURLData(url *URL) ([]*URLData, error) {
	rows, err := s.manager.Query(`SELECT * FROM URL_DATA WHERE short=?`, url.Short)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	alter := false
	var urlDatas []*URLData
	var urlData *URLData
	var timeStamp *time.Time

	for rows.Next() {
		urlData = new(URLData)
		timeStamp = new(time.Time)

		urlData.Alter = alter
		err = rows.Scan(&urlData.ShortURL, &urlData.IP, &urlData.UserAgent, &urlData.VisitLocation, &timeStamp)
		if err != nil {
			return nil, err
		}
		urlData.UserAgent = urlData.UserAgent[:strings.Index(urlData.UserAgent, ")")+1]
		urlData.VisitTime = timeStamp.Format("Mon Jan/2/2006 15:04 MST")

		urlDatas = append(urlDatas, urlData)
		alter = !alter
	}

	// happily ever after
	return urlDatas, nil
}

// AddSession adds a new session to the database
func (s *SQLite) AddSession(sess *Session) error {
	_, err := s.manager.Exec(`INSERT INTO SESSION (user_email, token, expires_at) VALUES (? ,?, ?)`,
		sess.UserEmail, sess.Token, sess.ExpiresAt)
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
func (s *SQLite) GetSession(sess *Session) (*Session, error) {
	rows, err := s.manager.Query(`SELECT * FROM SESSION WHERE token=?`, sess.Token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	session := new(Session)
	rows.Next()
	err = rows.Scan(&session.UserEmail, &session.Token, &session.ExpiresAt)
	if err != nil || session.ExpiresAt.Unix() < time.Now().Unix() { // expired session :)
		return nil, err
	}

	// happily ever after
	return session, nil
}
