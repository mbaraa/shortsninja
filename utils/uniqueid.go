package utils

import (
	"github.com/baraa-almasri/useless"
	"time"
)

// UniqueID is a wrapping for a somehow unique id generator
//
type UniqueID struct {
	randomizer *useless.RandASCII
}

// NewUniqueID returns a new UniqueID instance
//
func NewUniqueID(charRandomizer *useless.RandASCII) *UniqueID {
	return &UniqueID{charRandomizer}
}

// GetUniqueString returns a unique string generated from the current timestamp mixed with some random characters
// and adds an additional random chars to it if needed :)
//
func (uid *UniqueID) GetUniqueString(additionalLength ...int) string {
	return uid.getUniqueString(additionalLength...)
}

// GetUniqueStringWithPrefix same as GetUniqueString but adds a prefix to the generated string
//
func (uid *UniqueID) GetUniqueStringWithPrefix(prefix string, additionalLength ...int) string {
	return prefix + uid.getUniqueString(additionalLength...)
}

// getUniqueString returns a unique string generated from the current timestamp mixed with some random characters
// and adds an additional random chars to it if needed :)
func (uid *UniqueID) getUniqueString(additionalLength ...int) string {
	if additionalLength == nil {
		additionalLength = make([]int, 1)
	}
	return uid.convertTimestampToString(time.Now().Unix()) +
		uid.randomizer.GetRandomAlphanumString(additionalLength[0])
}

// getAlphanumCharFromInt returns an alphanumeric character from a given integer,
// when the integer actually represents an alphanumeric character it returns it as is,
// otherwise it shifts it to a known alphanumeric characters range,
// and when it can't be shifted it returns a random alphanumeric character
//
func (uid *UniqueID) getAlphanumCharFromInt(n int) uint8 {
	if n >= '0' && n <= '9' ||
		n >= 'A' && n <= 'Z' ||
		n >= 'a' && n <= 'z' {

		return uint8(n)
	}

	other := uint8(n%'z' + '0')

	if other < '0' {
		return other + ('0' - other)
	} else if other > '9' && other < 'A' {
		return other - (other - '9')
	}

	if other < 'A' && other > '9' {
		return other + ('A' - other)
	} else if other > 'Z' && other < 'a' {
		return other - (other - 'Z')
	}

	if other < 'a' && other > 'Z' {
		return other + ('a' - other)
	} else if other > 'z' {
		return other - (other - 'z')
	}
	return uid.randomizer.GetRandomAlphanumChar()
}

// convertTimestampToString partitions a given timestamp to a string
// where each character consists of 3 digits from the timestamp
//
func (uid *UniqueID) convertTimestampToString(ts int64) string {
	str := ""
	for ; ts > 0; ts /= 100 {
		str += string(uid.getAlphanumCharFromInt(
			int(ts%100) + 23), // +23 when the value=99 so add 23 to reach 'z'
		)
	}
	return str
}
