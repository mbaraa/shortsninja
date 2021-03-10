package handlers

import (
	"github.com/baraa-almasri/shortsninja/globals"
	"io/ioutil"
	"os"
)

func getFullURL(shortURL string) string {
	url, err := ioutil.ReadFile("./urls/" + shortURL)
	if err != nil {
		return "/play_meme_song/"
	}
	return string(url)
}

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
