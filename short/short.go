package short

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
)

var Shorter shorter

type shorter struct {
	ltos map[string]string
	stol map[string]string
}

// initSequence will panic when it can not open the sequence successfully.
func (shorter *shorter) mustInit() {
	shorter.ltos = make(map[string]string)
	shorter.stol = make(map[string]string)
}

func (shorter *shorter) Short(longURL string) (shortURL string, err error) {
	hash := generateHash(longURL)
	shortURL = hash[:7] // Take the first 7 characters of the hash
	shorter.ltos[longURL] = shortURL
	shorter.stol[shortURL] = longURL

	return shortURL, nil
}

func generateHash(url string) string {
	hasher := md5.New()
	hasher.Write([]byte(url))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (shorter *shorter) Expand(shortURL string) (longURL string, err error) {
	longURL, ok := shorter.stol[shortURL]
	if !ok {
		return "", errors.New("short URL not found")
	}
	return longURL, nil
}

func (s *shorter) GetShortenedURLs() []string {
	longURLs := make([]string, 0, len(s.ltos))
	for longURL := range s.ltos {
		longURLs = append(longURLs, longURL)
	}
	return longURLs
}

func Start() {
	Shorter.mustInit()
	log.Println("shorter starts")
}
