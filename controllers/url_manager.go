package controllers

import (
	"github.com/baraa-almasri/shortsninja/models"
	"github.com/baraa-almasri/useless"
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"strings"
)

// URLManager holds URL operations handlers
type URLManager struct {
	urlValidator *URLValidator
	reqData      *RequestDataManager
	userManager  *UserManager
	randomizer   *useless.RandASCII
	db           models.Database
}

// NewURLManager returns a new URLManager instance
func NewURLManager(urlValidator *URLValidator, requestDataManager *RequestDataManager,
	userManager *UserManager, randomStringGenerator *useless.RandASCII, db models.Database) *URLManager {
	return &URLManager{
		urlValidator: urlValidator,
		reqData:      requestDataManager,
		userManager:  userManager,
		randomizer:   randomStringGenerator,
		db:           db,
	}
}

// CreateShortURL adds a new short URL to the database
func (um *URLManager) CreateShortURL(url string, user *models.User) string {
	newURL := &models.URL{
		Short:     um.createUniqueShortURL(5),
		FullURL:   url,
		UserEmail: user.Email,
	}
	err := um.db.AddURL(newURL)
	if err != nil {
		um.CreateShortURL(url, user)
	}
	return newURL.Short
}

// createUniqueShortURL
func (um *URLManager) createUniqueShortURL(length int) string {
	uuidGen := uuid.New()
	short := strings.ReplaceAll(uuidGen.String(), "-", "")
	rand.Seed(int64(uuidGen.ClockSequence()))

	lastIndex := rand.Intn(len(short)-length+1) + length

	return short[lastIndex-length : lastIndex]
}

// GetFullURL returns the full URL for the given short URL
func (um *URLManager) GetFullURL(shortURL string) string {
	url := um.GetURL(shortURL)

	// happily ever after
	return url.FullURL
}

// GetURL returns a url as is from the database if no url exists, well the caller gets rick rolled :)
func (um *URLManager) GetURL(shortURL string) *models.URL {
	url, err := um.db.GetURL(shortURL)
	if err != nil {
		return &models.URL{FullURL: "/no_url/"} // get rick rolled :)
	}
	return url
}

// GetURLData returns a slice of URLData of the given URL
func (um *URLManager) GetURLData(url *models.URL) []*models.URLData {
	urlData, _ := um.db.GetURLData(url)
	return urlData
}

// TrackURLData stores short URL data if the short URL is assigned to a specific user
func (um *URLManager) TrackURLData(req *http.Request) {
	data := um.reqData.GetURLDataFromRequestData(req)
	data.ShortURL = req.URL.Path[1:]

	url, err := um.db.GetURL(data.ShortURL)
	if err != nil {
		return
	}

	user, err := um.db.GetUser(&models.User{Email: url.UserEmail})
	if err == nil && user.Email != "" {
		_ = um.db.AddURLData(data)
	}
}

// RemoveURL removes a corresponding short URL from the database
func (um *URLManager) RemoveURL(shortURL string, request *http.Request) {
	user := um.userManager.GetUserFromToken(request)
	url, err := um.db.GetURL(shortURL)

	if err == nil && user.Email == url.UserEmail {
		_ = um.db.RemoveURLData(url)
		_ = um.db.RemoveURL(url)
	}
}
