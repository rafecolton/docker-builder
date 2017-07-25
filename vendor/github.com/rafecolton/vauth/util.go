package vauth

import (
	"crypto/subtle"
)

// SecureCompare performs a constant time compare of two strings to limit timing attacks.
func SecureCompare(x string, y string) bool {
	if len(x) != len(y) {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(x), []byte(y)) == 1
}
