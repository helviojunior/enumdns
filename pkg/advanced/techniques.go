package advanced

import (
	"strings"
	"unicode"

	"golang.org/x/net/publicsuffix"
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
    if !labelMutationAllowed(domain, baseName) {
        return variations
    }
    if len(baseName) < 3 {
        return variations
    }

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
    if !labelMutationAllowed(domain, baseName) {
        return variations
    }
    if len(baseName) < 3 {
        return variations
    }

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
    if !labelMutationAllowed(domain, baseName) {
        return variations
    }
    if len(baseName) < 3 {
        return variations
    }

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
		'g': {'q', '9', 'ğ', 'ģ'},
		'd': {'ð', 'đ', 'ď'},
		'v': {'ν', 'υ', 'ѵ'},
		'x': {'×', 'χ', 'х'},
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
    if !labelMutationAllowed(domain, baseName) {
        return variations
    }
    if len(baseName) < 3 {
        return variations
    }

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
    if !labelMutationAllowed(domain, baseName) {
        return variations
    }
    if len(baseName) < 3 {
        return variations
    }

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
    if !labelMutationAllowed(domain, baseName) {
        return variations
    }
    if len(baseName) < 3 {
        return variations
    }

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

// ===================
// SUFFIX IMPERSONATION (e.g., gov.br -> g0v.br)
// ===================
type SuffixImpersonationTechnique struct{}

func (t *SuffixImpersonationTechnique) Name() string {
    return "suffix_impersonation"
}

func (t *SuffixImpersonationTechnique) GetRiskLevel() string {
    return "high"
}

func (t *SuffixImpersonationTechnique) GetConfidence() float64 {
    return 0.8
}

func (t *SuffixImpersonationTechnique) Generate(domain string, tlds []string) []Variation {
    var variations []Variation
    d := strings.TrimSuffix(strings.ToLower(domain), ".")
    if d == "" {
        return variations
    }
    label, suffix := getLabelAndSuffix(d)
    if label == "" || suffix == "" {
        return variations
    }
    // Only for gov.br by default unless focusSuffix matches
    if suffix != "gov.br" && focusSuffix != "gov.br" {
        return variations
    }
    // Extract left part (subdomains), everything before label.suffix
    tail := label + "." + suffix
    left := strings.TrimSuffix(d, tail)
    left = strings.TrimSuffix(left, ".")

    // Variations of "gov"
    candidates := map[string]struct{}{
        // direct subs
        "g0v": {}, "gov": {}, "gøv": {}, "gоv": {}, // 'о' Cyrillic
        "qov": {}, "hov": {}, "gou": {}, "goy": {},
        // insertions/deletions/duplications
        "gv": {}, "goov": {}, "govv": {}, "g-ov": {}, "go-v": {},
    }

    seen := map[string]struct{}{}
    for variant := range candidates {
        ns := variant + ".br"
        // Compose full domain: [left].label.[variant].br
        var full string
        if left != "" {
            full = left + "." + label + "." + ns
        } else {
            full = label + "." + ns
        }
        if _, ok := seen[full]; ok {
            continue
        }
        seen[full] = struct{}{}
        variations = append(variations, Variation{
            Domain:     full,
            Technique:  t.Name(),
            Confidence: t.GetConfidence(),
            Risk:       t.GetRiskLevel(),
            BaseDomain: domain,
        })
    }

    return variations
}

func (t *TLDVariationTechnique) GetRiskLevel() string {
	return "medium"
}

func (t *TLDVariationTechnique) GetConfidence() float64 {
	return 0.6
}

func (t *TLDVariationTechnique) Generate(domain string, tlds []string) []Variation {
    var variations []Variation
    base := strings.TrimSuffix(strings.ToLower(domain), ".")
    if base == "" {
        return variations
    }

    baseName := getBaseName(base)
    if baseName == "" {
        return variations
    }

    // Prefer provided target TLDs; fall back to a conservative default list
    targetTLDs := tlds
    if len(targetTLDs) == 0 {
        targetTLDs = []string{
            "com", "net", "org", "co", "info", "biz", "io", "me",
            "tk", "ml", "ga", "cf", "gq", "pw", "to", "ws", "ly",
        }
    }

    for _, tld := range targetTLDs {
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

type ThreatIndicators struct {
	SuspiciousTLD      bool    // .tk, .ml, etc.
	PhishingKeywords   bool    // secure-, login-, etc.
	RecentRegistration bool    // < 30 days
	UnicodeUsage       bool    // Non-ASCII chars
	HighSimilarity     float64 // Levenshtein distance
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
    label, _ := getLabelAndSuffix(domain)
    return label
}

// getLabelAndSuffix returns the registrable label and its effective suffix,
// with overrides for compound bases like gov.br and com.br.
// Examples:
// - recife.pe.gov.br -> ("pe", "gov.br")
// - example.com      -> ("example", "com")
// - sub.example.org  -> ("example", "org")
func getLabelAndSuffix(domain string) (string, string) {
    d := strings.TrimSuffix(strings.ToLower(domain), ".")
    if d == "" {
        return "", ""
    }
    // spanLast3 strategy: take the last 3 labels when available
    if spanLast3 {
        parts := strings.Split(d, ".")
        if len(parts) >= 3 {
            label := parts[len(parts)-3]
            suffix := parts[len(parts)-2] + "." + parts[len(parts)-1]
            return label, suffix
        }
        // Fallback to PSL below for shorter domains
    }
    etld1, err := publicsuffix.EffectiveTLDPlusOne(d)
    if err != nil || etld1 == "" {
        return "", ""
    }
    parts := strings.Split(etld1, ".")
    if len(parts) < 2 {
        return "", ""
    }

    // Default PSL-based label and suffix
    label := parts[0]
    suffix := strings.Join(parts[1:], ".")

    // Overrides for compound base suffixes where organizational level precedes a known base
    // e.g., X.Y.gov.br -> treat Y as label and gov.br as suffix
    compoundBases := map[string]struct{}{
        "gov.br": {},
        "com.br": {},
        "net.br": {},
        "org.br": {},
        "com.ar": {},
        "com.mx": {},
    }

    if len(parts) >= 3 {
        base2 := parts[len(parts)-2] + "." + parts[len(parts)-1]
        if _, ok := compoundBases[base2]; ok {
            label = parts[len(parts)-3]
            suffix = base2
        }
    }

    return label, suffix
}

// labelMutationAllowed controls whether to mutate the registrable label based on suffix rules.
// - For gov.br: avoid mutating very short labels (length < 3) to reduce noise; focus on suffix impersonation instead.
// - For com.br and single-label suffixes (e.g., com, net, org): allow mutation even for 2-char labels.
func labelMutationAllowed(domain, label string) bool {
    if label == "" {
        return false
    }
    _, suffix := getLabelAndSuffix(domain)
    if suffix == "gov.br" {
        return len(label) >= 3
    }
    // Default allow if label has at least 2 chars
    return len(label) >= 2
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
    "suffix_impersonation": &SuffixImpersonationTechnique{},
}
