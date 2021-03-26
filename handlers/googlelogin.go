package handlers

import (
	"context"
	"encoding/json"
	"github.com/baraa-almasri/shortsninja/globals"
	"github.com/baraa-almasri/shortsninja/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  globals.Config.GoogleCallbackHandler,
		ClientID:     globals.Config.GoogleClientID,
		ClientSecret: globals.Config.GoogleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	state = ""
)

// GoogleLogin handles login using google authentication
func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state = randomizer.GetRandomAlphanumString(32)
	url := googleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

// HandleCallback is called when authenticating with google
func HandleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != state {
		w.WriteHeader(http.StatusNotFound)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	token, err := googleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("couldn't get token!"))
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	dataResponse, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	defer dataResponse.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data := make(map[string]interface{})
	json.NewDecoder(dataResponse.Body).Decode(&data)

	token1 := randomizer.GetRandomAlphanumString(32)
	_ = globals.DBManager.AddSession(&models.Session{
		Token:     token1,
		IP:        r.Header.Get("X-FORWARDED-FOR"),
		UserEmail: data["email"].(string),
	})

	_ = globals.DBManager.AddUser(&models.User{
		Email:  data["email"].(string),
		Avatar: data["picture"].(string),
	})

	// go back to the home page with the user token
	http.Redirect(w, r, "/?token="+token1, http.StatusFound)
}
