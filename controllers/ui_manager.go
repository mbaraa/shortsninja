package controllers

import (
	"github.com/baraa-almasri/shortsninja/config"
	"github.com/baraa-almasri/shortsninja/models"
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

// HandleTracking renders the URLs tracking page of a specific user
func (ui *UIManager) HandleTracking(res http.ResponseWriter, req *http.Request) {
	user := ui.userMan.GetUserFromToken(req)
	ui.userMan.SetToken(res, req)

	urls := make([]*models.URL, 1)
	if user.Email != "" {
		urls = ui.userMan.GetURLsOfUser(user)
	}

	// no error to be handled since it's being called by the router only :)
	_ = ui.templates.ExecuteTemplate(res, "tracking", map[string]interface{}{
		"Avatar":  user.Avatar,
		"Email":   user.Email,
		"FontB64": ui.conf.Font,
		"URLs":    urls,
	})
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
