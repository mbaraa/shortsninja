package main

import (
	"github.com/baraa-almasri/shortsninja/config"
	"github.com/baraa-almasri/shortsninja/models"
	"github.com/baraa-almasri/shortsninja/routes"
	"github.com/rs/cors"
	"html/template"
	"log"
	"net/http"
)

// still testing :)
func main() {
	templates := template.Must(template.ParseGlob("./templates/*.html"))
	dbManager := models.NewSQLiteDB()
	conf := config.LoadConfig()

	mux := routes.NewRouter(dbManager, templates, conf)

	corsHandler := cors.Default().Handler(mux.GetRoutes())
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
