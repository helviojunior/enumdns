package advanced

type VariationGenerator struct {
	BaseDomain string
	Options    GeneratorOptions
	analyzer   *RiskAnalyzer
}

func NewVariationGenerator(domain string, opts GeneratorOptions) *VariationGenerator {
	return &VariationGenerator{
		BaseDomain: domain,
		Options:    opts,
		analyzer:   NewRiskAnalyzer(domain),
	}
}

func (vg *VariationGenerator) GenerateAll() []Variation {
	var allVariations []Variation

	// Executar cada técnica selecionada
	for _, techniqueName := range vg.Options.Techniques {
		if technique, exists := AvailableTechniques[techniqueName]; exists {
			variations := technique.Generate(vg.BaseDomain, vg.Options.TargetTLDs)
			allVariations = append(allVariations, variations...)
		}
	}

	// Analisar e filtrar resultados
	analyzed := vg.analyzeAndFilter(allVariations)

	return analyzed
}

func (vg *VariationGenerator) analyzeAndFilter(variations []Variation) []Variation {
	// Remove duplicatas
	seen := make(map[string]bool)
	var filtered []Variation

	for _, v := range variations {
		if !seen[v.Domain] && v.Domain != vg.BaseDomain {
			// Enriquecer com análise de risco
			analysis := vg.analyzer.AnalyzeVariation(v)
			v.Similarity = analysis.Similarity

			seen[v.Domain] = true
			filtered = append(filtered, v)
		}
	}

	// Limitar número de variações
	if len(filtered) > vg.Options.MaxVariations {
		// Ordenar por score de ameaça (descendente)
		// e manter apenas as mais relevantes
		filtered = vg.rankByThreatScore(filtered)[:vg.Options.MaxVariations]
	}

	return filtered
}

func (vg *VariationGenerator) rankByThreatScore(variations []Variation) []Variation {
	// Implementar ordenação por confidence + similarity
	// Para simplificar, vamos manter a ordem atual
	return variations
}
