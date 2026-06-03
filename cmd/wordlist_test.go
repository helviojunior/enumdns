package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// writeWordlistInput writes content to a temp file and returns its path.
func writeWordlistInput(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "in.txt")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writing input: %s", err)
	}
	return path
}

func collect(t *testing.T, content string, keepTLD bool) map[string]struct{} {
	t.Helper()
	wordlistOpts.MinLength = 1
	wordlistOpts.MaxLength = 63
	wordlistOpts.KeepTLD = keepTLD
	wordlistOpts.KeepNumeric = false

	out := map[string]struct{}{}
	if _, err := collectTokens(writeWordlistInput(t, content), map[string]struct{}{}, out); err != nil {
		t.Fatalf("collectTokens: %s", err)
	}
	return out
}

func assertHas(t *testing.T, out map[string]struct{}, words ...string) {
	t.Helper()
	for _, w := range words {
		if _, ok := out[w]; !ok {
			t.Errorf("expected label %q to be kept, but it was dropped", w)
		}
	}
}

func assertMissing(t *testing.T, out map[string]struct{}, words ...string) {
	t.Helper()
	for _, w := range words {
		if _, ok := out[w]; ok {
			t.Errorf("expected token %q to be dropped, but it was kept", w)
		}
	}
}

// Only the trailing TLD of a dotted name must be stripped. An internal label
// that merely shares a name with a gTLD (e.g. "dev" in ed.dev.br) is kept.
func TestCollectTokensStripsOnlyTrailingTLD(t *testing.T) {
	out := collect(t, "ed.dev.br\nwww.example.com\nfoo.dev\n", false)
	assertHas(t, out, "ed", "dev", "www", "example", "foo")
	assertMissing(t, out, "br", "com")
}

// A bare token (no dot) is never treated as a TLD, even when it is one.
func TestCollectTokensKeepsBareTLDToken(t *testing.T) {
	out := collect(t, "dev\ncloud\napp\n", false)
	assertHas(t, out, "dev", "cloud", "app")
}

// --keep-tld retains the trailing TLD as well.
func TestCollectTokensKeepTLD(t *testing.T) {
	out := collect(t, "ed.dev.br\nwww.example.com\n", true)
	assertHas(t, out, "ed", "dev", "br", "www", "example", "com")
}

// All-numeric tokens (IP octets, TTLs, SOA serials) are dropped by default but
// kept with --keep-numeric. Mixed alphanumeric labels are never affected.
func TestCollectTokensNumeric(t *testing.T) {
	const in = "144.22.130.226\nttl 3600\nns1-01\nweb01\n"

	dropped := collect(t, in, false)
	assertHas(t, dropped, "ttl", "ns1-01", "web01")
	assertMissing(t, dropped, "144", "22", "130", "226", "3600")

	wordlistOpts.KeepNumeric = true
	kept := map[string]struct{}{}
	if _, err := collectTokens(writeWordlistInput(t, in), map[string]struct{}{}, kept); err != nil {
		t.Fatalf("collectTokens: %s", err)
	}
	assertHas(t, kept, "144", "226", "3600", "ns1-01")
}
