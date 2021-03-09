package handlers

import (
	"net/http"
)

var (
	mux = http.NewServeMux()
)

func GetMux() *http.ServeMux {
	bunchHandle()
	return mux
}

func bunchHandle() {
	mux.HandleFunc("/play_meme_song/", handle_play_meme_song)
	mux.HandleFunc("/no_url/", handle_rick_roll)

}
