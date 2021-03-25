package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/baraa-almasri/shortsninja/globals"
	"github.com/baraa-almasri/shortsninja/models"
	"io/ioutil"
	"net/http"
)

// AddURL adds a URL to the short urls list and returns the assigned short url
func AddURL(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()["url"][0]
	shortURL := createAndUpdate(url, &models.User{})

	resp := make(map[string]interface{})
	resp["url"] = url
	resp["short"] = "shorts.ninja/" + shortURL
	resp["valid_url"] = isURLValid(url)

	_ = json.NewEncoder(w).Encode(resp)
}

// GetURL redirects to the full URL from the given shortURL
// if no url is assigned or the short URL is not valid it rick rolls the caller :)
func GetURL(w http.ResponseWriter, r *http.Request) {
	if !isShortURLValid(r.URL.Path[1:]) {
		http.Redirect(w, r, "/no_url/", http.StatusFound)
		return
	}

	url := getFullURL(r.URL.Path[1:])
	data := getRequestData(r)
	data.ShortURL = r.URL.Path[1:]

	_ = globals.DBManager.AddURLData(data)

	http.Redirect(w, r, url, http.StatusFound)
}

// RickRoll redirects to Rick Astley's - Never Gonna Give You Up YT Video, perfect RickRolling :)
func RickRoll(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusFound)
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	_ = globals.Templates.ExecuteTemplate(w, "shorten", getDummyUser())
}

func HandleTracking(w http.ResponseWriter, r *http.Request) {
	_ = globals.Templates.ExecuteTemplate(w, "tracking", getDummyUser())
}

func HandleAbout(w http.ResponseWriter, r *http.Request) {
	_ = globals.Templates.ExecuteTemplate(w, "about", getDummyUser())
}

func HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	_ = globals.Templates.ExecuteTemplate(w, "tracking", getDummyUser())
}

// TODO
// complete login and the other boiz :)
func Login(w http.ResponseWriter, r *http.Request) {
	token := &http.Cookie{
		Name:  "token",
		Value: "jherbvjlhr",
	}
	http.SetCookie(w, token)
	globals.Templates.ExecuteTemplate(w, "login", nil)
	//tmpl.Execute(w, data))
	/*globals.Templates.ExecuteTemplate(w, "login.html", map[string]interface{}{
		"User": "Blyat",
	})*/
}

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
		fmt.Println("sup from state")
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

	content, err := ioutil.ReadAll(dataResponse.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("couldn't parse response!"))
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	w.Write(content)
}

func CheckSession(w http.ResponseWriter, r *http.Request) {
	println("cookie: ", r.FormValue("token"))
	if r.FormValue("token") != "" {
		http.Redirect(w, r, "/no_url/", http.StatusFound)
		return
		w.Write([]byte("https://google.com"))
		return
	}
	w.Write([]byte(""))
}

func Signup(w http.ResponseWriter, r *http.Request) {

}
