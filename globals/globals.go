package globals

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	ShortURLs     []string
	UsedShortURLs []string
)

// LoadGlobals initializes ShortURLs and UsedShortURLs using the given urls file and the content of `./urls` directory
// returns nil urls slices and an error when an error is occurred
// and returns ShortURLs, UsedShortURLs, nil respectively for later use
func LoadGlobals(urlsFile *os.File) ([]string, []string, error) {
	shorts, err := loadShortURLs(urlsFile)
	if err != nil {
		return nil, nil, err
	}

	used, err := loadUsedShortURLs()
	if err != nil {
		return nil, nil, err
	}

	// happily ever after
	return shorts, used, nil
}

// IsShortURLExist returns true if the given short URL exist in the short URLs list
func IsShortURLExist(shortURL string) bool {
	return elementExistInArr(shortURL, ShortURLs)
}

// IsShortURLUsed returns true if the given short URL is assigned to some URL
func IsShortURLUsed(shortURL string) bool {
	return elementExistInArr(shortURL, UsedShortURLs)
}

// elementExistInArr returns true if the given element exists in the given slice
func elementExistInArr(element string, slice []string) bool {
	for i := range slice {
		if slice[i] == element {
			return true
		}
	}

	return false
}

// loadShortURLs loads short URLs file's content into the global short URLs string slice
func loadShortURLs(urlsFile *os.File) ([]string, error) {
	rawShorts, err := ioutil.ReadAll(urlsFile)
	if err != nil {
		return nil, err
	}

	ShortURLs = strings.Split(string(rawShorts), "\n")

	// happily ever after
	return ShortURLs, nil
}

// loadUsedShortURLs loads used short URLs using the existing files in the "urls" directory
func loadUsedShortURLs() ([]string, error) {
	rawUsedURLs, err := exec.Command("ls", "./urls").Output()
	if err != nil {
		return nil, err
	}

	UsedShortURLs = strings.Split(string(rawUsedURLs), "\n")

	// happily ever after
	return UsedShortURLs, nil
}
