package tools

import (
	"crypto/rand"
	"math/big"
)

// SliceHasStr checks if a slice has a string
func SliceHasStr(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}

// SliceHasInt checks if a slice has an int
func SliceHasInt(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}

// SliceHasInt checks if a slice has an int
func SliceHasUInt16(slice []uint16, item uint16) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}

// UniqueIntSlice returns a slice of unique ints
func UniqueIntSlice(slice []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, num := range slice {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}

	return result
}

// ShuffleStr shuffles a slice of strings
func ShuffleStr(slice []string) {
	// Fisher-Yates shuffle algorithm using crypto/rand
	for i := len(slice) - 1; i > 0; i-- {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		j := int(n.Int64())
		slice[i], slice[j] = slice[j], slice[i]
	}
}
