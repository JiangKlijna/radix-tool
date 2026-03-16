package utils

import (
	"math/rand"
)

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

func ShuffleString(s string) string {
	runes := []rune(s)
	for i := len(runes) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
