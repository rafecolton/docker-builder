package vauth

import (
	"github.com/martini-contrib/auth"
)

// SecureCompare performs a constant time compare of two strings to limit timing attacks.
func SecureCompare(given string, actual string) bool {
	return auth.SecureCompare(given, actual)
}
