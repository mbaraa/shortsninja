package controllers

import (
	"github.com/baraa-almasri/shortsninja/config"
	"github.com/baraa-almasri/shortsninja/models"
	"github.com/baraa-almasri/shortsninja/utils"
	"net/http"
	"time"
)

// AdminController represents admin functionalities
//
type AdminController struct {
	db   models.Database
	conf *config.Config
	uniq *utils.UniqueID
	ui   *UIManager
}

// NewAdminController returns a new AdminController instance
//
func NewAdminController(dbManager models.Database, conf *config.Config,
	uniqueID *utils.UniqueID, uiManager *UIManager) *AdminController {
	return &AdminController{
		db:   dbManager,
		conf: conf,
		uniq: uniqueID,
		ui:   uiManager,
	}
}

// Login logs in the admin user :)
// GET /admin/
//
func (admin *AdminController) Login(res http.ResponseWriter, req *http.Request) {
	if admin.checkToken(req) {
		http.Redirect(res, req, "/admin/users/", 302)
		return
	}
	admin.ui.renderPageFromUserIP("admin", res, admin.ui.getBasicUserData(new(models.User)))
}

// AuthenticateAdmin checks admin credentials
// POST /admin/auth/
//
func (admin *AdminController) AuthenticateAdmin(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		res.WriteHeader(400)
		return
	}

	if req.PostFormValue("password") == admin.conf.AdminPassword {
		token := admin.uniq.GetUniqueString(27)
		expireTime := time.Now().Add(time.Minute * 30)

		http.SetCookie(res, &http.Cookie{
			Name:    "_admin_token",
			Value:   token,
			Path:    "/admin/",
			Expires: expireTime,
		})

		_ = admin.db.AddSession(&models.Session{
			UserEmail: "admin",
			Token:     token,
			ExpiresAt: expireTime,
		})
	}

	// go to main admin page anyway
	http.Redirect(res, req, "/admin/", 302)
}

// ViewUsers lists all users in the database
// GET /admin/users/
//
func (admin *AdminController) ViewUsers(res http.ResponseWriter, req *http.Request) {
	if admin.checkToken(req) {
		users, err := admin.db.GetUsers()
		if err != nil {
			users = []*models.User{}
		}
		admin.ui.renderPageFromUserIP("admin.users", res,
			mergeMaps(admin.ui.getBasicUserData(new(models.User)), map[string]interface{}{
				"Users":      users,
				"UsersCount": len(users),
			}))
		return
	}
	http.Redirect(res, req, "/admin/", 302)
}

// ViewURLs lists all short URLs in the database
// GET /admin/urls/
//
func (admin *AdminController) ViewURLs(res http.ResponseWriter, req *http.Request) {
	if admin.checkToken(req) {
		urls, err := admin.db.GetAllURLs()
		if err != nil {
			urls = []*models.URL{}
		}
		admin.ui.renderPageFromUserIP("admin.urls", res,
			mergeMaps(admin.ui.getBasicUserData(new(models.User)), map[string]interface{}{
				"URLs":      urls,
				"URLsCount": len(urls),
			}))
		return
	}
	http.Redirect(res, req, "/admin/", 302)
}

// ViewSessions lists all active sessions in the database
// GET /admin/sessions/
//
func (admin *AdminController) ViewSessions(res http.ResponseWriter, req *http.Request) {
	if admin.checkToken(req) {
		sessions, err := admin.db.GetSessions()
		if err != nil {
			sessions = []*models.Session{}
		}
		admin.ui.renderPageFromUserIP("admin.sessions", res,
			mergeMaps(admin.ui.getBasicUserData(new(models.User)), map[string]interface{}{
				"Sessions":      sessions,
				"SessionsCount": len(sessions),
			}))
		return
	}
	http.Redirect(res, req, "/admin/", 302)
}

// Logout deletes the current admin session
// GET /admin/logout/
//
func (admin *AdminController) Logout(res http.ResponseWriter, req *http.Request) {
	token, err := req.Cookie("_admin_token")
	if err == nil {
		_ = admin.db.RemoveSession(&models.Session{Token: token.Value})
	}
	http.Redirect(res, req, "/", 302)
}

// RemoveUser deletes a given user from the database
// DELETE /admin/users/remove/?user=someUser
//
func (admin *AdminController) RemoveUser(res http.ResponseWriter, req *http.Request) {
	if userEmail, ok := req.URL.Query()["user"]; ok && admin.checkToken(req) {
		_ = admin.db.RemoveUser(&models.User{Email: userEmail[0]})
	}
}

// RemoveURL deletes a given URL from the database
// DELETE /admin/urls/remove/?url=shortURL
//
func (admin *AdminController) RemoveURL(res http.ResponseWriter, req *http.Request) {
	if short, ok := req.URL.Query()["url"]; ok && admin.checkToken(req) {
		_ = admin.db.RemoveURL(&models.URL{Short: short[0]})
	}
}

// RemoveSession deletes a session from the database
// DELETE /admin/sessions/remove/?session=sessionToken
//
func (admin *AdminController) RemoveSession(res http.ResponseWriter, req *http.Request) {
	if token, ok := req.URL.Query()["session"]; ok && admin.checkToken(req) {
		_ = admin.db.RemoveSession(&models.Session{Token: token[0]})
	}
}

func (admin *AdminController) checkToken(req *http.Request) bool {
	token, err := req.Cookie("_admin_token")
	if err != nil {
		return false
	}

	session, _ := admin.db.GetSession(&models.Session{
		Token: token.Value,
	})

	return session != nil
}
