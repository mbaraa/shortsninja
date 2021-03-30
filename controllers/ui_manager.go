package controllers

import (
	"github.com/baraa-almasri/shortsninja/config"
	"html/template"
	"net/http"
)

// UIManager holds the UI handlers
type UIManager struct {
	userMan   *UserManager
	templates *template.Template
	conf      *config.Config
}

// NewUIManager returns a new UIManager instance
func NewUIManager(userManager *UserManager, templates *template.Template, config *config.Config) *UIManager {
	return &UIManager{
		userMan:   userManager,
		templates: templates,
		conf:      config,
	}
}

// GetPageByName returns a handler function depending on page name
func (ui *UIManager) GetPageByName(pageName string) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		ui.renderPageFromSessionToken(pageName, res, req)
	}
}

// Deprecated
// use GetPageByName instead, since it's more generic and really awesome :)
//
// HandleHome renders the shortening page of a specific or an anonymous user
func (ui *UIManager) HandleHome(w http.ResponseWriter, r *http.Request) {
	ui.renderPageFromSessionToken("shorten", w, r)
}

// Deprecated
// use GetPageByName instead, since it's more generic and really awesome :)
//
// HandleTracking renders the URLs tracking page of a specific user
func (ui *UIManager) HandleTracking(w http.ResponseWriter, r *http.Request) {
	ui.renderPageFromSessionToken("tracking", w, r)
}

// Deprecated
// use GetPageByName instead, since it's more generic and really awesome :)
//
// HandleAbout renders the about page
func (ui *UIManager) HandleAbout(w http.ResponseWriter, r *http.Request) {
	ui.renderPageFromSessionToken("about", w, r)
}

// Deprecated
// use GetPageByName instead, since it's more generic and really awesome :)
//
// HandleUserInfo renders the user info page
func (ui *UIManager) HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	ui.renderPageFromSessionToken("login", w, r)
}

// renderPageFromSessionToken generates the required web page with the given user's data
func (ui *UIManager) renderPageFromSessionToken(pageName string, res http.ResponseWriter, req *http.Request) {
	user := ui.userMan.GetUserFromToken(req)
	ui.userMan.SetToken(res, req)

	// no error to be handled since it's being called by the router only :)
	_ = ui.templates.ExecuteTemplate(res, pageName, map[string]string{
		"Avatar":  user.Avatar,
		"Email":   user.Email,
		"FontB64": ui.conf.Font,
	})
}
