package util

import (
	"crypto/rand"
	"encoding/base32"

	"github.com/go-faster/errors"
)

func GenerateRandomString(length int) (string, error) {
	randomBytes := make([]byte, length)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", errors.Wrap(err, "random string")
	}

	return base32.StdEncoding.EncodeToString(randomBytes)[:length], nil
}
