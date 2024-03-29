package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt returns a random int between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString returns a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner returns a random owner name (random string of length 6)
func RandomOwner() string {
	return RandomString(6)
}

// RandomAmount returns a random amount of money (between 1 and 1000)
func RandomAmount() int64 {
	return RandomInt(1, 1000)
}

// RandomCurrency returns a random currency
func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD, PLN}
	n := len(currencies)

	return currencies[rand.Intn(n)]
}

// RandomEmail returns a random email
func RandomEmail() string {
	return RandomString(6) + "@example.com"
}
