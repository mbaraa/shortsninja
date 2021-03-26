package config

import (
	"encoding/json"
	"os"
)

// Config represents the configuration of the server
type Config struct {
	Font                  string
	GoogleClientID        string
	GoogleClientSecret    string
	GoogleCallbackHandler string
	IPInfoIoToken         string
}

// LoadConfig loads the configuration from the file `./config.json`
// if the file doesn't exist the program crashes
func LoadConfig() *Config {
	confFile, err := os.Open("./config.json")
	if err != nil {
		panic(err)
	}

	conf := make(map[string]interface{})

	err = json.NewDecoder(confFile).Decode(&conf)
	if err != nil {
		panic(err)
	}

	return &Config{
		Font:                  conf["font_b64"].(string),
		GoogleClientID:        conf["google_client_id"].(string),
		GoogleClientSecret:    conf["google_client_secret"].(string),
		GoogleCallbackHandler: conf["google_callback_handler"].(string),
		IPInfoIoToken:         conf["ipinfo_io_token"].(string),
	}
}
