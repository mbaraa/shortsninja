package initialsetup

import (
	"errors"
	"github.com/baraa-almasri/useless"
	"io"
	"math"
	"os"
)

// GenerateShortURLsFile generates a list of possible short URLs of the given criteria, and puts it in the given file
// when done returns the map with a nil error
// possible errors:
// from getEmptyShortURLs, io errors
// numChars is how much alphanumeric characters are in the short URL handler
// amount is how much short URLs handlers you want to generate,
// note that I tried to calculate it using numChars but it got stuck in a loop :(
// so when you provide an amount make sure it's at least 10^numChars lesser than 62^numChars
func GenerateShortURLsFile(numChars, amount int, urlsFile *os.File) error {
	urls, err := getEmptyShortURLs(numChars, amount)
	if err != nil {
		return err
	}

	for i := range urls {
		err = appendToFile(urlsFile, i)
		if err != nil {
			return err
		}
	}

	// happily ever after
	return nil
}

// getEmptyShortURLs generates a map of short URLs using the given criteria returns the generated map and no error when done
// current possible error is an impossible number of URLs
// the rest is the same as GenerateShortURLsFile
func getEmptyShortURLs(numChars, amount int) (map[string]string, error) {
	if amount > int(math.Pow(62, float64(numChars))) {
		return nil, errors.New("amount of urls can't be obtained")
	}

	randomizer := useless.NewRandASCII()
	urls := make(map[string]string)
	for i := 0; i < amount; i++ {
		short := randomizer.GetRandomAlphanumString(numChars)
		if urls[short] != "" {
			i--
			continue
		}
		urls[short] = "/play_meme_song"
	}

	// happily ever after
	return urls, nil
}

// appendToFile appends a string with a line feed to a given file
func appendToFile(file *os.File, cont string) error {
	_, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(cont + "\n"))
	if err != nil {
		return err
	}

	// happily ever after
	return nil
}
