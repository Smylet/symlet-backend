package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + r.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[r.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

// RandomVerificationCode generates a numeric code of a given length for use as a verification code.
// RandomCode generates a random unsigned integer of a specified length (number of digits).
func RandomCode(length int) uint {
	// Calculate the maximum value based on the length: 10^length - 1
	max := 1
	for i := 0; i < length; i++ {
		max *= 10
	}

	// Generate a random number in the range [0, max - 1].
	// We're assuming that 'max' will not be greater than the maximum value a uint can hold.
	randomNumber := uint(r.Intn(max))

	return randomNumber
}
