package tools

import (
	"testing"
	"time"
)

func TestFloat64ToTime(t *testing.T) {
	tests := []struct {
		input    float64
		expected time.Time
	}{
		{0, time.Time{}}, // Zero value for f=0
		{1234567890, time.Unix(0, int64(1234567890*float64(time.Second)))},
		{1234567890.123, time.Unix(0, int64(1234567890.123*float64(time.Second)))},
	}

	for _, test := range tests {
		result := Float64ToTime(test.input)
		if !result.Equal(test.expected) {
			t.Errorf("Float64ToTime(%f) = %v, expected %v", test.input, result, test.expected)
		}
	}
}
