package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/baraa-almasri/shortsninja/globals"
	"github.com/baraa-almasri/shortsninja/models"
	"github.com/baraa-almasri/useless"
	"github.com/baraa-almasri/useless/songs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"os"
	"strings"
)

var (
	memes             = songs.NewMemeSongs()
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("CALL_BACK_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	randomizer = useless.NewRandASCII()
	state      = ""
)

// createAndUpdate creates a new short url that doesn't exist in the db,
// adds the new short URL to the database and returns the assigned short URL
func createAndUpdate(url string, user *models.User) string {
	// storing the generated short url so it can be returned :)
	var newURL *models.URL

	// loop until the generated short url doesn't exist in the db
	for {
		newURL = &models.URL{
			Short:     randomizer.GetRandomAlphanumString(5),
			FullURL:   url,
			UserEmail: user.Email,
		}
		if globals.DBManager.AddURL(newURL) == nil {
			break
		}
	}

	return newURL.Short
}

// getURLData returns the full URL for the given short URL
func getURLData(shortURL string) string {
	url, err := globals.DBManager.GetURL(shortURL)
	if err != nil || strings.Contains(shortURL, ".") { // wow, much security!
		return "/play_meme_song"
	}

	// happily ever after
	return url
}

// getRequestData returns a map with the needed request headers
func getRequestData(req *http.Request) map[string]interface{} {
	data := make(map[string]interface{})
	ip := req.Header.Get("X-FORWARDED-FOR")
	data["location"] = getIPLocation(ip)
	data["user_agent"] = req.Header.Get("User-Agent")

	return data
}

// getIPLocation return a string of the IP's location using ipinfo.io
func getIPLocation(ip string) string {
	resp, err := http.Get(fmt.Sprintf("https://ipinfo.io/%s?token=%s", ip, globals.IPInfoToken))
	if err != nil {
		return "NULL/NULL"
	}

	defer resp.Body.Close()

	ipData := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&ipData)
	if err != nil {
		return "NULL/NULL"
	}

	return fmt.Sprintf("%s/%s", ipData["region"], ipData["country"])
}
