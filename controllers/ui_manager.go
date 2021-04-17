package controllers

import (
	"github.com/baraa-almasri/shortsninja/config"
	"github.com/baraa-almasri/shortsninja/models"
	"html/template"
	"net/http"
)

// UIManager holds the UI handlers
type UIManager struct {
	userManager *UserManager
	urlManager  *URLManager
	templates   *template.Template
	conf        *config.Config
}

// NewUIManager returns a new UIManager instance
func NewUIManager(userManager *UserManager, urlManager *URLManager,
	templates *template.Template, config *config.Config) *UIManager {

	return &UIManager{
		userManager: userManager,
		urlManager:  urlManager,
		templates:   templates,
		conf:        config,
	}
}

// GetPageHandlerByName returns a handler function depending on page name
func (ui *UIManager) GetPageHandlerByName(pageName string) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		user := ui.userManager.GetUserFromRequest(req)
		ui.renderPageFromUserIP(pageName, res, ui.getBasicUserData(user))
	}
}

// HandleTracking renders the URLs tracking page of a specific user
func (ui *UIManager) HandleTracking(res http.ResponseWriter, req *http.Request) {
	user := ui.userManager.GetUserFromRequest(req)
	urls := ui.userManager.GetURLsOfUser(user)

	ui.renderPageFromUserIP("tracking", res, mergeMaps(ui.getBasicUserData(user), map[string]interface{}{
		"URLs": urls,
	}))
}

// HandleURLDataTracking renders the URLs tracking page of a specific user
func (ui *UIManager) HandleURLDataTracking(res http.ResponseWriter, req *http.Request) {
	user := ui.userManager.GetUserFromRequest(req)
	var urlData []*models.URLData

	if shortURL := req.URL.Query()["short"]; user != nil && shortURL != nil {
		url := ui.urlManager.getFullURL(shortURL[0])
		if user.Email != url.UserEmail {
			goto ignoreData
		}
		urlData = ui.urlManager.GetURLData(url)
	}

ignoreData:
	ui.renderPageFromUserIP("url_data", res, mergeMaps(ui.getBasicUserData(user), map[string]interface{}{
		"URLData": urlData,
	}))
}

// HandleUserInfo renders the user info page
func (ui *UIManager) HandleUserInfo(res http.ResponseWriter, req *http.Request) {
	user := ui.userManager.GetUserFromRequest(req)

	ui.renderPageFromUserIP("login", res, mergeMaps(ui.getBasicUserData(user), map[string]interface{}{
		"Created": user.Created,
		"NumURLs": user.NumURLs,
	}))
}

// renderPageFromUserIP generates the required web page with the given user's data
func (ui *UIManager) renderPageFromUserIP(pageName string, res http.ResponseWriter, data map[string]interface{}) {
	_ = ui.templates.ExecuteTemplate(res, pageName, data)
}

func (ui *UIManager) getBasicUserData(user *models.User) map[string]interface{} {
	return map[string]interface{}{
		"Avatar":  user.Avatar,
		"Email":   user.Email,
		"FontB64": ui.conf.Font,
	}
}

func mergeMaps(src map[string]interface{}, dist map[string]interface{}) map[string]interface{} {
	for key, value := range src {
		dist[key] = value
	}
	return dist
}
