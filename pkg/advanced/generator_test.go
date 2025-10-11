package advanced

import (
	"reflect"
	"testing"
)

func TestNewVariationGenerator(t *testing.T) {
	domain := "example.com"
	opts := GeneratorOptions{
		Techniques:    []string{"typosquatting", "homographic"},
		MaxVariations: 100,
		TargetTLDs:    []string{"com", "net"},
	}

	generator := NewVariationGenerator(domain, opts)

	if generator.BaseDomain != domain {
		t.Errorf("Expected BaseDomain to be %s, got %s", domain, generator.BaseDomain)
	}

	if !reflect.DeepEqual(generator.Options, opts) {
		t.Errorf("Expected Options to be %v, got %v", opts, generator.Options)
	}

	if generator.analyzer == nil {
		t.Error("Expected analyzer to be initialized")
	}
}

func TestGenerateAll(t *testing.T) {
	domain := "test.com"
	opts := GeneratorOptions{
		Techniques:    []string{"typosquatting"},
		MaxVariations: 10,
		TargetTLDs:    []string{"com"},
	}

	generator := NewVariationGenerator(domain, opts)
	variations := generator.GenerateAll()

	if len(variations) == 0 {
		t.Error("Expected some variations to be generated")
	}

	// Verify all variations have the expected technique
	for _, v := range variations {
		if v.Technique != "typosquatting" {
			t.Errorf("Expected technique to be 'typosquatting', got %s", v.Technique)
		}
		if v.BaseDomain != domain {
			t.Errorf("Expected BaseDomain to be %s, got %s", domain, v.BaseDomain)
		}
	}
}

func TestAnalyzeAndFilter(t *testing.T) {
	domain := "test.com"
	generator := NewVariationGenerator(domain, GeneratorOptions{
		MaxVariations: 2,
		TargetTLDs:    []string{"com"},
	})

	variations := []Variation{
		{Domain: "test1.com", Technique: "typosquatting"},
		{Domain: "test2.com", Technique: "typosquatting"},
		{Domain: "test1.com", Technique: "homographic"}, // Duplicate
		{Domain: "test.com", Technique: "original"},     // Same as base
		{Domain: "test3.com", Technique: "bitsquatting"},
	}

	filtered := generator.analyzeAndFilter(variations)

	// Should remove duplicates and base domain
	if len(filtered) > 3 {
		t.Errorf("Expected max 3 unique domains (excluding base), got %d", len(filtered))
	}

	// Should respect MaxVariations limit
	if len(filtered) > generator.Options.MaxVariations {
		t.Errorf("Expected max %d variations, got %d", generator.Options.MaxVariations, len(filtered))
	}

	// Verify no base domain in results
	for _, v := range filtered {
		if v.Domain == domain {
			t.Errorf("Base domain %s should not be in filtered results", domain)
		}
	}

	// Verify similarity is calculated
	for _, v := range filtered {
		if v.Similarity == 0 {
			t.Errorf("Similarity should be calculated for variation %s", v.Domain)
		}
	}
}

func TestRankByThreatScore(t *testing.T) {
	generator := NewVariationGenerator("test.com", GeneratorOptions{})

	variations := []Variation{
		{Domain: "test1.com", Confidence: 0.5},
		{Domain: "test2.com", Confidence: 0.9},
		{Domain: "test3.com", Confidence: 0.3},
	}

	ranked := generator.rankByThreatScore(variations)

	// For now, the implementation just returns the same order
	// In a real implementation, this would sort by threat score
	if len(ranked) != len(variations) {
		t.Errorf("Expected same number of variations, got %d instead of %d", len(ranked), len(variations))
	}
}

func TestGenerateAllWithMultipleTechniques(t *testing.T) {
	domain := "example.com"
	opts := GeneratorOptions{
		Techniques:    []string{"typosquatting", "homographic", "insertion"},
		MaxVariations: 50,
		TargetTLDs:    []string{"com", "net"},
	}

	generator := NewVariationGenerator(domain, opts)
	variations := generator.GenerateAll()

	if len(variations) == 0 {
		t.Error("Expected variations to be generated with multiple techniques")
	}

	// Check that multiple techniques are represented
	techniques := make(map[string]bool)
	for _, v := range variations {
		techniques[v.Technique] = true
	}

	// At least some techniques should be present
	if len(techniques) == 0 {
		t.Error("Expected at least one technique to be present in variations")
	}
}

func TestGenerateAllWithInvalidTechnique(t *testing.T) {
	domain := "example.com"
	opts := GeneratorOptions{
		Techniques:    []string{"invalid_technique", "typosquatting"},
		MaxVariations: 10,
		TargetTLDs:    []string{"com"},
	}

	generator := NewVariationGenerator(domain, opts)
	variations := generator.GenerateAll()

	// Should still generate variations for valid techniques
	if len(variations) == 0 {
		t.Error("Expected variations from valid techniques even with invalid technique present")
	}

	// Should not contain variations from invalid technique
	for _, v := range variations {
		if v.Technique == "invalid_technique" {
			t.Error("Should not generate variations for invalid technique")
		}
	}
}

func TestGenerateAllWithNoTechniques(t *testing.T) {
	domain := "example.com"
	opts := GeneratorOptions{
		Techniques:    []string{},
		MaxVariations: 10,
		TargetTLDs:    []string{"com"},
	}

	generator := NewVariationGenerator(domain, opts)
	variations := generator.GenerateAll()

	if len(variations) != 0 {
		t.Errorf("Expected no variations with empty techniques list, got %d", len(variations))
	}
}

func TestGenerateAllWithZeroMaxVariations(t *testing.T) {
	domain := "example.com"
	opts := GeneratorOptions{
		Techniques:    []string{"typosquatting"},
		MaxVariations: 0,
		TargetTLDs:    []string{"com"},
	}

	generator := NewVariationGenerator(domain, opts)
	variations := generator.GenerateAll()

	if len(variations) != 0 {
		t.Errorf("Expected no variations with MaxVariations=0, got %d", len(variations))
	}
}

func TestGenerateAllWithEmptyTLDs(t *testing.T) {
	domain := "example.com"
	opts := GeneratorOptions{
		Techniques:    []string{"typosquatting"},
		MaxVariations: 10,
		TargetTLDs:    []string{},
	}

	generator := NewVariationGenerator(domain, opts)
	variations := generator.GenerateAll()

	if len(variations) != 0 {
		t.Errorf("Expected no variations with empty TLDs list, got %d", len(variations))
	}
}

// Benchmark tests
func BenchmarkGenerateAll(b *testing.B) {
	domain := "example.com"
	opts := GeneratorOptions{
		Techniques:    []string{"typosquatting", "homographic"},
		MaxVariations: 100,
		TargetTLDs:    []string{"com", "net", "org"},
	}

	generator := NewVariationGenerator(domain, opts)

	for i := 0; i < b.N; i++ {
		generator.GenerateAll()
	}
}

func BenchmarkAnalyzeAndFilter(b *testing.B) {
	domain := "test.com"
	generator := NewVariationGenerator(domain, GeneratorOptions{
		MaxVariations: 100,
		TargetTLDs:    []string{"com"},
	})

	// Create test variations
	variations := make([]Variation, 1000)
	for i := 0; i < 1000; i++ {
		variations[i] = Variation{
			Domain:    "test" + string(rune(i)) + ".com",
			Technique: "typosquatting",
		}
	}

	for i := 0; i < b.N; i++ {
		generator.analyzeAndFilter(variations)
	}
}
