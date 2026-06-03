package tools

import (
	"testing"
)

func TestUnescapeDNSName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// No backslash -> returned unchanged (fast path).
		{"example.com", "example.com"},
		{"", ""},
		// miekg/dns decimal escapes for raw UTF-8 bytes (IDN / Latin chars).
		{`s\195\178doc.cf`, "sòdoc.cf"},
		{`s\195\182doc.co`, "södoc.co"},
		{`s\197\141doc.net.br`, "sōdoc.net.br"},
		{`mar\196\171sa.ml`, "marīsa.ml"},
		{`hostmaster.m\196\131risa.cf`, "hostmaster.mărisa.cf"},
		// Escaped special characters keep their literal meaning.
		{`a\.b.example.com`, "a.b.example.com"},
		{`a\\b.example.com`, `a\b.example.com`},
		// Trailing lone backslash / incomplete escape is preserved as-is byte.
		{`end\195`, "end\xc3"},
	}

	for _, tc := range tests {
		got := UnescapeDNSName(tc.input)
		if got != tc.expected {
			t.Errorf("UnescapeDNSName(%q) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

// A decoded name passed through again must not change (idempotent for real
// names, which contain no backslash after decoding).
func TestUnescapeDNSNameIdempotent(t *testing.T) {
	for _, s := range []string{"sòdoc.cf", "marīsa.ml", "example.com"} {
		if got := UnescapeDNSName(s); got != s {
			t.Errorf("UnescapeDNSName(%q) not idempotent: got %q", s, got)
		}
	}
}
