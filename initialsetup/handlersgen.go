package initialsetup

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// TODO
// write more docs
// I swear I will write :)

func GenerateAllUsingURLsFile(urlsFile *os.File) error {
	rawURLs, err := ioutil.ReadAll(urlsFile)
	if err != nil {
		return err
	}

	urls := strings.Split(string(rawURLs), "\n")
	err = GenerateAll(urls)
	if err != nil {
		return err
	}

	// happily ever after
	return nil
}

// GenerateAll generates url handlers and a proper multiplexer
func GenerateAll(urls []string) error {
	err := GenerateMultiplexerFile(urls)
	if err != nil {
		return err
	}
	err = GenerateHandlerFunctionsFile(urls)
	if err != nil {
		return err
	}

	// happily ever after
	return nil
}

// GenerateHandlerFunctionsFile generates a handlers' file
func GenerateHandlerFunctionsFile(urls []string) error {
	f, err := os.Create("./handlers/handlers.go")
	if err != nil {
		return err
	}

	handlersFileContent := `package handlers

import (
	"github.com/baraa-almasri/useless/songs"
	"io/ioutil"
	"net/http"
)

func getFullURL(shortURL string) string {
	url, _ := ioutil.ReadFile("./urls/" + shortURL)
	return string(url)
}

func handle_play_meme_song(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, songs.NewMemeSongs().GetRandomSong(), http.StatusFound)
}

func handle_rick_roll(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusFound)
}
`
	for _, shortURL := range urls {
		handlersFileContent += generateHandlerFunction(shortURL)
	}

	_, err = f.Write([]byte(handlersFileContent))
	if err != nil {
		return err
	}

	// happily ever after
	return nil
}

// generateHandlerFunction
func generateHandlerFunction(shortURL string) string {
	return fmt.Sprintf(`
func handle_%s(w http.ResponseWriter, r *http.Request) {
	url := getFullURL("%s")
	http.Redirect(w, r, url, http.StatusFound)
}`, shortURL, shortURL)
}

// GenerateMultiplexerFile generates urls multiplexer from a given map :)
func GenerateMultiplexerFile(urls []string) error {
	f, err := os.Create("./handlers/multiplexer.go")
	if err != nil {
		return err
	}

	muxFileCont := getMuxPrefix()

	for _, shortURL := range urls {
		muxFileCont += generateHandlerStatement(shortURL)
	}
	muxFileCont += "\n}"

	_, err = f.Write([]byte(muxFileCont))
	if err != nil {
		return err
	}

	// happily ever after
	return nil
}

// generateHandlerStatement returns a handling statement for a given short url
func generateHandlerStatement(shortURL string) string {
	return fmt.Sprintf("\n\tmux.HandleFunc(\"/%s/\", handle_%s)", shortURL, shortURL)
}

// getMuxPrefix returns multiplexer file content before the handling statements
func getMuxPrefix() string {
	return `package handlers

import (
	"net/http"
)

var (
	mux = http.NewServeMux()
)

func GetMux() *http.ServeMux {
	bunchHandle()
	return mux
}

func bunchHandle() {
	mux.HandleFunc("/play_meme_song/", handle_play_meme_song)
	mux.HandleFunc("/no_url/", handle_rick_roll)
`
}
