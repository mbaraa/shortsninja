package controllers

import (
	"github.com/baraa-almasri/shortsninja/models"
	"github.com/baraa-almasri/useless"
	"net/http"
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

// CreateAndUpdate creates a new short url that doesn't exist in the db,
// adds the new short URL to the database and returns the assigned short URL
// TODO
// split into createUUID() and addURLToDB()
func (um *URLManager) CreateAndUpdate(url string, user *models.User) string {
	// storing the generated short url so it can be returned :)
	var newURL *models.URL
	// loop until the generated short url doesn't exist in the db
	for {
		newURL = &models.URL{
			Short:     um.randomizer.GetRandomAlphanumString(5),
			FullURL:   url,
			UserEmail: user.Email,
		}
		if um.db.AddURL(newURL) == nil {
			break
		}
	}
	return newURL.Short
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
