package advanced

import (
	"strings"
	// Removido: "unicode/utf8" - não estava sendo usado
)

// Analisador de risco e similaridade
type RiskAnalyzer struct {
	BaseDomain string
}

type AnalysisResult struct {
	Variation   Variation
	Similarity  float64
	ThreatScore float64
	Indicators  []string
}

func NewRiskAnalyzer(baseDomain string) *RiskAnalyzer {
	return &RiskAnalyzer{
		BaseDomain: baseDomain,
	}
}

func (ra *RiskAnalyzer) AnalyzeVariation(variation Variation) AnalysisResult {
	result := AnalysisResult{
		Variation:  variation,
		Indicators: []string{},
	}

	// Calcular similaridade
	result.Similarity = ra.calculateSimilarity(variation.Domain)

	// Calcular score de ameaça
	result.ThreatScore = ra.calculateThreatScore(variation)

	// Identificar indicadores específicos
	result.Indicators = ra.identifyThreatIndicators(variation)

	return result
}

func (ra *RiskAnalyzer) calculateSimilarity(domain string) float64 {
	baseName := getBaseName(ra.BaseDomain)
	targetName := getBaseName(domain)

	return levenshteinSimilarity(baseName, targetName)
}

func (ra *RiskAnalyzer) calculateThreatScore(variation Variation) float64 {
	score := variation.Confidence

	// Ajustar baseado na técnica
	switch variation.Technique {
	case "homographic", "typosquatting":
		score += 0.2
	case "subdomain_pattern":
		score += 0.3
	case "bitsquatting":
		score += 0.1
	}

	// Ajustar baseado no TLD
	if isSuspiciousTLD(variation.Domain) {
		score += 0.2
	}

	// Normalizar para 0-1
	if score > 1.0 {
		score = 1.0
	}

	return score
}

func (ra *RiskAnalyzer) identifyThreatIndicators(variation Variation) []string {
	var indicators []string

	// Verificar TLDs suspeitos
	if isSuspiciousTLD(variation.Domain) {
		indicators = append(indicators, "suspicious_tld")
	}

	// Verificar padrões de phishing
	if containsPhishingPatterns(variation.Domain) {
		indicators = append(indicators, "phishing_pattern")
	}

	// Verificar alta similaridade
	if ra.calculateSimilarity(variation.Domain) > 0.8 {
		indicators = append(indicators, "high_similarity")
	}

	// Verificar caracteres Unicode suspeitos
	if containsUnicodeTricks(variation.Domain) {
		indicators = append(indicators, "unicode_tricks")
	}

	return indicators
}

// ===================
// UTILITY FUNCTIONS
// ===================

func levenshteinSimilarity(s1, s2 string) float64 {
	distance := levenshteinDistance(s1, s2)
	maxLen := len(s1)
	if len(s2) > maxLen {
		maxLen = len(s2)
	}

	if maxLen == 0 {
		return 1.0
	}

	return 1.0 - float64(distance)/float64(maxLen)
}

func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}

	for j := range matrix[0] {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}

func isSuspiciousTLD(domain string) bool {
	suspiciousTLDs := []string{
		".tk", ".ml", ".ga", ".cf", ".gq", // Freenom
		".click", ".download", ".stream", ".review",
		".top", ".win", ".date", ".racing",
	}

	domain = strings.ToLower(domain)
	for _, tld := range suspiciousTLDs {
		if strings.HasSuffix(domain, tld) {
			return true
		}
	}
	return false
}

func containsPhishingPatterns(domain string) bool {
	phishingPatterns := []string{
		"secure", "login", "verify", "update", "confirm",
		"account", "banking", "paypal", "amazon", "apple",
		"microsoft", "google", "facebook", "twitter",
	}

	domain = strings.ToLower(domain)
	for _, pattern := range phishingPatterns {
		if strings.Contains(domain, pattern) {
			return true
		}
	}
	return false
}

func containsUnicodeTricks(domain string) bool {
	// Verificar se contém caracteres não-ASCII
	for _, r := range domain {
		if r > 127 {
			return true
		}
	}

	// Verificar combinações suspeitas
	suspiciousCombos := []string{
		"а", "е", "о", "р", "с", "х", "у", // Cyrillic que se parecem com ASCII
		"ο", "α", "ρ", "τ", "υ", "ι", // Greek
	}

	for _, combo := range suspiciousCombos {
		if strings.Contains(domain, combo) {
			return true
		}
	}

	return false
}
