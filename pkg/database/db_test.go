package database

import (
	"testing"
)

func TestConnectionErrorHandling(t *testing.T) {
	// Test with invalid URI
	_, err := Connection("invalid://uri", false, false)
	if err == nil {
		t.Error("Expected error for invalid URI")
	}

	// Test with unsupported scheme
	_, err = Connection("unsupported://test", false, false)
	if err == nil {
		t.Error("Expected error for unsupported scheme")
	}

	// Test with nonexistent sqlite file when shouldExist is true
	_, err = Connection("sqlite://nonexistent.db", true, false)
	if err == nil {
		t.Error("Expected error for nonexistent sqlite file")
	}
}
