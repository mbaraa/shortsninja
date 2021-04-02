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

// GetPageByName returns a handler function depending on page name
func (ui *UIManager) GetPageHandlerByName(pageName string) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		ui.renderPageFromUserIP(pageName, res, req)
	}
}

// HandleTracking renders the URLs tracking page of a specific user
func (ui *UIManager) HandleTracking(res http.ResponseWriter, req *http.Request) {
	user := ui.userManager.GetUserFromIP(req)

	urls := make([]*models.URL, 1)
	if user.Email != "" {
		urls = ui.userManager.GetURLsOfUser(user)
		var urlData []*models.URLData

		for _, url := range urls {
			urlData = ui.urlManager.GetURLData(url)
			url.Visits = len(urlData)
		}

	}

	// no error to be handled since it's being called by the router only :)
	_ = ui.templates.ExecuteTemplate(res, "tracking", map[string]interface{}{
		"Avatar":  user.Avatar,
		"Email":   user.Email,
		"FontB64": ui.conf.Font,
		"URLs":    urls,
	})
}

// HandleURLDataTracking renders the URLs tracking page of a specific user
func (ui *UIManager) HandleURLDataTracking(res http.ResponseWriter, req *http.Request) {
	user := ui.userManager.GetUserFromIP(req)

	urlData := make([]*models.URLData, 1)

	if shortURL := req.URL.Query()["short"]; user.Email != "" && shortURL != nil {
		url := ui.urlManager.GetURL(shortURL[0])
		if user.Email != url.UserEmail {
			goto ignoreData
		}
		urlData = ui.urlManager.GetURLData(url)

	}

ignoreData:
	// no error to be handled since it's being called by the router only :)
	_ = ui.templates.ExecuteTemplate(res, "url_data", map[string]interface{}{
		"Avatar":  user.Avatar,
		"Email":   user.Email,
		"FontB64": ui.conf.Font,
		"URLData": urlData,
	})
}

// HandleUserInfo renders the user info page
func (ui *UIManager) HandleUserInfo(res http.ResponseWriter, req *http.Request) {

	user := ui.userManager.GetUserFromIP(req)

	urls := make([]*models.URL, 1)
	if user.Email != "" {
		urls = ui.userManager.GetURLsOfUser(user)
	}
	numURLs := len(urls)

	// no error to be handled since it's being called by the router only :)
	_ = ui.templates.ExecuteTemplate(res, "login", map[string]interface{}{
		"Avatar":  user.Avatar,
		"Email":   user.Email,
		"Created": user.Created,
		"FontB64": ui.conf.Font,
		"NumURLs": numURLs,
	})

}

// renderPageFromUserIP generates the required web page with the given user's data
func (ui *UIManager) renderPageFromUserIP(pageName string, res http.ResponseWriter, req *http.Request) {
	user := ui.userManager.GetUserFromIP(req)

	// no error to be handled since it's being called by the router only :)
	_ = ui.templates.ExecuteTemplate(res, pageName, map[string]string{
		"Avatar":  user.Avatar,
		"Email":   user.Email,
		"FontB64": ui.conf.Font,
	})
}
