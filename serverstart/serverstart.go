package serverstart

import (
	"github.com/baraa-almasri/shortsninja/globals"
	"os"
)

// LoadURLs loads short URLs file's content into the global short URLs string slice
// also loads used short URLs using the existing files in the "urls" directory
func LoadURLs(urlsFile *os.File) error {
	_, _, err := globals.LoadGlobals(urlsFile)
	if err != nil {
		return err
	}

	// happily ever after
	return nil
}
