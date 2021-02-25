package handlers

import "net/http"

var mux = &http.ServeMux{}

func GetMux() *http.ServeMux {
	bunchHandle()
	return mux
}

func bunchHandle() {
	mux.HandleFunc("/dummy", handle_dummy)

}