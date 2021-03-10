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
