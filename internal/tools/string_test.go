package tools

import (
	"testing"
)

func TestLeftTrucate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"hello world", 10, "d"},     // Remove first 10 chars, leave "d"
		{"hello", 10, "hello"},       // String shorter than max, return as is
		{"hello world", 5, " world"}, // Remove first 5 chars, leave " world"
		{"", 5, ""},                  // Empty string
		{"a", 0, "a"},                // maxLen 0, return original
		{"hello", 0, "hello"},        // maxLen 0, return original
	}

	for _, test := range tests {
		result := LeftTrucate(test.input, test.maxLen)
		if result != test.expected {
			t.Errorf("LeftTrucate(%q, %d) = %q, expected %q", test.input, test.maxLen, result, test.expected)
		}
	}
}

func TestFormatInt64(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "0"},
		{123, "123"},
		{1234, "1.234"},
		{1234567, "1.234.567"},
		{-1234, "-1.234"},
		{1000000000, "1.000.000.000"},
	}

	for _, test := range tests {
		result := FormatInt64(test.input)
		if result != test.expected {
			t.Errorf("FormatInt64(%d) = %q, expected %q", test.input, result, test.expected)
		}
	}
}
