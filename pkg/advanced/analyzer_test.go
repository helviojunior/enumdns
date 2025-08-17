package advanced

import (
	"testing"
)

func TestNewRiskAnalyzer(t *testing.T) {
	baseDomain := "example.com"
	analyzer := NewRiskAnalyzer(baseDomain)

	if analyzer.BaseDomain != baseDomain {
		t.Errorf("Expected BaseDomain to be %s, got %s", baseDomain, analyzer.BaseDomain)
	}
}

func TestCalculateSimilarity(t *testing.T) {
	analyzer := NewRiskAnalyzer("example.com")

	tests := []struct {
		name     string
		domain   string
		expected float64
		tolerance float64
	}{
		{
			name:      "Identical domain",
			domain:    "example.com",
			expected:  1.0,
			tolerance: 0.0,
		},
		{
			name:      "Similar domain (typo)",
			domain:    "exanple.com",
			expected:  0.87, // approximately
			tolerance: 0.05,
		},
		{
			name:      "Different TLD",
			domain:    "example.org",
			expected:  1.0, // same base name
			tolerance: 0.05,
		},
		{
			name:      "Completely different",
			domain:    "google.com",
			expected:  0.29, // approximately
			tolerance: 0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.calculateSimilarity(tt.domain)
			if result < tt.expected-tt.tolerance || result > tt.expected+tt.tolerance {
				t.Errorf("calculateSimilarity() = %v, expected around %v ± %v", result, tt.expected, tt.tolerance)
			}
		})
	}
}

func TestCalculateThreatScore(t *testing.T) {
	analyzer := NewRiskAnalyzer("example.com")

	tests := []struct {
		name      string
		variation Variation
		minScore  float64
		maxScore  float64
	}{
		{
			name: "High confidence homographic",
			variation: Variation{
				Domain:     "еxample.com", // Cyrillic 'е'
				Technique:  "homographic",
				Confidence: 0.8,
			},
			minScore: 0.8,
			maxScore: 1.0,
		},
		{
			name: "Low confidence bitsquatting",
			variation: Variation{
				Domain:     "fxample.com",
				Technique:  "bitsquatting",
				Confidence: 0.3,
			},
			minScore: 0.3,
			maxScore: 0.5,
		},
		{
			name: "Suspicious TLD domain",
			variation: Variation{
				Domain:     "example.tk",
				Technique:  "tld_variation",
				Confidence: 0.5,
			},
			minScore: 0.7,
			maxScore: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := analyzer.calculateThreatScore(tt.variation)
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("calculateThreatScore() = %v, expected between %v and %v", score, tt.minScore, tt.maxScore)
			}
		})
	}
}

func TestIdentifyThreatIndicators(t *testing.T) {
	analyzer := NewRiskAnalyzer("example.com")

	tests := []struct {
		name               string
		variation          Variation
		expectedIndicators []string
	}{
		{
			name: "Suspicious TLD",
			variation: Variation{
				Domain:    "example.tk",
				Technique: "tld_variation",
			},
			expectedIndicators: []string{"suspicious_tld"},
		},
		{
			name: "Phishing pattern",
			variation: Variation{
				Domain:    "secure-example.com",
				Technique: "insertion",
			},
			expectedIndicators: []string{"phishing_pattern"},
		},
		{
			name: "Unicode tricks",
			variation: Variation{
				Domain:    "еxample.com", // Cyrillic 'е'
				Technique: "homographic",
			},
			expectedIndicators: []string{"unicode_tricks"},
		},
		{
			name: "Multiple indicators",
			variation: Variation{
				Domain:    "secure-еxample.tk", // Phishing + Unicode + Suspicious TLD
				Technique: "homographic",
			},
			expectedIndicators: []string{"suspicious_tld", "phishing_pattern", "unicode_tricks"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indicators := analyzer.identifyThreatIndicators(tt.variation)
			
			// Check if all expected indicators are present
			for _, expected := range tt.expectedIndicators {
				found := false
				for _, indicator := range indicators {
					if indicator == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected indicator %s not found in %v", expected, indicators)
				}
			}
		})
	}
}

func TestAnalyzeVariation(t *testing.T) {
	analyzer := NewRiskAnalyzer("example.com")
	
	variation := Variation{
		Domain:     "еxample.com", // Cyrillic 'е'
		Technique:  "homographic",
		Confidence: 0.8,
	}

	result := analyzer.AnalyzeVariation(variation)

	if result.Variation != variation {
		t.Errorf("Expected variation to be preserved")
	}

	if result.Similarity <= 0 || result.Similarity > 1 {
		t.Errorf("Similarity should be between 0 and 1, got %v", result.Similarity)
	}

	if result.ThreatScore <= 0 || result.ThreatScore > 1 {
		t.Errorf("ThreatScore should be between 0 and 1, got %v", result.ThreatScore)
	}

	if len(result.Indicators) == 0 {
		t.Errorf("Expected some threat indicators for this variation")
	}
}

func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		name     string
		s1       string
		s2       string
		expected int
	}{
		{
			name:     "Identical strings",
			s1:       "test",
			s2:       "test",
			expected: 0,
		},
		{
			name:     "Single substitution",
			s1:       "test",
			s2:       "best",
			expected: 1,
		},
		{
			name:     "Single insertion",
			s1:       "test",
			s2:       "tests",
			expected: 1,
		},
		{
			name:     "Single deletion",
			s1:       "tests",
			s2:       "test",
			expected: 1,
		},
		{
			name:     "Empty strings",
			s1:       "",
			s2:       "",
			expected: 0,
		},
		{
			name:     "One empty string",
			s1:       "test",
			s2:       "",
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := levenshteinDistance(tt.s1, tt.s2)
			if result != tt.expected {
				t.Errorf("levenshteinDistance(%s, %s) = %d, expected %d", tt.s1, tt.s2, result, tt.expected)
			}
		})
	}
}

func TestIsSuspiciousTLD(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		expected bool
	}{
		{
			name:     "Freenom TLD (.tk)",
			domain:   "example.tk",
			expected: true,
		},
		{
			name:     "Suspicious TLD (.click)",
			domain:   "example.click",
			expected: true,
		},
		{
			name:     "Normal TLD (.com)",
			domain:   "example.com",
			expected: false,
		},
		{
			name:     "Normal TLD (.org)",
			domain:   "example.org",
			expected: false,
		},
		{
			name:     "Case insensitive (.TK)",
			domain:   "example.TK",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSuspiciousTLD(tt.domain)
			if result != tt.expected {
				t.Errorf("isSuspiciousTLD(%s) = %v, expected %v", tt.domain, result, tt.expected)
			}
		})
	}
}

func TestContainsPhishingPatterns(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		expected bool
	}{
		{
			name:     "Contains 'secure'",
			domain:   "secure-example.com",
			expected: true,
		},
		{
			name:     "Contains 'login'",
			domain:   "example-login.com",
			expected: true,
		},
		{
			name:     "Contains 'paypal'",
			domain:   "paypal-example.com",
			expected: true,
		},
		{
			name:     "No phishing patterns",
			domain:   "example.com",
			expected: false,
		},
		{
			name:     "Case insensitive",
			domain:   "SECURE-example.com",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsPhishingPatterns(tt.domain)
			if result != tt.expected {
				t.Errorf("containsPhishingPatterns(%s) = %v, expected %v", tt.domain, result, tt.expected)
			}
		})
	}
}

func TestContainsUnicodeTricks(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		expected bool
	}{
		{
			name:     "Cyrillic character",
			domain:   "еxample.com", // Cyrillic 'е'
			expected: true,
		},
		{
			name:     "Greek character",
			domain:   "αpple.com", // Greek 'α'
			expected: true,
		},
		{
			name:     "Normal ASCII",
			domain:   "example.com",
			expected: false,
		},
		{
			name:     "Mixed unicode",
			domain:   "gοοgle.com", // Mixed with Greek 'ο'
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsUnicodeTricks(tt.domain)
			if result != tt.expected {
				t.Errorf("containsUnicodeTricks(%s) = %v, expected %v", tt.domain, result, tt.expected)
			}
		})
	}
}

func TestMinFunction(t *testing.T) {
	tests := []struct {
		name     string
		a, b, c  int
		expected int
	}{
		{
			name:     "First is minimum",
			a:        1,
			b:        2,
			c:        3,
			expected: 1,
		},
		{
			name:     "Second is minimum",
			a:        3,
			b:        1,
			c:        2,
			expected: 1,
		},
		{
			name:     "Third is minimum",
			a:        3,
			b:        2,
			c:        1,
			expected: 1,
		},
		{
			name:     "All equal",
			a:        2,
			b:        2,
			c:        2,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := min(tt.a, tt.b, tt.c)
			if result != tt.expected {
				t.Errorf("min(%d, %d, %d) = %d, expected %d", tt.a, tt.b, tt.c, result, tt.expected)
			}
		})
	}
}

// Benchmark tests for performance
func BenchmarkLevenshteinDistance(b *testing.B) {
	s1 := "example"
	s2 := "exanple"
	
	for i := 0; i < b.N; i++ {
		levenshteinDistance(s1, s2)
	}
}

func BenchmarkAnalyzeVariation(b *testing.B) {
	analyzer := NewRiskAnalyzer("example.com")
	variation := Variation{
		Domain:     "еxample.com",
		Technique:  "homographic",
		Confidence: 0.8,
	}
	
	for i := 0; i < b.N; i++ {
		analyzer.AnalyzeVariation(variation)
	}
}