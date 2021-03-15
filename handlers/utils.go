package handlers

import (
	"github.com/baraa-almasri/shortsninja/globals"
	"github.com/baraa-almasri/useless/songs"
	"io/ioutil"
	"os"
	"strings"
)

var (
	memes = songs.NewMemeSongs()
)

// getFullURL retrieves a URL from its corresponding shortURL file
func getFullURL(shortURL string) string {
	url, err := ioutil.ReadFile("./urls/" + shortURL)
	if err != nil || strings.Contains(shortURL, ".") { // wow, much security!
		return "/play_meme_song/"
	}
	return string(url)
}

// createAndUpdate creates a new url file that doesn't exist in the usedURLs slice and updates the usedURLs slice
// and returns the assigned short URL
func createAndUpdate(url string) string {
	for _, v := range globals.ShortURLs {
		if !globals.IsShortURLUsed(v) {
			f, _ := os.Create("./urls/" + v)
			_, _ = f.Write([]byte(url))
			_ = f.Close()

			globals.UsedShortURLs = append(globals.UsedShortURLs, v)
			return v
		}
	}

	return ""
}
