package log

import (
	"testing"
)

func TestDebugf(t *testing.T) {
	// Test that it doesn't panic
	Debugf("test debug message: %s", "value")
}

func TestInfof(t *testing.T) {
	// Test that it doesn't panic
	Infof("test info message: %s", "value")
}

func TestWarnf(t *testing.T) {
	// Test that it doesn't panic
	Warnf("test warn message: %s", "value")
}

func TestErrorf(t *testing.T) {
	// Test that it doesn't panic
	Errorf("test error message: %s", "value")
}

func TestFatalf(t *testing.T) {
	// Skip fatal test as it would exit the program
	t.Skip("Skipping Fatalf test as it would exit the program")
}

func TestDebug(t *testing.T) {
	// Test that it doesn't panic
	Debug("test debug message")
}

func TestInfo(t *testing.T) {
	// Test that it doesn't panic
	Info("test info message")
}

func TestWarn(t *testing.T) {
	// Test that it doesn't panic
	Warn("test warn message")
}

func TestError(t *testing.T) {
	// Test that it doesn't panic
	Error("test error message")
}

func TestWith(t *testing.T) {
	logger := With("key", "value")
	if logger == nil {
		t.Error("With should return a logger")
	}
}
