package controllers

import (
	"regexp"
	"strings"
)

// URLValidator is used to check URL or short URL validity
type URLValidator struct{}

// NewURLValidator returns a new URLValidator instance
func NewURLValidator() *URLValidator {
	return new(URLValidator)
}

// IsURLValid returns true when the given URL is valid
func (*URLValidator) IsURLValid(url string) bool {
	validURLPatt := regexp.MustCompile(
		`[a-z0-9]{0,255}[.]?[a-z0-9]{1,255}[.][a-z0-9]{1,255}[:]?[0-9]{0,5}[a-zA-Z0-9/$\-_.+!*‘(),#?=:%]{1,1000}`)

	return (strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") ||
		strings.HasPrefix(url, "ftp://")) && (validURLPatt.MatchString(url))
}

// IsShortURLValid returns true when the short url is valid
func (*URLValidator) IsShortURLValid(short string) bool {
	shortURLPattern, _ := regexp.Compile(`[A-Z0-9a-z\-]{4,5}`)
	return short == shortURLPattern.FindString(short)
}
