package linkshortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortURL(t *testing.T) {
	shortURLs := make(map[string]string)
	// we chek the uniqueness of each link
	for i := 0; i < 1000000; i++ {
		url, err := GenerateShortURL()
		assert.NoError(t, err, "got error when generating a short url")
		_, ok := shortURLs[url]
		assert.False(t, ok, "the short url is not unique", url, shortURLs[url])
		shortURLs[url] = url
	}
}
