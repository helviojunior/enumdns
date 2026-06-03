package models

import (
	"testing"
)

// An IDN whose labels arrive as miekg/dns "\DDD" escapes must be decoded to raw
// UTF-8 in every serialized/hashed form, so it is stored as "sòdoc.cf" and never
// as "s\195\178doc.cf".
func TestResultNameUnescaped(t *testing.T) {
	r := &Result{
		FQDN:   `s\195\178doc.cf`,
		RType:  "CNAME",
		Target: `cn\195\178.example.com`,
	}

	if got, want := r.String(), "sòdoc.cf: cnò.example.com"; got != want {
		t.Errorf("String() = %q; want %q", got, want)
	}

	b, err := r.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON error: %v", err)
	}
	js := string(b)
	for _, frag := range []string{`"fqdn":"sòdoc.cf"`, `"target":"cnò.example.com"`} {
		if !contains(js, frag) {
			t.Errorf("MarshalJSON output %s missing %q", js, frag)
		}
	}
}

// SOA name hashing (used for the obs_dns _id) must be over the decoded name, so
// re-scanning a homoglyph zone yields the same id as the corrected records.
func TestSOAHashUnescaped(t *testing.T) {
	escaped := SOA{Name: `m\196\131risa.cf`}
	decoded := SOA{Name: "mărisa.cf"}
	if escaped.GetHash() != decoded.GetHash() {
		t.Errorf("escaped SOA hash %s != decoded SOA hash %s",
			escaped.GetHash(), decoded.GetHash())
	}
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
