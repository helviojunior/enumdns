package advanced

import (
	"strings"
	"unicode"
)

// Interface para técnicas de geração
type Technique interface {
	Name() string
	Generate(domain string, tlds []string) []Variation
	GetRiskLevel() string
	GetConfidence() float64
}

// ===================
// TYPOSQUATTING
// ===================
type TyposquattingTechnique struct{}

func (t *TyposquattingTechnique) Name() string {
	return "typosquatting"
}

func (t *TyposquattingTechnique) GetRiskLevel() string {
	return "high"
}

func (t *TyposquattingTechnique) GetConfidence() float64 {
	return 0.8
}

func (t *TyposquattingTechnique) Generate(domain string, tlds []string) []Variation {
	var variations []Variation
	baseName := getBaseName(domain)

	// Mapa de teclas adjacentes no teclado QWERTY
	keyboardMap := map[rune][]rune{
		'q': {'w', 'a'}, 'w': {'q', 'e', 's'}, 'e': {'w', 'r', 'd'},
		'r': {'e', 't', 'f'}, 't': {'r', 'y', 'g'}, 'y': {'t', 'u', 'h'},
		'u': {'y', 'i', 'j'}, 'i': {'u', 'o', 'k'}, 'o': {'i', 'p', 'l'},
		'p': {'o', 'l'}, 'a': {'q', 's', 'z'}, 's': {'a', 'w', 'd', 'z', 'x'},
		'd': {'s', 'e', 'f', 'x', 'c'}, 'f': {'d', 'r', 'g', 'c', 'v'},
		'g': {'f', 't', 'h', 'v', 'b'}, 'h': {'g', 'y', 'j', 'b', 'n'},
		'j': {'h', 'u', 'k', 'n', 'm'}, 'k': {'j', 'i', 'l', 'm'},
		'l': {'k', 'o', 'p', 'm'}, 'z': {'a', 's', 'x'},
		'x': {'z', 's', 'd', 'c'}, 'c': {'x', 'd', 'f', 'v'},
		'v': {'c', 'f', 'g', 'b'}, 'b': {'v', 'g', 'h', 'n'},
		'n': {'b', 'h', 'j', 'm'}, 'm': {'n', 'j', 'k', 'l'},
	}

	for i, char := range baseName {
		if adjacentKeys, exists := keyboardMap[char]; exists {
			for _, adjKey := range adjacentKeys {
				variation := baseName[:i] + string(adjKey) + baseName[i+1:]
				for _, tld := range tlds {
					variations = append(variations, Variation{
						Domain:     variation + "." + tld,
						Technique:  t.Name(),
						Confidence: t.GetConfidence(),
						Risk:       t.GetRiskLevel(),
						BaseDomain: domain,
					})
				}
			}
		}
	}

	return variations
}

// ===================
// BITSQUATTING
// ===================
type BitsquattingTechnique struct{}

func (t *BitsquattingTechnique) Name() string {
	return "bitsquatting"
}

func (t *BitsquattingTechnique) GetRiskLevel() string {
	return "medium"
}

func (t *BitsquattingTechnique) GetConfidence() float64 {
	return 0.6
}

func (t *BitsquattingTechnique) Generate(domain string, tlds []string) []Variation {
	var variations []Variation
	baseName := getBaseName(domain)

	for i, char := range baseName {
		// Flip each bit of the character
		for bit := 0; bit < 8; bit++ {
			flipped := char ^ (1 << bit)
			if unicode.IsPrint(flipped) && flipped != char {
				variation := baseName[:i] + string(flipped) + baseName[i+1:]
				for _, tld := range tlds {
					variations = append(variations, Variation{
						Domain:     variation + "." + tld,
						Technique:  t.Name(),
						Confidence: t.GetConfidence(),
						Risk:       t.GetRiskLevel(),
						BaseDomain: domain,
					})
				}
			}
		}
	}

	return variations
}

// ===================
// HOMOGRAPHIC
// ===================
type HomographicTechnique struct{}

func (t *HomographicTechnique) Name() string {
	return "homographic"
}

func (t *HomographicTechnique) GetRiskLevel() string {
	return "high"
}

func (t *HomographicTechnique) GetConfidence() float64 {
	return 0.9
}

func (t *HomographicTechnique) Generate(domain string, tlds []string) []Variation {
	var variations []Variation
	baseName := getBaseName(domain)

	// Mapa de caracteres homográficos (visualmente similares)
	homographicMap := map[rune][]rune{
		'a': {'à', 'á', 'ä', 'â', 'ā', 'ă', 'ą', 'α'},
		'e': {'è', 'é', 'ê', 'ë', 'ē', 'ě', 'ę'},
		'i': {'ì', 'í', 'î', 'ï', 'ī', 'į', '1', 'l'},
		'o': {'ò', 'ó', 'ô', 'õ', 'ö', 'ō', 'ő', '0'},
		'u': {'ù', 'ú', 'û', 'ü', 'ū', 'ů', 'ű'},
		'c': {'ç', 'ć', 'č', 'ĉ', 'ċ'},
		'n': {'ñ', 'ń', 'ň', 'ņ'},
		's': {'š', 'ś', 'ş', 'ŝ', '$'},
		'l': {'ł', '1', 'i', '|'},
		'0': {'o', 'O', 'ο', 'о'},
		'1': {'l', 'i', 'I', '|'},
	}

	for i, char := range baseName {
		if homographs, exists := homographicMap[char]; exists {
			for _, homograph := range homographs {
				variation := baseName[:i] + string(homograph) + baseName[i+1:]
				for _, tld := range tlds {
					variations = append(variations, Variation{
						Domain:     variation + "." + tld,
						Technique:  t.Name(),
						Confidence: t.GetConfidence(),
						Risk:       t.GetRiskLevel(),
						BaseDomain: domain,
					})
				}
			}
		}
	}

	return variations
}

// ===================
// INSERTION
// ===================
type InsertionTechnique struct{}

func (t *InsertionTechnique) Name() string {
	return "insertion"
}

func (t *InsertionTechnique) GetRiskLevel() string {
	return "medium"
}

func (t *InsertionTechnique) GetConfidence() float64 {
	return 0.7
}

func (t *InsertionTechnique) Generate(domain string, tlds []string) []Variation {
	var variations []Variation
	baseName := getBaseName(domain)

	// Caracteres comuns para inserção
	commonInserts := []string{"s", "e", "r", "t", "i", "a", "n", "l", "o"}

	for i := 0; i <= len(baseName); i++ {
		for _, insert := range commonInserts {
			variation := baseName[:i] + insert + baseName[i:]
			for _, tld := range tlds {
				variations = append(variations, Variation{
					Domain:     variation + "." + tld,
					Technique:  t.Name(),
					Confidence: t.GetConfidence(),
					Risk:       t.GetRiskLevel(),
					BaseDomain: domain,
				})
			}
		}
	}

	return variations
}

// ===================
// DELETION
// ===================
type DeletionTechnique struct{}

func (t *DeletionTechnique) Name() string {
	return "deletion"
}

func (t *DeletionTechnique) GetRiskLevel() string {
	return "medium"
}

func (t *DeletionTechnique) GetConfidence() float64 {
	return 0.7
}

func (t *DeletionTechnique) Generate(domain string, tlds []string) []Variation {
	var variations []Variation
	baseName := getBaseName(domain)

	// Character deletion
	for i := 0; i < len(baseName); i++ {
		variation := baseName[:i] + baseName[i+1:]
		if len(variation) > 2 { // Avoid very short domains
			for _, tld := range tlds {
				variations = append(variations, Variation{
					Domain:     variation + "." + tld,
					Technique:  t.Name(),
					Confidence: t.GetConfidence(),
					Risk:       t.GetRiskLevel(),
					BaseDomain: domain,
				})
			}
		}
	}

	return variations
}

// ===================
// TRANSPOSITION
// ===================
type TranspositionTechnique struct{}

func (t *TranspositionTechnique) Name() string {
	return "transposition"
}

func (t *TranspositionTechnique) GetRiskLevel() string {
	return "high"
}

func (t *TranspositionTechnique) GetConfidence() float64 {
	return 0.8
}

func (t *TranspositionTechnique) Generate(domain string, tlds []string) []Variation {
	var variations []Variation
	baseName := getBaseName(domain)

	// Character transposition (swap adjacent characters)
	for i := 0; i < len(baseName)-1; i++ {
		chars := []rune(baseName)
		chars[i], chars[i+1] = chars[i+1], chars[i]
		variation := string(chars)
		for _, tld := range tlds {
			variations = append(variations, Variation{
				Domain:     variation + "." + tld,
				Technique:  t.Name(),
				Confidence: t.GetConfidence(),
				Risk:       t.GetRiskLevel(),
				BaseDomain: domain,
			})
		}
	}

	return variations
}

// ===================
// TLD VARIATIONS
// ===================
type TLDVariationTechnique struct{}

func (t *TLDVariationTechnique) Name() string {
	return "tld_variation"
}

func (t *TLDVariationTechnique) GetRiskLevel() string {
	return "medium"
}

func (t *TLDVariationTechnique) GetConfidence() float64 {
	return 0.6
}

func (t *TLDVariationTechnique) Generate(domain string, tlds []string) []Variation {
	var variations []Variation
	baseName := getBaseName(domain)

	// TLDs comuns para phishing
	phishingTLDs := []string{
		"tk", "ml", "ga", "cf", "gq", // Freenom TLDs
		"co", "cc", "biz", "info", "org", "net",
		"io", "me", "ly", "to", "ws", "click", "download",
	}

	for _, tld := range phishingTLDs {
		variations = append(variations, Variation{
			Domain:     baseName + "." + tld,
			Technique:  t.Name(),
			Confidence: t.GetConfidence(),
			Risk:       t.GetRiskLevel(),
			BaseDomain: domain,
		})
	}

	return variations
}

// ===================
// SUBDOMAIN PATTERNS
// ===================
type SubdomainPatternTechnique struct{}

func (t *SubdomainPatternTechnique) Name() string {
	return "subdomain_pattern"
}

func (t *SubdomainPatternTechnique) GetRiskLevel() string {
	return "high"
}

func (t *SubdomainPatternTechnique) GetConfidence() float64 {
	return 0.8
}

func (t *SubdomainPatternTechnique) Generate(domain string, tlds []string) []Variation {
	var variations []Variation
	baseName := getBaseName(domain)

	// Padrões comuns de phishing
	phishingPrefixes := []string{
		"secure", "login", "account", "verify", "update",
		"confirm", "support", "help", "mail", "www",
		"app", "mobile", "portal", "admin", "api",
		"my", "user", "client", "member", "auth",
	}

	for _, prefix := range phishingPrefixes {
		// Formato: prefix + domain
		variations = append(variations, Variation{
			Domain:     prefix + baseName + ".com",
			Technique:  t.Name(),
			Confidence: t.GetConfidence(),
			Risk:       t.GetRiskLevel(),
			BaseDomain: domain,
		})

		// Formato: prefix-domain
		variations = append(variations, Variation{
			Domain:     prefix + "-" + baseName + ".com",
			Technique:  t.Name(),
			Confidence: t.GetConfidence(),
			Risk:       t.GetRiskLevel(),
			BaseDomain: domain,
		})
	}

	return variations
}

// ===================
// UTILITY FUNCTIONS
// ===================
func getBaseName(domain string) string {
	parts := strings.Split(domain, ".")
	return parts[0]
}

// Registry de técnicas disponíveis
var AvailableTechniques = map[string]Technique{
	"typosquatting":     &TyposquattingTechnique{},
	"bitsquatting":      &BitsquattingTechnique{},
	"homographic":       &HomographicTechnique{},
	"insertion":         &InsertionTechnique{},
	"deletion":          &DeletionTechnique{},
	"transposition":     &TranspositionTechnique{},
	"tld_variation":     &TLDVariationTechnique{},
	"subdomain_pattern": &SubdomainPatternTechnique{},
}
