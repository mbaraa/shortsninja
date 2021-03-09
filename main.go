package main

import (
	"github.com/baraa-almasri/shortsninja/handlers"
	"github.com/baraa-almasri/shortsninja/initialsetup"
	"net/http"
	"os"
)

// still testing :)
func main() {
	mux := handlers.GetMux()
	http.ListenAndServe(":8080", mux)

	return
	f, _ := os.Open("./urls.txt")
	defer f.Close()

	initialsetup.GenerateMuxFileUsingURLsFile(f)
}
