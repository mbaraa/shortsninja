package controllers

import (
	"context"
	"encoding/json"
	"github.com/mbaraa/shortsninja/config"
	"github.com/mbaraa/shortsninja/models"
	"github.com/mbaraa/shortsninja/utils"
	"github.com/mbaraa/useless"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"time"
)

// GoogleLogin holds google login handlers
type GoogleLogin struct {
	googleOauthConfig *oauth2.Config
	state             string
	randomizer        *useless.RandASCII
	config            *config.Config
	reqData           *RequestDataManager
	db                models.Database
}

// NewGoogleLogin returns a new GoogleLogin instance
func NewGoogleLogin(randomizer *useless.RandASCII, config *config.Config,
	requestDataManager *RequestDataManager, dbManager models.Database) *GoogleLogin {
	return &GoogleLogin{
		googleOauthConfig: &oauth2.Config{
			RedirectURL:  config.GoogleCallbackHandler,
			ClientID:     config.GoogleClientID,
			ClientSecret: config.GoogleClientSecret,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
		state:      "",
		randomizer: randomizer,
		config:     config,
		reqData:    requestDataManager,
		db:         dbManager,
	}
}

// LoginWithGoogle handles login using google authentication
func (g *GoogleLogin) LoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	g.state = g.randomizer.GetRandomAlphanumString(32)
	url := g.googleOauthConfig.AuthCodeURL(g.state)
	http.Redirect(w, r, url, http.StatusFound)
}

// HandleCallback is called when authenticating with google
func (g *GoogleLogin) HandleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != g.state {
		http.Redirect(w, r, "/user_info/", http.StatusFound)
		return
	}

	token, err := g.googleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		http.Redirect(w, r, "/user_info/", http.StatusFound)
		return
	}

	dataResponse, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		http.Redirect(w, r, "/user_info/", http.StatusFound)
		return
	}
	defer dataResponse.Body.Close()

	data := make(map[string]interface{})
	json.NewDecoder(dataResponse.Body).Decode(&data)

	token1 := utils.NewUniqueID(g.randomizer).GetUniqueString(27)
	expireTime := time.Now().AddDate(0, 1, 0)

	_ = g.db.AddSession(&models.Session{
		UserEmail: data["email"].(string),
		Token:     token1,
		ExpiresAt: expireTime,
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token1,
		Expires: expireTime,
		Path:    "/",
	})

	_ = g.db.AddUser(&models.User{
		Email:  data["email"].(string),
		Avatar: data["picture"].(string),
	})

	// go back to the home page
	http.Redirect(w, r, "/", 308)
}
