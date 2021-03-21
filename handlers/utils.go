package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/baraa-almasri/shortsninja/globals"
	"github.com/baraa-almasri/useless/songs"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	memes = songs.NewMemeSongs()
)

// createAndUpdate creates a new url file that doesn't exist in the usedURLs slice and updates the usedURLs slice
// and returns the assigned short URL
func createAndUpdate(url, userName string) string {
	for _, v := range globals.ShortURLs {
		if !globals.IsShortURLUsed(v) {
			f, _ := os.Create("./urls/" + v)
			_, _ = f.Write([]byte(url + " " + userName))
			_ = f.Close()

			f, _ = os.Create("./tracking/" + userName + "/" + v + ".csv")
			_, _ = f.Write([]byte("short_url, full_url, visit_location, user_agent, visit_time"))
			_ = f.Close()

			globals.UsedShortURLs = append(globals.UsedShortURLs, v)
			return v
		}
	}

	return ""
}

// getURLData returns the full URL and username for the given short URL
func getURLData(shortURL string) (string, string) {
	urlFile, err := os.Open("./urls/" + shortURL)
	if err != nil || strings.Contains(shortURL, ".") { // wow, much security!
		return "/play_meme_song", "0"
	}

	var url, user string
	_, err = fmt.Fscanf(urlFile, "%s %s", &url, &user)
	if err != nil {
		return "/play_meme_song", "0"
	}

	// happily ever after
	return url, user
}

func appendData(shortURL, userName string, data map[string]interface{}) {
	f, err := os.OpenFile(fmt.Sprintf("./tracking/%s/%s.csv", userName, shortURL), os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}

	err = appendToFile(f, fmt.Sprintf("%s, %s, %s, %s, %d",
		shortURL, data["url"].(string), data["location"].(string),
		data["user_agent"].(string)[:strings.Index(data["user_agent"].(string), " ")],
		data["visit_time"].(int64)))
	if err != nil {
		panic(err)
	}
}

func getRequestData(req *http.Request) map[string]interface{} {
	data := make(map[string]interface{})
	ip := req.Header.Get("X-FORWARDED-FOR")
	data["location"] = getIPLocation(ip)
	data["user_agent"] = req.Header.Get("User-Agent")
	data["visit_time"] = time.Now().Unix()

	return data
}

func getIPLocation(ip string) string {
	resp, err := http.Get(fmt.Sprintf("https://ipinfo.io/%s?token=7254267b98ddeb", ip))
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

// Deprecated
func appendToFile(file *os.File, cont string) error {
	_, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(cont + "\n"))
	if err != nil {
		return err
	}
	file.Close()
	// happily ever after
	return nil
}
