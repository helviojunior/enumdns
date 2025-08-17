package advanced

import (
	"strings"
	"testing"
)

// Test TyposquattingTechnique
func TestTyposquattingTechnique(t *testing.T) {
	technique := &TyposquattingTechnique{}

	if technique.Name() != "typosquatting" {
		t.Errorf("Expected name 'typosquatting', got %s", technique.Name())
	}

	if technique.GetRiskLevel() != "high" {
		t.Errorf("Expected risk level 'high', got %s", technique.GetRiskLevel())
	}

	if technique.GetConfidence() != 0.8 {
		t.Errorf("Expected confidence 0.8, got %f", technique.GetConfidence())
	}

	variations := technique.Generate("test.com", []string{"com", "net"})

	if len(variations) == 0 {
		t.Error("Expected typosquatting variations to be generated")
	}

	// Check that all variations use the correct technique
	for _, v := range variations {
		if v.Technique != "typosquatting" {
			t.Errorf("Expected technique 'typosquatting', got %s", v.Technique)
		}
		if v.BaseDomain != "test.com" {
			t.Errorf("Expected BaseDomain 'test.com', got %s", v.BaseDomain)
		}
		if v.Confidence != 0.8 {
			t.Errorf("Expected confidence 0.8, got %f", v.Confidence)
		}
		if v.Risk != "high" {
			t.Errorf("Expected risk 'high', got %s", v.Risk)
		}
	}

	// Check that we have variations for both TLDs
	comCount := 0
	netCount := 0
	for _, v := range variations {
		if strings.HasSuffix(v.Domain, ".com") {
			comCount++
		}
		if strings.HasSuffix(v.Domain, ".net") {
			netCount++
		}
	}

	if comCount == 0 {
		t.Error("Expected .com variations")
	}
	if netCount == 0 {
		t.Error("Expected .net variations")
	}
}

// Test BitsquattingTechnique
func TestBitsquattingTechnique(t *testing.T) {
	technique := &BitsquattingTechnique{}

	if technique.Name() != "bitsquatting" {
		t.Errorf("Expected name 'bitsquatting', got %s", technique.Name())
	}

	if technique.GetRiskLevel() != "medium" {
		t.Errorf("Expected risk level 'medium', got %s", technique.GetRiskLevel())
	}

	if technique.GetConfidence() != 0.6 {
		t.Errorf("Expected confidence 0.6, got %f", technique.GetConfidence())
	}

	variations := technique.Generate("test.com", []string{"com"})

	if len(variations) == 0 {
		t.Error("Expected bitsquatting variations to be generated")
	}

	// Check properties
	for _, v := range variations {
		if v.Technique != "bitsquatting" {
			t.Errorf("Expected technique 'bitsquatting', got %s", v.Technique)
		}
		if v.Confidence != 0.6 {
			t.Errorf("Expected confidence 0.6, got %f", v.Confidence)
		}
		if v.Risk != "medium" {
			t.Errorf("Expected risk 'medium', got %s", v.Risk)
		}
	}
}

// Test HomographicTechnique
func TestHomographicTechnique(t *testing.T) {
	technique := &HomographicTechnique{}

	if technique.Name() != "homographic" {
		t.Errorf("Expected name 'homographic', got %s", technique.Name())
	}

	if technique.GetRiskLevel() != "high" {
		t.Errorf("Expected risk level 'high', got %s", technique.GetRiskLevel())
	}

	if technique.GetConfidence() != 0.9 {
		t.Errorf("Expected confidence 0.9, got %f", technique.GetConfidence())
	}

	variations := technique.Generate("test.com", []string{"com"})

	if len(variations) == 0 {
		t.Error("Expected homographic variations to be generated")
	}

	// Check that we have variations for characters with homographs
	// The domain "test" should have homographs for 'e' and 't'
	for _, v := range variations {
		if v.Technique != "homographic" {
			t.Errorf("Expected technique 'homographic', got %s", v.Technique)
		}
	}
}

// Test InsertionTechnique
func TestInsertionTechnique(t *testing.T) {
	technique := &InsertionTechnique{}

	if technique.Name() != "insertion" {
		t.Errorf("Expected name 'insertion', got %s", technique.Name())
	}

	variations := technique.Generate("ab.com", []string{"com"})

	if len(variations) == 0 {
		t.Error("Expected insertion variations to be generated")
	}

	// Check that variations are longer than original
	originalBase := "ab"
	for _, v := range variations {
		if v.Technique != "insertion" {
			t.Errorf("Expected technique 'insertion', got %s", v.Technique)
		}

		baseName := getBaseName(v.Domain)
		if len(baseName) <= len(originalBase) {
			t.Errorf("Expected inserted variation to be longer than original. Got %s", baseName)
		}
	}
}

// Test DeletionTechnique
func TestDeletionTechnique(t *testing.T) {
	technique := &DeletionTechnique{}

	if technique.Name() != "deletion" {
		t.Errorf("Expected name 'deletion', got %s", technique.Name())
	}

	variations := technique.Generate("test.com", []string{"com"})

	if len(variations) == 0 {
		t.Error("Expected deletion variations to be generated")
	}

	// Check that variations are shorter than original
	originalBase := "test"
	for _, v := range variations {
		if v.Technique != "deletion" {
			t.Errorf("Expected technique 'deletion', got %s", v.Technique)
		}

		baseName := getBaseName(v.Domain)
		if len(baseName) >= len(originalBase) {
			t.Errorf("Expected deleted variation to be shorter than original. Got %s", baseName)
		}
	}
}

// Test TranspositionTechnique
func TestTranspositionTechnique(t *testing.T) {
	technique := &TranspositionTechnique{}

	if technique.Name() != "transposition" {
		t.Errorf("Expected name 'transposition', got %s", technique.Name())
	}

	variations := technique.Generate("test.com", []string{"com"})

	if len(variations) == 0 {
		t.Error("Expected transposition variations to be generated")
	}

	// Check that we have expected transpositions
	expectedVariations := []string{"etst.com", "tset.com", "tets.com"}
	foundCount := 0

	for _, v := range variations {
		if v.Technique != "transposition" {
			t.Errorf("Expected technique 'transposition', got %s", v.Technique)
		}

		for _, expected := range expectedVariations {
			if v.Domain == expected {
				foundCount++
				break
			}
		}
	}

	if foundCount == 0 {
		t.Error("Expected at least one known transposition variation")
	}
}

// Test TLDVariationTechnique
func TestTLDVariationTechnique(t *testing.T) {
	technique := &TLDVariationTechnique{}

	if technique.Name() != "tld_variation" {
		t.Errorf("Expected name 'tld_variation', got %s", technique.Name())
	}

	variations := technique.Generate("test.com", []string{"com"})

	if len(variations) == 0 {
		t.Error("Expected TLD variations to be generated")
	}

	// Check that we have different TLD variations
	tldsSeen := make(map[string]bool)
	for _, v := range variations {
		if v.Technique != "tld_variation" {
			t.Errorf("Expected technique 'tld_variation', got %s", v.Technique)
		}

		parts := strings.Split(v.Domain, ".")
		if len(parts) >= 2 {
			tld := parts[len(parts)-1]
			tldsSeen[tld] = true
		}
	}

	if len(tldsSeen) == 0 {
		t.Error("Expected multiple TLD variations")
	}

	// Check for some known suspicious TLDs
	suspiciousTLDs := []string{"tk", "ml", "ga", "cf"}
	foundSuspicious := false
	for _, tld := range suspiciousTLDs {
		if tldsSeen[tld] {
			foundSuspicious = true
			break
		}
	}

	if !foundSuspicious {
		t.Error("Expected at least one suspicious TLD variation")
	}
}

// Test SubdomainPatternTechnique
func TestSubdomainPatternTechnique(t *testing.T) {
	technique := &SubdomainPatternTechnique{}

	if technique.Name() != "subdomain_pattern" {
		t.Errorf("Expected name 'subdomain_pattern', got %s", technique.Name())
	}

	variations := technique.Generate("test.com", []string{"com"})

	if len(variations) == 0 {
		t.Error("Expected subdomain pattern variations to be generated")
	}

	// Check that we have expected patterns
	foundSecure := false
	foundLogin := false

	for _, v := range variations {
		if v.Technique != "subdomain_pattern" {
			t.Errorf("Expected technique 'subdomain_pattern', got %s", v.Technique)
		}

		if strings.Contains(v.Domain, "secure") {
			foundSecure = true
		}
		if strings.Contains(v.Domain, "login") {
			foundLogin = true
		}
	}

	if !foundSecure {
		t.Error("Expected 'secure' pattern variation")
	}
	if !foundLogin {
		t.Error("Expected 'login' pattern variation")
	}
}

// Test getBaseName utility function
func TestGetBaseName(t *testing.T) {
	tests := []struct {
		domain   string
		expected string
	}{
		{"example.com", "example"},
		{"test.co.uk", "test"},
		{"subdomain.example.org", "subdomain"},
		{"simple", "simple"},
		{"", ""},
	}

	for _, test := range tests {
		result := getBaseName(test.domain)
		if result != test.expected {
			t.Errorf("getBaseName(%s) = %s, expected %s", test.domain, result, test.expected)
		}
	}
}

// Test AvailableTechniques registry
func TestAvailableTechniques(t *testing.T) {
	expectedTechniques := []string{
		"typosquatting", "bitsquatting", "homographic",
		"insertion", "deletion", "transposition",
		"tld_variation", "subdomain_pattern",
	}

	for _, name := range expectedTechniques {
		if technique, exists := AvailableTechniques[name]; !exists {
			t.Errorf("Expected technique %s to be available", name)
		} else {
			// Test that the technique name matches the key
			if technique.Name() != name {
				t.Errorf("Technique %s reports name %s", name, technique.Name())
			}
		}
	}

	// Test that we have the expected number of techniques
	if len(AvailableTechniques) != len(expectedTechniques) {
		t.Errorf("Expected %d techniques, got %d", len(expectedTechniques), len(AvailableTechniques))
	}
}

// Test edge cases
func TestTechniquesWithEmptyDomain(t *testing.T) {
	techniques := []Technique{
		&TyposquattingTechnique{},
		&BitsquattingTechnique{},
		&HomographicTechnique{},
		&InsertionTechnique{},
		&DeletionTechnique{},
		&TranspositionTechnique{},
		&TLDVariationTechnique{},
		&SubdomainPatternTechnique{},
	}

	for _, technique := range techniques {
		variations := technique.Generate("", []string{"com"})
		// Most techniques should handle empty domain gracefully
		// Some might return empty results, others might still generate something
		if len(variations) > 100 {
			t.Errorf("Technique %s generated too many variations (%d) for empty domain",
				technique.Name(), len(variations))
		}
	}
}

func TestTechniquesWithEmptyTLDs(t *testing.T) {
	techniques := []Technique{
		&TyposquattingTechnique{},
		&BitsquattingTechnique{},
		&HomographicTechnique{},
		&InsertionTechnique{},
		&DeletionTechnique{},
		&TranspositionTechnique{},
	}

	for _, technique := range techniques {
		variations := technique.Generate("test.com", []string{})
		// Should return no variations with empty TLD list
		if len(variations) != 0 {
			t.Errorf("Technique %s should return no variations with empty TLDs, got %d",
				technique.Name(), len(variations))
		}
	}

	// TLDVariationTechnique and SubdomainPatternTechnique generate their own TLDs
	// so they should be tested separately
	tldTech := &TLDVariationTechnique{}
	tldVariations := tldTech.Generate("test.com", []string{})
	if len(tldVariations) == 0 {
		t.Error("TLDVariationTechnique should generate variations even with empty TLD list")
	}

	subdomainTech := &SubdomainPatternTechnique{}
	subVariations := subdomainTech.Generate("test.com", []string{})
	if len(subVariations) == 0 {
		t.Error("SubdomainPatternTechnique should generate variations even with empty TLD list")
	}
}

// Benchmark tests
func BenchmarkTyposquattingGenerate(b *testing.B) {
	technique := &TyposquattingTechnique{}
	for i := 0; i < b.N; i++ {
		technique.Generate("example.com", []string{"com", "net", "org"})
	}
}

func BenchmarkHomographicGenerate(b *testing.B) {
	technique := &HomographicTechnique{}
	for i := 0; i < b.N; i++ {
		technique.Generate("example.com", []string{"com", "net", "org"})
	}
}

func BenchmarkBitsquattingGenerate(b *testing.B) {
	technique := &BitsquattingTechnique{}
	for i := 0; i < b.N; i++ {
		technique.Generate("example.com", []string{"com", "net", "org"})
	}
}
