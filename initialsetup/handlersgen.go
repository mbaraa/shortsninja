package initialsetup

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// GenerateMuxFileUsingURLsFile generates a multiplexer file from a given short URLs text file
// returns the short URLs slice for later use, and an occurred error{io}
//
// TODO
// reject non-valid short URLs file
func GenerateMuxFileUsingURLsFile(urlsFile *os.File) ([]string, error) {
	rawURLs, err := ioutil.ReadAll(urlsFile)
	if err != nil {
		return nil, err
	}

	urls := strings.Split(string(rawURLs), "\n")
	err = GenerateMuxFile(urls)
	if err != nil {
		return nil, err
	}

	// happily ever after
	return urls, nil
}

// GenerateMuxFile generates urls multiplexer from a given short URLs string slice :)
// returns an occurred error{io}
// TODO
// append urls handlers instead of overriding!
func GenerateMuxFile(urls []string) error {
	// error ignored since it's only {file exists}
	_ = os.Mkdir("./mux", 0755)

	f, err := os.Create("./mux/mux.go")
	defer f.Close()
	if err != nil {
		return err
	}

	muxFileCont := getMuxPrefix()

	for _, shortURL := range urls {
		muxFileCont += generateHandlerStatement(shortURL, generateHandlerFunction(shortURL))
	}
	muxFileCont += "\n}"

	_, err = f.Write([]byte(muxFileCont))
	if err != nil {
		return err
	}

	// happily ever after
	return nil
}

// generateHandlerFunction returns a handler function of the given short URL that will be used in the multiplexer
func generateHandlerFunction(shortURL string) string {
	return fmt.Sprintf(`func (w http.ResponseWriter, r *http.Request) {
		url := getFullURL("%s")
		http.Redirect(w, r, url, http.StatusFound)
	}`, shortURL)
}

// generateHandlerStatement returns a handling statement for a given short url
func generateHandlerStatement(shortURL, handlerFunc string) string {
	return fmt.Sprintf("\n\tmux.HandleFunc(\"/%s/\", %s)", shortURL, handlerFunc)
}

// getMuxPrefix returns multiplexer file content before the handling statements
// much wow!
func getMuxPrefix() string {
	return `// auto generated file
// generated at server's initialization!
package handlers 

import (
	"github.com/baraa-almasri/useless/songs"
	"io/ioutil"
	"net/http"
)

var (
	mux = http.NewServeMux()
)

func GetMux() *http.ServeMux {
	bunchHandle()
	return mux
}

func getFullURL(shortURL string) string {
	url, err := ioutil.ReadFile("./urls/" + shortURL)
	if err != nil {
		return "/play_meme_song/"
	}
	return string(url)
}

func bunchHandle() {
	mux.HandleFunc("/play_meme_song/", func (w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, songs.NewMemeSongs().GetRandomSong(), http.StatusFound)
	})

	mux.HandleFunc("/no_url/", func (w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusFound)
	})
`
}
