package linkshortener

import (
	"crypto/rand"
	"strings"
)

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"abcdefghijklmnopqrstuvwxyz" +
	"0123456789")

func GenerateShortURL() (string, error) {
	const size = 7
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	var result strings.Builder
	_, err = result.Write([]byte("http://localhost:8080/"))
	if err != nil {
		return "", err
	}
	for _, b := range b {
		if _, err := result.WriteRune(letters[int(b)%len(letters)]); err != nil {
			return "", err
		}
	}
	return result.String(), nil
}
