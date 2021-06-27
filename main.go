package main

import (
	"github.com/mbaraa/shortsninja/config"
	"github.com/mbaraa/shortsninja/models"
	"github.com/mbaraa/shortsninja/routes"
	"github.com/rs/cors"
	"html/template"
	"log"
	"net/http"
)

func main() {
	templates := template.Must(template.ParseGlob("./templates/*"))
	dbManager := models.NewSQLiteDB()
	conf := config.LoadConfig()

	mux := routes.NewRouter(dbManager, templates, conf)

	//go runMainHTTPSServer(mux)
	//runHTTPServer()

	runLocally(mux)
}

func runLocally(router *routes.Router) {
	corsHandler := cors.Default().Handler(router.GetRoutes())
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

func runMainHTTPSServer(router *routes.Router) {
	corsHandler := cors.Default().Handler(router.GetRoutes())
	log.Fatal(http.ListenAndServeTLS("", "./shorts.ninja.crt", "./shorts.ninja.key", corsHandler))
}

func runHTTPServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "https://shorts.ninja/"+req.URL.Path[1:], http.StatusFound)
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}
