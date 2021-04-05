package routes

import (
	"encoding/json"
	"github.com/baraa-almasri/shortsninja/config"
	"github.com/baraa-almasri/shortsninja/controllers"
	"github.com/baraa-almasri/shortsninja/models"
	"github.com/baraa-almasri/useless"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

// Router represents the program's multiplexer
type Router struct {
	dbMan              models.Database
	templates          *template.Template
	conf               *config.Config
	requestDataManager *controllers.RequestDataManager
	userManager        *controllers.UserManager
	uiManager          *controllers.UIManager
	urlManager         *controllers.URLManager
	googleLoginManager *controllers.GoogleLogin
	urlValidator       *controllers.URLValidator
	multiplexer        *mux.Router
}

// NewRouter returns a new Router instance
func NewRouter(dbManager models.Database, templates *template.Template, config *config.Config) *Router {
	randomizer := useless.NewRandASCII()

	requestDataManager := controllers.NewRequestDataManager(config, dbManager)
	userManager := controllers.NewUserManager(dbManager, requestDataManager)
	urlManager := controllers.NewURLManager(
		controllers.NewURLValidator(), requestDataManager, userManager, randomizer, dbManager,
	)
	uiManager := controllers.NewUIManager(userManager, urlManager, templates, config)
	urlValidator := controllers.NewURLValidator()
	googleLogin := controllers.NewGoogleLogin(randomizer, config, requestDataManager, dbManager)

	return &Router{
		dbMan:              dbManager,
		templates:          templates,
		conf:               config,
		requestDataManager: requestDataManager,
		userManager:        userManager,
		uiManager:          uiManager,
		urlManager:         urlManager,
		googleLoginManager: googleLogin,
		urlValidator:       urlValidator,
		multiplexer:        mux.NewRouter(),
	}
}

// GetRoutes returns a gorilla mux router with the wanted routes
func (router *Router) GetRoutes() *mux.Router {
	return router.handleRoutes()
}

func (router *Router) handleRoutes() *mux.Router {
	router.handleURLOps()
	router.handleUI()
	router.handleUserOps()

	return router.multiplexer
}

// createShortURL creates a short url for the given url
// GET /shorten/?url=http://someurl.com
func (router *Router) createShortURL(res http.ResponseWriter, req *http.Request) {
	user := router.userManager.GetUserFromIP(req)
	url := req.URL.Query()["url"][0]

	resp := make(map[string]interface{})
	if resp["valid_url"] = router.urlValidator.IsURLValid(url); resp["valid_url"].(bool) {
		shortURL := router.urlManager.CreateShortURL(url, user)

		resp["url"] = url
		resp["short"] = "shorts.ninja/" + shortURL
	}
	_ = json.NewEncoder(res).Encode(resp)
}

// getFullURL retrieves the original url of the given short url
// GET /{[A-Z;0-9;a-z]{4,5}}
func (router *Router) getFullURL(res http.ResponseWriter, req *http.Request) {
	if !router.urlValidator.IsShortURLValid(req.URL.Path[1:]) {
		http.Redirect(res, req, "/no_url/", http.StatusFound)
		return
	}

	url := router.urlManager.GetFullURL(req.URL.Path[1:])
	router.urlManager.TrackURLData(req)

	http.Redirect(res, req, url, http.StatusFound)
}

// removeURL removes the given url from the database
// DELETE /remove/?short=shortHandler
func (router *Router) removeURL(res http.ResponseWriter, req *http.Request) {
	router.urlManager.RemoveURL(req.URL.Query()["short"][0], req)
}

// getURLData loads a specific URL's data into a weird ass table
// GET /url_data/?short=shortURL
func (router *Router) getURLData(res http.ResponseWriter, req *http.Request) {
	router.uiManager.HandleURLDataTracking(res, req)
}

// rickRoll redirects to Rick Astley's - Never Gonna Give You Up YT Video, perfect RickRolling :)
// GET /no_url/
func (router *Router) rickRoll(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusFound)
}

func (router *Router) handleURLOps() {
	router.multiplexer.HandleFunc("/shorten/", router.createShortURL).Methods("GET")
	router.multiplexer.HandleFunc("/no_url/", router.rickRoll).Methods("GET")
	router.multiplexer.HandleFunc("/{[A-Z;0-9;a-z]{4,5}}", router.getFullURL).Methods("GET")
	router.multiplexer.HandleFunc("/remove/", router.removeURL).Methods("DELETE")
}

func (router *Router) handleUI() {
	router.multiplexer.HandleFunc("/", router.uiManager.GetPageHandlerByName("shorten")).Methods("GET")
	router.multiplexer.HandleFunc("/about/", router.uiManager.GetPageHandlerByName("about")).Methods("GET")
	router.multiplexer.HandleFunc("/tracking/", router.uiManager.HandleTracking).Methods("GET")
	router.multiplexer.HandleFunc("/user_info/", router.uiManager.HandleUserInfo).Methods("GET")
	router.multiplexer.HandleFunc("/url_data/", router.uiManager.HandleURLDataTracking).Methods("GET")
}

func (router *Router) handleUserOps() {
	router.multiplexer.HandleFunc("/login/", router.googleLoginManager.LoginWithGoogle).Methods("GET")
	router.multiplexer.HandleFunc("/login_callback/", router.googleLoginManager.HandleCallback).Methods("GET")
	router.multiplexer.HandleFunc("/logout/", router.userManager.Logout).Methods("GET")
}
