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

// Global option to control segmentation strategy from cmd package
var spanLast3 bool

// SetSpanLast3 toggles the strategy to operate over the last 3 labels
// (label = 3rd from right, suffix = last 2 labels) when possible.
func SetSpanLast3(v bool) { spanLast3 = v }

// focusSuffix, when set, allows techniques to adjust behavior for a specific suffix
// e.g., "gov.br" to emphasize suffix impersonation.
var focusSuffix string

func SetFocusSuffix(s string) { focusSuffix = s }
