package controllers

import (
	"encoding/json"
	"github.com/baraa-almasri/shortsninja/models"
	"github.com/baraa-almasri/shortsninja/utils"
	"github.com/baraa-almasri/useless"
	"net/http"
)

// URLManager holds URL operations handlers
type URLManager struct {
	urlValidator *URLValidator
	reqData      *RequestDataManager
	userManager  *UserManager
	randomizer   *useless.RandASCII
	uniqueID     *utils.UniqueID
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
		uniqueID:     utils.NewUniqueID(randomStringGenerator),
		db:           db,
	}
}

// HandleCreateShortURL creates a short url for the given url
// GET /shorten/?url=http://someurl.com
//
func (um *URLManager) HandleCreateShortURL(res http.ResponseWriter, req *http.Request) {
	user := um.userManager.GetUserFromRequest(req)
	url := req.URL.Query()["url"][0]

	resp := make(map[string]interface{})
	if resp["valid_url"] = um.urlValidator.IsURLValid(url); resp["valid_url"].(bool) {
		shortURL := um.createShortURL(url, user)

		resp["url"] = url
		resp["short"] = "shorts.ninja/" + shortURL
	}
	_ = json.NewEncoder(res).Encode(resp)
}

// createShortURL creates and adds a new short URL handler to the database
// if the url exists for the given user it returns its short handler
// if the generated short url exists it regenerates a new short handler
//
func (um *URLManager) createShortURL(url string, user *models.User) string {
	if existingURL := um.getShortURL(url, user); existingURL != nil {
		return existingURL.Short
	}

	newURL := &models.URL{
		Short:     um.uniqueID.GetUniqueString(),
		FullURL:   url,
		UserEmail: user.Email,
	}
	err := um.db.IncrementUserURLsCount(user)
	err = um.db.AddURL(newURL)
	if err != nil {
		um.createShortURL(url, user)
	}
	return newURL.Short
}

// checkURLExistence returns true if the URL exists for a certain user, and false otherwise
func (um *URLManager) checkURLExistence(fullURL string, user *models.User) bool {
	return um.getShortURL(fullURL, user) != nil
}

// getShortURL returns a pointer to an existing URL struct, and nil otherwise
func (um *URLManager) getShortURL(fullURL string, user *models.User) *models.URL {
	urls, err := um.db.GetURLs(user)
	if err != nil {
		return nil
	}

	for _, url := range urls {
		if url.FullURL == fullURL {
			return url
		}
	}
	return nil
}

// HandleGetFullURL retrieves the original url of the given short url
// GET /{[A-Z;0-9;a-z]{4,5}}
func (um *URLManager) HandleGetFullURL(res http.ResponseWriter, req *http.Request) {
	if !um.urlValidator.IsShortURLValid(req.URL.Path[1:]) {
		http.Redirect(res, req, "/no_url/", 302)
		return
	}

	url := um.getFullURL(req.URL.Path[1:])
	um.TrackURLData(req)

	http.Redirect(res, req, url.FullURL, 302)
}

// getFullURL returns a url as is from the database if no url exists, well the caller gets rick rolled :)
//
func (um *URLManager) getFullURL(shortURL string) *models.URL {
	url, err := um.db.GetURL(shortURL)
	if err != nil {
		return &models.URL{FullURL: "/no_url/"} // get rick rolled :)
	}
	return url
}

// GetURLData returns a slice of URLData of the given URL
//
func (um *URLManager) GetURLData(url *models.URL) []*models.URLData {
	urlData, _ := um.db.GetURLData(url)
	return urlData
}

// TrackURLData stores short URL data if the short URL is assigned to a specific user
//
func (um *URLManager) TrackURLData(req *http.Request) {
	data := um.reqData.GetURLDataFromRequestData(req)
	data.ShortURL = req.URL.Path[1:]

	url, err := um.db.GetURL(data.ShortURL)
	if err != nil {
		return
	}

	_ = um.db.IncrementURLVisits(data.ShortURL)

	user, err := um.db.GetUser(&models.User{Email: url.UserEmail})
	if err == nil && user.Email != "" {
		_ = um.db.AddURLData(data)
	}
}

// HandleRickRoll redirects to Rick Astley's - Never Gonna Give You Up YT Video, perfect RickRolling :)
// GET /no_url/
//
func (um *URLManager) HandleRickRoll(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusFound)
}

// HandleRemoveURL removes the given url from the database
// DELETE /remove/?short=shortHandler
//
func (um *URLManager) HandleRemoveURL(res http.ResponseWriter, req *http.Request) {
	um.removeURL(req.URL.Query()["short"][0], req)
}

// removeURL removes a corresponding short URL from the database
//
func (um *URLManager) removeURL(shortURL string, request *http.Request) {
	user := um.userManager.GetUserFromRequest(request)
	url, err := um.db.GetURL(shortURL)

	if err == nil && user.Email == url.UserEmail {
		_ = um.db.RemoveURLData(url)
		_ = um.db.RemoveURL(url)
		_ = um.db.DecrementUserURLsCount(user)
	}
}
