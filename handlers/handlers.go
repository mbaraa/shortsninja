// not a generated file!
package handlers

import (
	"github.com/baraa-almasri/useless/songs"
	"net/http"
)

func AddURL(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()["url"][0]
	createAndUpdate(url)
}

func GetURL(w http.ResponseWriter, r *http.Request) {
	url := getFullURL(r.URL.Path)
	http.Redirect(w, r, url, http.StatusFound)
}

func PlayMemeSong(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, songs.NewMemeSongs().GetRandomSong(), http.StatusFound)
}

func RickRoll(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusFound)
}
