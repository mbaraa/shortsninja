package handlers

import (
	"encoding/json"
	"github.com/baraa-almasri/shortsninja/globals"
	"github.com/baraa-almasri/shortsninja/models"
	"net/http"
	"strings"
)

// AddURL adds a URL to the short urls list and returns the assigned short url
func AddURL(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	if r.URL.Query()["token"] != nil {
		callerIP := getIP(r)
		callerIP = callerIP[:strings.Index(callerIP, ":")]
		user = getUser(r.URL.Query()["token"][0], callerIP)
	}
	url := r.URL.Query()["url"][0]

	resp := make(map[string]interface{})
	if resp["valid_url"] = isURLValid(url); resp["valid_url"].(bool) {
		shortURL := createAndUpdate(url, user)

		resp["url"] = url
		resp["short"] = "shorts.ninja/" + shortURL
	}
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

// HandleHome renders the shortening page of a specific or an anonymous user
func HandleHome(w http.ResponseWriter, r *http.Request) {
	ip, token := getIPAndToken(w, r)
	renderPageFromSessionToken("shorten", token, ip, w, r)
}

// HandleTracking renders the URLs tracking page of a specific user
func HandleTracking(w http.ResponseWriter, r *http.Request) {
	ip, token := getIPAndToken(w, r)
	renderPageFromSessionToken("tracking", token, ip, w, r)
}

// HandleAbout renders the about page
func HandleAbout(w http.ResponseWriter, r *http.Request) {
	ip, token := getIPAndToken(w, r)
	renderPageFromSessionToken("about", token, ip, w, r)
}

// HandleUserInfo renders the user info page
func HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	ip, token := getIPAndToken(w, r)
	renderPageFromSessionToken("login", token, ip, w, r)
}

// RickRoll redirects to Rick Astley's - Never Gonna Give You Up YT Video, perfect RickRolling :)
func RickRoll(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusFound)
}
