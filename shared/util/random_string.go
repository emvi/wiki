package util

import (
	"math/rand"
	"time"
)

const (
	randStrChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// GenRandomString generates a random string with length n for a fixed alphabet.
func GenRandomString(n uint) string {
	rand.Seed(time.Now().UnixNano())
	randStr := make([]byte, n)

	for i := range randStr {
		randStr[i] = randStrChars[rand.Intn(len(randStrChars))]
	}

	return string(randStr)
}
