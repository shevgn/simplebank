// Package util provides helper functions
package util

import (
	"math/rand"
	"strings"
)

func init() {
}

const letters = "abcdefghijklmnopqrstuvwxyz"

// RandomInt generates random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates random string of given length
func RandomString(length int) string {
	var sb strings.Builder
	k := len(letters)

	for range length {
		sb.WriteByte(letters[rand.Intn(k)])
	}

	return sb.String()
}

// RandomOwner generates random owner
func RandomOwner() string {
	return RandomString(8)
}

// RandomBalance generates random balance
func RandomBalance() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates random currency
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "UAH"}

	return currencies[rand.Intn(len(currencies))]
}
