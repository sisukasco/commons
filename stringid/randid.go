package stringid

import (
	"math/rand"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const allchars = "abcdefghijklmnopqrstuvwxyz" + "0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandID(length int) string {
	return RandString(length)
}

func RandString(length int) string {
	if length <= 2 {
		return stringWithCharset(length, alphabet)
	}
	var id string
	for i := 0; i < 100; i++ {
		id = stringWithCharset(1, alphabet) +
			stringWithCharset(length-1, allchars)
		if !IsProfanity(id) {
			break
		}
	}
	return id
}
