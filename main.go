package main

import (
	"errors"
	"fmt"
	"github.com/baraa-almasri/shortsninja/globals"
	"github.com/baraa-almasri/shortsninja/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

// still testing :)
func main() {
	f, _ := os.Open("./urls.txt")
	urls, _, _ := globals.LoadGlobals(f)
	if urls[0] == "" {
		fmt.Println(errors.New("wow looks like `urls.txt` is empty!" +
			"\ngenerate using the admin panel(lol not ready)" +
			"\nor generate manually using `initialsetup.GenerateShortURLsFile()` or `initialsetup.UpdateShortURLsFile()`"))
	}

	m := mux.NewRouter()
	m.HandleFunc("/shorten", handlers.AddURL).Methods("GET")
	m.HandleFunc("/play_meme_song/", handlers.PlayMemeSong).Methods("GET")
	m.HandleFunc("/no_url", handlers.RickRoll).Methods("GET")
	m.HandleFunc("/{[A-Z;0-9;a-z]{4}}", handlers.GetURL).Methods("GET")

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("nothing to see here!"))
	})

	log.Fatal(http.ListenAndServe(":8080", m))
}
