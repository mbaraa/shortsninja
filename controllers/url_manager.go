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
	userMan      *UserManager
	randomizer   *useless.RandASCII
	db           models.Database
}

// NewURLManager returns a new URLManager instance
func NewURLManager(urlValidator *URLValidator, requestDataManager *RequestDataManager,
	userManager *UserManager, randomStringGenerator *useless.RandASCII, db models.Database) *URLManager {
	return &URLManager{
		urlValidator: urlValidator,
		reqData:      requestDataManager,
		userMan:      userManager,
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
	url, err := um.db.GetURL(shortURL)
	if err != nil {
		return "/no_url/" // get rick rolled :)
	}

	// happily ever after
	return url
}

// TrackURLData stores short URL data if the short URL is assigned to a specific user
func (um *URLManager) TrackURLData(req *http.Request) {
	data := um.reqData.GetURLDataFromRequestData(req)
	data.ShortURL = req.URL.Path[1:]

	user := um.userMan.GetUserFromToken(req)
	if user.Email != "" {
		_ = um.db.AddURLData(data)
	}
}
