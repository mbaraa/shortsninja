package routes

import (
	"github.com/baraa-almasri/shortsninja/config"
	"github.com/baraa-almasri/shortsninja/controllers"
	"github.com/baraa-almasri/shortsninja/models"
	"github.com/baraa-almasri/shortsninja/utils"
	"github.com/baraa-almasri/useless"
	"github.com/gorilla/mux"
	"html/template"
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
	urlValidator       *utils.URLValidator
	multiplexer        *mux.Router
	admin              *controllers.AdminController
}

// NewRouter returns a new Router instance
func NewRouter(dbManager models.Database, templates *template.Template, config *config.Config) *Router {
	randomizer := useless.NewRandASCII()

	requestDataManager := controllers.NewRequestDataManager(config, dbManager)
	userManager := controllers.NewUserManager(dbManager, requestDataManager)
	urlManager := controllers.NewURLManager(
		utils.NewURLValidator(), requestDataManager, userManager, randomizer, dbManager)
	uiManager := controllers.NewUIManager(userManager, urlManager, templates, config)
	urlValidator := utils.NewURLValidator()
	googleLogin := controllers.NewGoogleLogin(randomizer, config, requestDataManager, dbManager)
	admin := controllers.NewAdminController(dbManager, config, utils.NewUniqueID(randomizer), uiManager)

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
		admin:              admin,
	}
}

// GetRoutes returns a gorilla mux router with the wanted routes
func (router *Router) GetRoutes() *mux.Router {
	router.handleURLOps()
	router.handleUI()
	router.handleUserOps()
	router.handleAdminOps()

	return router.multiplexer
}

func (router *Router) handleURLOps() {
	router.multiplexer.HandleFunc("/shorten/", router.urlManager.HandleCreateShortURL).Methods("GET")
	router.multiplexer.HandleFunc("/no_url/", router.urlManager.HandleRickRoll).Methods("GET")
	router.multiplexer.HandleFunc("/{[A-Z;0-9;a-z]+}", router.urlManager.HandleGetFullURL).Methods("GET")
	router.multiplexer.HandleFunc("/remove/", router.urlManager.HandleRemoveURL).Methods("DELETE")
}

func (router *Router) handleUI() {
	router.multiplexer.HandleFunc("/", router.uiManager.GetPageHandlerByName("shorten")).Methods("GET")
	router.multiplexer.HandleFunc("/about/", router.uiManager.GetPageHandlerByName("about")).Methods("GET")
	router.multiplexer.HandleFunc("/tracking/", router.uiManager.HandleTracking).Methods("GET")
	router.multiplexer.HandleFunc("/user_info/", router.uiManager.HandleUserInfo).Methods("GET")
	router.multiplexer.HandleFunc("/url_data/", router.uiManager.HandleURLDataTracking).Methods("GET")
	router.multiplexer.HandleFunc("/admin/", router.admin.Login).Methods("GET")
}

func (router *Router) handleUserOps() {
	router.multiplexer.HandleFunc("/login/", router.googleLoginManager.LoginWithGoogle).Methods("GET")
	router.multiplexer.HandleFunc("/login_callback/", router.googleLoginManager.HandleCallback).Methods("GET")
	router.multiplexer.HandleFunc("/logout/", router.userManager.Logout).Methods("GET")
	router.multiplexer.HandleFunc("/GTFO/", router.userManager.GTFO).Methods("GET")
}

func (router *Router) handleAdminOps() {
	router.multiplexer.HandleFunc("/admin/auth/", router.admin.AuthenticateAdmin).Methods("POST")
	router.multiplexer.HandleFunc("/admin/users/", router.admin.ViewUsers).Methods("GET")
	router.multiplexer.HandleFunc("/admin/users/remove/", router.admin.RemoveUser).Methods("DELETE")
	router.multiplexer.HandleFunc("/admin/urls/", router.admin.ViewURLs).Methods("GET")
	router.multiplexer.HandleFunc("/admin/urls/remove/", router.admin.RemoveURL).Methods("DELETE")
	router.multiplexer.HandleFunc("/admin/sessions/", router.admin.ViewSessions).Methods("GET")
	router.multiplexer.HandleFunc("/admin/sessions/remove/", router.admin.RemoveSession).Methods("DELETE")
	router.multiplexer.HandleFunc("/admin/logout/", router.admin.Logout).Methods("GET")
}
