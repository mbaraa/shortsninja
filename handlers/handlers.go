package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/baraa-almasri/shortsninja/globals"
	"github.com/baraa-almasri/shortsninja/models"
	"io/ioutil"
	"net/http"
	"regexp"
)

// AddURL adds a URL to the short urls list and returns the assigned short url
func AddURL(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()["url"][0]
	user := r.URL.Query()["user"][0]

	short := createAndUpdate(url, &models.User{
		Email: user,
	})

	_ = json.NewEncoder(w).Encode(&models.URL{
		Short:   short,
		FullURL: url,
	})
}

// GetURL redirects to the full URL from the given shortURL
// if no url is assigned it plays a meme song :)
func GetURL(w http.ResponseWriter, r *http.Request) {
	// check if the short URL is actually valid!
	shortURLPattern, _ := regexp.Compile("[A-Z;0-9;a-z]{5}")
	shortURLPattern.FindString(r.URL.Path[1:])
	if shortURLPattern.FindString(r.URL.Path[1:]) != r.URL.Path[1:] {
		http.Redirect(w, r, "/play_meme_song/", http.StatusFound)
		return
	}

	url := getURLData(r.URL.Path[1:])
	data := getRequestData(r)
	data["url"] = url

	http.Redirect(w, r, url, http.StatusFound)
}

// PlayMemeSong plays a random meme song, wow!
func PlayMemeSong(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, memes.GetRandomSong(), http.StatusFound)
}

// RickRoll redirects to Rick Astley's - Never Gonna Give You Up YT Video, perfect RickRolling :)
func RickRoll(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusFound)
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	globals.Templates.ExecuteTemplate(w, "shorten", &models.User{
		Email:  "hexagon16.rpm@gmail.com",
		Avatar: "https://i1.wp.com/tech-ish.com/wp-content/uploads/2014/10/Google-Logo.jpg?fit=1575%2C1575&ssl=1",
	})
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
