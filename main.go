package main

import (
	"github.com/baraa-almasri/shortsninja/globals"
	"github.com/baraa-almasri/shortsninja/handlers"
	"net/http"
	"os"
)

// still testing :)
func main() {
	/*
		op, _ := exec.Command("ls", "./urls").Output()
		fmt.Println(string(op))

		return*/
	f, _ := os.Open("./urls.txt")
	globals.LoadGlobals(f)
	mux := handlers.GetMux()
	mux.HandleFunc("/shorten/", handlers.AddURL)
	http.ListenAndServe(":8080", mux)

	return

	//initialsetup.GenerateMuxFileUsingURLsFile(f)
}
