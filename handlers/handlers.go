package handlers

import (
	"github.com/baraa-almasri/useless/songs"
	"io/ioutil"
	"net/http"
)

func getFullURL(shortURL string) string {
	url, _ := ioutil.ReadFile("./urls/" + shortURL)
	return string(url)
}

func handle_play_meme_song(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, songs.NewMemeSongs().GetRandomSong(), http.StatusFound)
}

func handle_rick_roll(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusFound)
}
