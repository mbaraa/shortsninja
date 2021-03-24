package main

import (
	"github.com/baraa-almasri/shortsninja/db"
	"github.com/baraa-almasri/shortsninja/globals"
	"github.com/baraa-almasri/shortsninja/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"html/template"
	"log"
	"net/http"
)

// still testing :)
func main() {
	globals.Templates, _ = template.ParseGlob("./templates/*.html")
	globals.DBManager = db.NewSQLiteDB()

	m := mux.NewRouter()
	m.HandleFunc("/shorten/", handlers.AddURL).Methods("GET")
	m.HandleFunc("/play_meme_song/", handlers.PlayMemeSong).Methods("GET")
	m.HandleFunc("/no_url/", handlers.RickRoll).Methods("GET")
	m.HandleFunc("/{[A-Z;0-9;a-z]{5}}", handlers.GetURL).Methods("GET")

	m.HandleFunc("/", handlers.HandleHome).Methods("GET")
	//m.HandleFunc("/signup/", handlers.Signup).Methods("GET")
	//m.HandleFunc("/check_session/", handlers.CheckSession).Methods("GET")
	//m.HandleFunc("/login/", handlers.GoogleLogin).Methods("GET")
	//m.HandleFunc("/login_callback/", handlers.HandleCallback).Methods("GET")

	corsHandler := cors.Default().Handler(m)

	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
