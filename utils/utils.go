package utils

import (
	"math/rand"
	"os"
)

// RemoveDuplicates removes duplicate characters from a string, preserving order
func RemoveDuplicates(s string) string {
	seen := make(map[rune]bool)
	result := []rune{}
	for _, r := range s {
		if !seen[r] {
			seen[r] = true
			result = append(result, r)
		}
	}
	return string(result)
}

// HasDuplicateChars checks whether a string contains duplicate characters
func HasDuplicateChars(s string) bool {
	seen := make(map[rune]bool)
	for _, r := range s {
		if seen[r] {
			return true
		}
		seen[r] = true
	}
	return false
}

// ShuffleString randomly shuffles the characters in a string
func ShuffleString(s string) string {
	runes := []rune(s)
	for i := len(runes) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ReadFile reads the contents of the file named by filename
func ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// WriteFile writes data to a file named by filename with the given permissions
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}
