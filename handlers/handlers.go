package handlers

import (
	"net/http"
)

// AddURL adds a URL to the short urls list and returns the assigned short url
func AddURL(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()["url"][0]
	short := createAndUpdate(url, "admin")
	w.Write([]byte("http://shorts.ninja/" + short))
}

// GetURL redirects to the full URL from the given shortURL
// if no url is assigned it plays a meme song :)
func GetURL(w http.ResponseWriter, r *http.Request) {
	url, user := getURLData(r.URL.Path[1:])
	data := getRequestData(r)
	data["url"] = url

	appendData(r.URL.Path[1:], user, data)

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
