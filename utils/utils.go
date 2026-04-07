package utils

import (
	"fmt"
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

// BytesToRunes converts []byte to []rune directly (each byte becomes a rune with same value)
func BytesToRunes(b []byte) []rune {
	runes := make([]rune, len(b))
	for i, v := range b {
		runes[i] = rune(v)
	}
	return runes
}

// RunesToBytes converts []rune to []byte (rune values must be in range 0-255)
func RunesToBytes(r []rune) ([]byte, error) {
	bytes := make([]byte, len(r))
	for i, v := range r {
		if v < 0 || v > 255 {
			return nil, fmt.Errorf("rune value %d at position %d is out of byte range (0-255)", v, i)
		}
		bytes[i] = byte(v)
	}
	return bytes, nil
}
