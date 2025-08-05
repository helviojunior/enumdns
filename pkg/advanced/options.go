package advanced

type GeneratorOptions struct {
	Techniques    []string
	MaxVariations int
	TargetTLDs    []string
}

type Variation struct {
	Domain     string
	Technique  string
	Confidence float64
	Risk       string // "high", "medium", "low"
	BaseDomain string
	Similarity float64
}

type TechniqueConfig struct {
	Name       string
	Enabled    bool
	Weight     float64
	MaxResults int
}
