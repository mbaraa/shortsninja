package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/baraa-almasri/shortsninja/globals"
	"github.com/baraa-almasri/shortsninja/models"
	"github.com/baraa-almasri/useless"
	"net/http"
	"regexp"
	"strings"
)

var (
	randomizer = useless.NewRandASCII()
)

// renderPageFromSessionToken generates the required web page with the given user's data
func renderPageFromSessionToken(pageName, token, ip string, w http.ResponseWriter, r *http.Request) {
	user := getUser(token, ip)

	if r.URL.Query()["token"] != nil {
		token = r.URL.Query()["token"][0]
		user = getUser(token, ip)

		if pageName == "shorten" {
			http.SetCookie(w, &http.Cookie{
				Name:  "token",
				Value: token,
			})
		}
	}
	tempData := map[string]string{
		"Avatar":  user.Avatar,
		"Email":   user.Email,
		"FontB64": globals.Config.Font,
	}
	_ = globals.Templates.ExecuteTemplate(w, pageName, tempData)
}

func getIPAndToken(w http.ResponseWriter, r *http.Request) (ip, token string) {
	token = ""
	if r.Header.Get("Cookie") != "" {
		token = r.Header.Get("Cookie")[len("token="):]
		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: token,
		})
	}
	ip = r.Header.Get("X-FORWARDED-FOR")

	return
}

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

// getFullURL returns the full URL for the given short URL
func getFullURL(shortURL string) string {
	url, err := globals.DBManager.GetURL(shortURL)
	if err != nil || strings.Contains(shortURL, ".") { // wow, much security!
		return "/no_url/" // get rick rolled :)
	}

	// happily ever after
	return url
}

// getRequestData returns a map with the needed request headers
func getRequestData(req *http.Request) *models.URLData {
	return &models.URLData{
		IP:            req.Header.Get("X-FORWARDED-FOR"),
		VisitLocation: getIPLocation(req.Header.Get("X-FORWARDED-FOR")),
		UserAgent:     req.Header.Get("User-Agent"),
	}
}

// getIPLocation return a string of the IP's location using ipinfo.io
func getIPLocation(ip string) string {
	resp, err := http.Get(fmt.Sprintf("https://ipinfo.io/%s?token=%s", ip, globals.Config.IPInfoIoToken))
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

// isURLValid returns true when the given URL is valid
func isURLValid(url string) bool {
	validURLPatt := regexp.MustCompile(
		`[a-z]{0,255}[.]?[a-z]{1,255}[.][a-z]{1,255}[a-z;A-Z;0-9;&;=;#;/;?;-]{1,1000}`)

	return (strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) &&
		(validURLPatt.MatchString(url))
}

// isShortURLValid returns true when the short url is valid
func isShortURLValid(short string) bool {
	shortURLPattern, _ := regexp.Compile("[A-Z;0-9;a-z]{4,5}")
	return short == shortURLPattern.FindString(short)
}

// getUser returns a user using the given session token, if no user exists
// or the token doesn't match its previous IP it returns a dummy user with shorts ninja icon
func getUser(token, ip string) *models.User {
	realSession, err := globals.DBManager.GetSession(token)
	if err != nil {
		return getDummyUser()
	}

	user, _ := globals.DBManager.GetUser(&models.User{Email: realSession.UserEmail})

	if user == nil || realSession == nil || realSession.IP != ip {
		return getDummyUser()
	}
	return user
}

// getDummyUser returns an unsigned-in user!
func getDummyUser() *models.User {
	return &models.User{
		Email:        "",
		Avatar:       "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",
		CreationDate: nil,
	}
}
