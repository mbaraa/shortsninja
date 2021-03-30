package controllers

import (
	"context"
	"encoding/json"
	"github.com/baraa-almasri/shortsninja/config"
	"github.com/baraa-almasri/shortsninja/models"
	"github.com/baraa-almasri/useless"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
)

// GoogleLogin holds google login handlers
type GoogleLogin struct {
	googleOauthConfig *oauth2.Config
	state             string
	randomizer        *useless.RandASCII
	config            *config.Config
	reqData           *RequestDataManager
	dbMan             models.Database
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
		dbMan:      dbManager,
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
		w.WriteHeader(http.StatusNotFound)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	token, err := g.googleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		w.Write([]byte("couldn't get token!"))
		http.Redirect(w, r, "/user_info/", http.StatusFound)
		return
	}

	dataResponse, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	defer dataResponse.Body.Close()
	if err != nil {
		http.Redirect(w, r, "/user_info/", http.StatusFound)
		return
	}

	data := make(map[string]interface{})
	json.NewDecoder(dataResponse.Body).Decode(&data)

	token1 := g.randomizer.GetRandomAlphanumString(32)
	callerIP := g.reqData.GetIP(r)
	_ = g.dbMan.AddSession(&models.Session{
		Token:     token1,
		IP:        callerIP,
		UserEmail: data["email"].(string),
	})

	_ = g.dbMan.AddUser(&models.User{
		Email:  data["email"].(string),
		Avatar: data["picture"].(string),
	})

	// go back to the home page with the user token
	http.Redirect(w, r, "/?token="+token1, http.StatusFound)
}
