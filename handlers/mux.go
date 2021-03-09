// test file!
package handlers

import (
	"github.com/baraa-almasri/useless/songs"
	"io/ioutil"
	"net/http"
)

var (
	mux = http.NewServeMux()
)

func GetMux() *http.ServeMux {
	bunchHandle()
	return mux
}

func getFullURL(shortURL string) string {
	url, err := ioutil.ReadFile("./urls/" + shortURL)
	if err != nil {
		return "/play_meme_song/"
	}
	return string(url)
}

func bunchHandle() {
	mux.HandleFunc("/play_meme_song/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, songs.NewMemeSongs().GetRandomSong(), http.StatusFound)
	})

	mux.HandleFunc("/no_url/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusFound)
	})

}
