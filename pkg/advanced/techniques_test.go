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

	variations := technique.Generate("test.com", []string{"com", "net"})

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

	variations := technique.Generate("abc.com", []string{"com"})

	if len(variations) == 0 {
		t.Error("Expected insertion variations to be generated")
	}

	// Check that variations are longer than original
	originalBase := "abc"
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

	// We expect at least one alternative TLD generated
	if len(tldsSeen) == 0 {
		t.Error("Expected at least one alternative TLD variation")
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
		{"subdomain.example.org", "example"},
		{"simple", ""},
		{"", ""},
	}

	for _, test := range tests {
		result := getBaseName(test.domain)
		if result != test.expected {
			t.Errorf("getBaseName(%s) = %s, expected %s", test.domain, result, test.expected)
		}
	}
}

func TestGetBaseNameMultiLabelSuffix(t *testing.T) {
	domain := "recife.pe.gov.br"
	expected := "pe"
	got := getBaseName(domain)
	if got != expected {
		t.Errorf("getBaseName(%s) = %s, expected %s", domain, got, expected)
	}
}

// Test AvailableTechniques registry
func TestAvailableTechniques(t *testing.T) {
	expectedTechniques := []string{
		"typosquatting", "bitsquatting", "homographic",
		"insertion", "deletion", "transposition",
		"tld_variation", "subdomain_pattern", "suffix_impersonation",
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

func TestSpanLast3Extraction(t *testing.T) {
	// Ensure default off
	SetSpanLast3(false)
	lbl, sfx := getLabelAndSuffix("updates.microsoft.com")
	if lbl != "microsoft" || sfx != "com" {
		t.Errorf("Expected PSL extraction (microsoft, com); got (%s, %s)", lbl, sfx)
	}
	// Enable last3 and validate
	SetSpanLast3(true)
	lbl, sfx = getLabelAndSuffix("updates.microsoft.com")
	if lbl != "updates" || sfx != "microsoft.com" {
		t.Errorf("Expected last3 extraction (updates, microsoft.com); got (%s, %s)", lbl, sfx)
	}
	// Reset
	SetSpanLast3(false)
}

func TestSuffixImpersonationGovBR(t *testing.T) {
	tech := &SuffixImpersonationTechnique{}
	vars := tech.Generate("recife.pe.gov.br", []string{})
	if len(vars) == 0 {
		t.Error("Expected suffix impersonation variations for gov.br")
	}
	found := false
	for _, v := range vars {
		if strings.Contains(v.Domain, ".g0v.br") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected at least one g0v.br impersonation candidate")
	}
}

func TestSuffixImpersonationFocusSuffix(t *testing.T) {
	tech := &SuffixImpersonationTechnique{}
	// Without focus, non-gov should not produce variations
	vars := tech.Generate("example.com", []string{})
	if len(vars) != 0 {
		t.Error("Expected no variations for non-gov.br without focus")
	}
	// With focus gov.br, allow generation over non-gov domain (e.g., example.g0v.br)
	SetFocusSuffix("gov.br")
	defer SetFocusSuffix("")
	vars = tech.Generate("example.com", []string{})
	if len(vars) == 0 {
		t.Error("Expected variations for non-gov domain when focus-suffix=gov.br")
	}
	seen := false
	for _, v := range vars {
		if strings.HasSuffix(v.Domain, ".br") {
			seen = true
			break
		}
	}
	if !seen {
		t.Error("Expected .br impersonation domains when focus-suffix=gov.br")
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
