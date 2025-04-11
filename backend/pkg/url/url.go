package url

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"regexp"
)

var regexURL *regexp.Regexp

func init() {
	regexURL = regexp.MustCompile(`^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&//=]*)$`)
}

func Validate(url string) bool {
	return regexURL.Match([]byte(url))
}

func GetHashURL(url string) string {
	h := sha256.New()
	h.Write([]byte(url))
	b := h.Sum(nil)
	return string(b)
}

func GenShortURL(url string) string {
	b := make([]byte, 4) // 4 байта => 6 base64 символов
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)[:6]
}
