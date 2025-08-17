package tools

import (
	"testing"
)

func TestFormatInt(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{123, "123"},
		{1234, "1.234"},
		{1234567, "1.234.567"},
		{-1234, "-1.234"},
	}

	for _, test := range tests {
		result := FormatInt(test.input)
		if result != test.expected {
			t.Errorf("FormatInt(%d) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestSliceHasStr(t *testing.T) {
	slice := []string{"apple", "banana", "orange"}

	if !SliceHasStr(slice, "banana") {
		t.Error("SliceHasStr should return true for existing item")
	}

	if SliceHasStr(slice, "grape") {
		t.Error("SliceHasStr should return false for non-existing item")
	}

	// Test empty slice
	emptySlice := []string{}
	if SliceHasStr(emptySlice, "apple") {
		t.Error("SliceHasStr should return false for empty slice")
	}
}

func TestRandSleep(t *testing.T) {
	// This function doesn't return anything and just sleeps
	// We'll just test that it doesn't panic
	RandSleep()
	// If we reach here, the function didn't panic
}

func TestSliceHasInt(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	if !SliceHasInt(slice, 3) {
		t.Error("SliceHasInt should return true for existing item")
	}

	if SliceHasInt(slice, 6) {
		t.Error("SliceHasInt should return false for non-existing item")
	}

	// Test empty slice
	emptySlice := []int{}
	if SliceHasInt(emptySlice, 1) {
		t.Error("SliceHasInt should return false for empty slice")
	}
}

func TestSliceHasUInt16(t *testing.T) {
	slice := []uint16{80, 443, 8080, 8443}

	if !SliceHasUInt16(slice, 443) {
		t.Error("SliceHasUInt16 should return true for existing item")
	}

	if SliceHasUInt16(slice, 9999) {
		t.Error("SliceHasUInt16 should return false for non-existing item")
	}

	// Test empty slice
	emptySlice := []uint16{}
	if SliceHasUInt16(emptySlice, 80) {
		t.Error("SliceHasUInt16 should return false for empty slice")
	}
}

func TestUniqueIntSlice(t *testing.T) {
	tests := []struct {
		input    []int
		expected []int
	}{
		{[]int{1, 2, 3, 2, 1}, []int{1, 2, 3}},
		{[]int{}, []int{}},
		{[]int{1}, []int{1}},
		{[]int{1, 1, 1}, []int{1}},
		{[]int{1, 2, 3}, []int{1, 2, 3}},
	}

	for _, test := range tests {
		result := UniqueIntSlice(test.input)
		if len(result) != len(test.expected) {
			t.Errorf("UniqueIntSlice(%v) length = %d, expected %d", test.input, len(result), len(test.expected))
			continue
		}

		for _, expected := range test.expected {
			if !SliceHasInt(result, expected) {
				t.Errorf("UniqueIntSlice(%v) missing expected value %d", test.input, expected)
			}
		}
	}
}

func TestShuffleStr(t *testing.T) {
	slice := []string{"a", "b", "c", "d", "e"}
	original := make([]string, len(slice))
	copy(original, slice)

	ShuffleStr(slice)

	// Check that all original elements are still present
	for _, orig := range original {
		if !SliceHasStr(slice, orig) {
			t.Errorf("ShuffleStr modified slice contents, missing %s", orig)
		}
	}

	// Check that slice length is preserved
	if len(slice) != len(original) {
		t.Errorf("ShuffleStr changed slice length from %d to %d", len(original), len(slice))
	}
}
