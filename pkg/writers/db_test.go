package writers

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/helviojunior/enumdns/pkg/models"
)

// TestDbWriterRewriteMutatedResult guards against a regression where a result
// object that was already written (and thus had its DB-managed ID populated via
// RETURNING) is mutated and written again. Because the logical identity is the
// `hash` and inserts use ON CONFLICT(hash), sending the stale ID must not cause a
// "UNIQUE constraint failed: results.id" error.
func TestDbWriterRewriteMutatedResult(t *testing.T) {
	dir := t.TempDir()
	uri := "sqlite://" + filepath.Join(dir, "results.db")

	w, err := NewDbWriter(uri, false)
	if err != nil {
		t.Fatalf("NewDbWriter: %v", err)
	}

	r := &models.Result{
		TestId: "t1", FQDN: "x.sec4us.com.br.", RType: "AAAA",
		IPv6: "2a01:111:f403:c801::1", Exists: true, ProbedAt: time.Now(),
	}
	if err := w.Write(r); err != nil {
		t.Fatalf("first write: %v", err)
	}

	// Mutate so the hash changes (e.g. enrichment performed after the first write),
	// keeping whatever ID the previous write populated.
	r.CloudProduct = "Microsoft Office 365"
	if err := w.Write(r); err != nil {
		t.Fatalf("rewrite of mutated result must not fail: %v", err)
	}
}

// TestDbWriterRewriteSameSOA ensures the shared/cached SOA object can be written
// repeatedly (once per host of a zone) without primary-key conflicts.
func TestDbWriterRewriteSameSOA(t *testing.T) {
	dir := t.TempDir()
	uri := "sqlite://" + filepath.Join(dir, "soa.db")

	w, err := NewDbWriter(uri, false)
	if err != nil {
		t.Fatalf("NewDbWriter: %v", err)
	}

	soa := &models.SOA{TestId: "t1", Name: "sec4us.com.br", ProbedAt: time.Now()}
	for i := 0; i < 3; i++ {
		if err := w.WriteSOA(soa); err != nil {
			t.Fatalf("WriteSOA #%d: %v", i, err)
		}
	}
}
