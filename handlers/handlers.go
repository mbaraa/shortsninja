package handlers

import "net/http"

func handle_dummy(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<!DOCTYPE html><head><title>Redirecting...</title><script> window.open(\"https://google.com\", \"_top\"); </script></head></html>"))
}