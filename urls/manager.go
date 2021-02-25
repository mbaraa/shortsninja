package urls

import (
	"fmt"
	"github.com/baraa-almasri/shortsninja/db"
	"github.com/baraa-almasri/useless"
	"os"
	"time"
)

// URLManager manages URLs :)
type URLManager struct{
	db db.Database
}

// NewShortener returns a new URLManager instance
func NewShortener(db db.Database) *URLManager {
	return &URLManager{db}
}

// ShortenNewURL shortens a given URL, generates an appropriate handler function and a handling statement
// and adds it to the database if it doesn't exist, otherwise it generates an another one
// also userID = 0 means no_user
// the short URL consists of 4 chars
func (u *URLManager) ShortenNewURL(url string) string {
	r := useless.NewRandASCII()
	generateAndAdd:
	shortURL := r.GetRandomString(4)

	err := u.db.AddURL(0, shortURL, url, time.Now())
	if err != nil {
		goto generateAndAdd
	}

	// add the new url to the list
	urls, err := u.db.GetURLs()
	urls[shortURL] = url

	err = u.generateHandlerFunctions(urls)
	if err != nil {
		panic(err)
	}

	err = u.generateHandlerStatements(urls)
	if err != nil {
		panic(err)
	}

	return shortURL
}

// generateHandlerStatements generates a multiplexer handler file
func (u *URLManager) generateHandlerStatements(shortURLs map[string]string) error {
	f, err := os.Create("./handlers/multiplexer.go")
	if err != nil {
		return err
	}

	multiplexerFileContent := u.getMuxPrefix()
	for shortURL, _ := range shortURLs {
		multiplexerFileContent += u.generateHandlerStatement(shortURL)
	}
	multiplexerFileContent += "\n}"

	_, err = f.Write([]byte(multiplexerFileContent))
	if err != nil {
		return err
	}

	return nil
}

// generateHandlerFunctions generates a handlers' file
func (u *URLManager) generateHandlerFunctions(urls map[string]string) error {
	f, err := os.Create("./handlers/handlers.go")
	if err != nil {
		return err
	}

	handlersFileContent := "package handlers\n\nimport \"net/http\"\n\n"
	for shortURL, url := range urls {
		handlersFileContent += u.generateHandlerFunction(url, shortURL)
	}

	_, err = f.Write([]byte(handlersFileContent))
	if err != nil {
		return err
	}

	return nil
}

// getMuxPrefix returns multiplexer file content before the handling statements
func (u *URLManager) getMuxPrefix() string {
	return `package handlers

import "net/http"

var mux = &http.NewServeMux()

func GetMux() *http.ServeMux {
	bunchHandle()
	return mux
}

func bunchHandle() {
`
}

// generateHandlerStatement returns a handling statement for a given short url
func (u *URLManager) generateHandlerStatement(shortURL string) string {
	return fmt.Sprintf("\n\tmux.HandleFunc(\"/%s\", handle_%s)", shortURL, shortURL)
}

// getPagePrefix will return redirection page prefix
func (u *URLManager) getPagePrefix() string {
	return "<!DOCTYPE html><head><title>Redirecting...</title><script> window.open(\\\""
}

// getPagePostfix will return redirection page postfix
func (u *URLManager) getPagePostfix() string {
	return "\\\", \\\"_top\\\"); </script></head></html>"
}

// getURLPage returns redirection web page that will go to a given url
func (u *URLManager) getURLPage(url string) string {
	return fmt.Sprintf("%s%s%s", u.getPagePrefix(), url, u.getPagePostfix())
}

// generateHandlerFunction
func (u *URLManager) generateHandlerFunction(url, shortURL string) string {
	return fmt.Sprintf(`
func handle_%s(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("%s"))
}`, shortURL, u.getURLPage(url))
}

// BunchGenerate generates empty url fields
func (u *URLManager) BunchGenerate(amount int) {
	for i := 0; i < amount; i++ {
		u.ShortenNewURL("")
	}
}

// TODO
// just do it!
func (u *URLManager) RemoveURL(shortURL string) string {
	return ""
}