package handlers

import (
	"github.com/baraa-almasri/shortsninja/globals"
	"os"
)

func elementExistInArr(element string, slice []string) bool {
	for i := range slice {
		if slice[i] == element {
			return true
		}
	}

	return false
}

func createAndUpdate(url string) {
	for _, v := range globals.ShortURLs {
		if !elementExistInArr(v, globals.UsedShortURLs) {
			f, _ := os.Create("./urls/" + v)
			_, _ = f.Write([]byte(url))
			_ = f.Close()

			globals.UsedShortURLs = append(globals.UsedShortURLs, v)
			return
		}
	}
}
