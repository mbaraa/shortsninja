package main

import (
	"github.com/baraa-almasri/shortsninja/handlers"
	"net/http"
)

func main() {
	mux := handlers.GetMux()
	http.ListenAndServe(":8080", mux)
}
