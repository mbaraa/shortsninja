// not a generated file!
package handlers

import (
	"net/http"
)

func AddURL(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()["url"][0]
	createAndUpdate(url)
}
