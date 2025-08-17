package ascii

import (
	"strings"
	"testing"
)

func TestLogo(t *testing.T) {
	logo := Logo()

	if logo == "" {
		t.Error("Logo should not be empty")
	}

	// Check if logo contains expected parts
	expectedParts := []string{"Ver:"}
	for _, part := range expectedParts {
		if !strings.Contains(logo, part) {
			t.Errorf("Logo should contain '%s'", part)
		}
	}
}

func TestLogoHelp(t *testing.T) {
	helpText := "This is help text"
	logoWithHelp := LogoHelp(helpText)

	if logoWithHelp == "" {
		t.Error("LogoHelp should not be empty")
	}

	// Should contain both logo and help text
	if !strings.Contains(logoWithHelp, helpText) {
		t.Error("LogoHelp should contain help text")
	}

	// Should contain logo parts
	if !strings.Contains(logoWithHelp, "Ver:") {
		t.Error("LogoHelp should contain logo")
	}
}

func TestScapeAnsi(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"hello", "hello"},
		{"\033[32mgreen\033[0m", "green"},
		{"\033[31mred\033[32mgreen\033[0m", "redgreen"},
		{"normal text", "normal text"},
		{"\033[1;31mbold red\033[0m", "bold red"},
	}

	for _, test := range tests {
		result := ScapeAnsi(test.input)
		if result != test.expected {
			t.Errorf("ScapeAnsi(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestGetNextSpinner(t *testing.T) {
	// Test that GetNextSpinner returns different values
	spinner1 := GetNextSpinner("")
	spinner2 := GetNextSpinner(spinner1)

	if spinner1 == "" {
		t.Error("GetNextSpinner should not return empty string")
	}

	if spinner2 == "" {
		t.Error("GetNextSpinner should not return empty string")
	}

	// After several calls, we should eventually get different values
	spinners := make(map[string]bool)
	current := ""
	for i := 0; i < 10; i++ {
		s := GetNextSpinner(current)
		spinners[s] = true
		current = s
	}

	if len(spinners) < 2 {
		t.Error("GetNextSpinner should cycle through different spinner characters")
	}
}

func TestColoredSpin(t *testing.T) {
	result := ColoredSpin("test")

	if result == "" {
		t.Error("ColoredSpin should not return empty string")
	}

	if !strings.Contains(result, "test") {
		t.Error("ColoredSpin should contain input text")
	}
}

func TestMarkdown(t *testing.T) {
	input := "This is **bold** and *italic* text"
	result := Markdown(input)

	if result == "" {
		t.Error("Markdown should not return empty string")
	}

	// The result should be different from input (processed in some way)
	if result == input {
		t.Error("Markdown should process the input")
	}
}

func TestClearLine(t *testing.T) {
	// Test that ClearLine doesn't panic
	ClearLine()
	// If we reach here, the function didn't panic
}

func TestClear(t *testing.T) {
	// Test that Clear doesn't panic
	Clear()
	// If we reach here, the function didn't panic
}

func TestShowCursor(t *testing.T) {
	// Test that ShowCursor doesn't panic
	ShowCursor()
	// If we reach here, the function didn't panic
}

func TestHideCursor(t *testing.T) {
	// Test that HideCursor doesn't panic
	HideCursor()
	// If we reach here, the function didn't panic
}

func TestSetConsoleColors(t *testing.T) {
	// Test that SetConsoleColors doesn't panic
	SetConsoleColors()
	// If we reach here, the function didn't panic
}
