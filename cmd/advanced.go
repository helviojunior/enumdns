package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/bob-reis/enumdns/internal/ascii"
	"github.com/bob-reis/enumdns/internal/tools"
	"github.com/bob-reis/enumdns/pkg/advanced"
	"github.com/bob-reis/enumdns/pkg/database"
	"github.com/bob-reis/enumdns/pkg/log"
	"github.com/bob-reis/enumdns/pkg/models"
	"github.com/bob-reis/enumdns/pkg/readers"
	"github.com/bob-reis/enumdns/pkg/runner"
	"github.com/bob-reis/enumdns/pkg/writers"
	"github.com/spf13/cobra"
	"golang.org/x/net/publicsuffix"
	"gorm.io/gorm"
)

// Opções específicas do módulo threat-analysis
var threatAnalysisOpts = struct {
	Typosquatting  bool
	Bitsquatting   bool
	Homographic    bool
	AllTechniques  bool
	MaxVariations  int
	TLDs           []string
	ShowProgress   bool
	WarnLimits     bool
	BrandCombo     bool
	SpanLast3      bool
	EmitCandidates bool
	FocusSuffix    string
}{
	MaxVariations: 1000,
	TLDs:          []string{"com", "net", "org", "co", "io", "com.br", "net.br", "org.br"},
	ShowProgress:  true,
	WarnLimits:    true,
}

var advancedWriters = []writers.Writer{}

var threatAnalysisCmd = &cobra.Command{
	Use:   "threat-analysis",
	Short: "Advanced domain threat analysis for typosquatting and malicious domains",
	Long: ascii.LogoHelp(ascii.Markdown(`
# threat-analysis

Perform comprehensive domain threat analysis including typosquatting detection,
homographic attacks, bitsquatting, and suspicious domain discovery.

This module generates domain variations using multiple advanced techniques and 
verifies their existence to identify potential security threats and malicious domains.
`)),
	Example: `
   - enumdns threat-analysis -d example.com --typosquatting -o threats.txt
   - enumdns threat-analysis -d example.com --all-techniques --write-db
   - enumdns threat-analysis -L domains.txt --bitsquatting --write-jsonl
   - enumdns threat-analysis -d recife.pe.gov.br --all-techniques --focus-suffix=gov.br --emit-candidates -o gov-br.txt
   - enumdns threat-analysis -d yeslinux.com.br --all-techniques --target-tlds com,net,org,co,info,io,com.br,net.br,org.br
   - enumdns threat-analysis -d updates.microsoft.com --all-techniques --span-last3`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		// Chamar o PersistentPreRunE do comando pai
		if err = rootCmd.PersistentPreRunE(cmd, args); err != nil {
			return err
		}

		if opts.PrivateDns {
			log.Warnf("DNS server: %s (private)", fileOptions.DnsServer)
		} else {
			log.Warn("DNS server: " + fileOptions.DnsServer)
		}

		// Configure writers using the common function
		if err := ConfigureWriters(&advancedWriters); err != nil {
			return err
		}

		fileOptions.DnsServer = opts.DnsServer + ":" + fmt.Sprintf("%d", opts.DnsPort)

		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Validações específicas do advanced
		if opts.DnsSuffix == "" && fileOptions.DnsSuffixFile == "" {
			return errors.New("a domain (-d) or domain list file (-L) must be specified")
		}

		if fileOptions.DnsSuffixFile != "" {
			if !tools.FileExists(fileOptions.DnsSuffixFile) {
				return errors.New("domain list file is not readable")
			}
		}

		// Verificar se pelo menos uma técnica foi selecionada
		if !threatAnalysisOpts.Typosquatting && !threatAnalysisOpts.Bitsquatting &&
			!threatAnalysisOpts.Homographic && !threatAnalysisOpts.AllTechniques {
			log.Warn("No specific techniques selected, enabling all techniques")
			threatAnalysisOpts.AllTechniques = true
		}

		// Validar limite de variações
		if threatAnalysisOpts.MaxVariations > 10000 {
			log.Warnf("MaxVariations (%d) is very high and may impact performance. Consider using a lower value.", threatAnalysisOpts.MaxVariations)
		}

		if threatAnalysisOpts.MaxVariations < 10 {
			log.Warnf("MaxVariations (%d) is very low and may miss important variations. Consider using a higher value.", threatAnalysisOpts.MaxVariations)
		}

		// Pass segmentation preference to advanced package
		advanced.SetSpanLast3(threatAnalysisOpts.SpanLast3)
		if threatAnalysisOpts.FocusSuffix != "" {
			advanced.SetFocusSuffix(strings.Trim(strings.ToLower(threatAnalysisOpts.FocusSuffix), ". "))
		}
		return nil
	},
	Run: runThreatAnalysis,
}

func validateDNSConnectivity() {
	_, err := tools.GetValidDnsSuffix(fileOptions.DnsServer, "google.com.", opts.Proxy)
	if err != nil {
		log.Error("Error checking DNS connectivity", "err", err)
		os.Exit(2)
	}
}

func loadAndValidateDomains() []string {
	domains := []string{}

	if fileOptions.DnsSuffixFile != "" {
		log.Debugf("Reading domain list file: %s", fileOptions.DnsSuffixFile)
		reader := readers.NewFileReader(fileOptions)
		if err := reader.ReadDnsList(&domains); err != nil {
			log.Error("error reading domain list", "err", err)
			os.Exit(2)
		}
	} else {
		s, err := tools.GetValidDnsSuffix(fileOptions.DnsServer, opts.DnsSuffix, opts.Proxy)
		if err != nil {
			log.Error("invalid domain", "domain", opts.DnsSuffix, "err", err)
			os.Exit(2)
		}
		domains = append(domains, s)
	}

	if len(domains) == 0 {
		log.Error("No valid domains to analyze")
		os.Exit(2)
	}

	return domains
}

func generateVariations(domains []string) []advanced.Variation {
	techniques := getEnabledTechniques()
	log.Infof("Enabled techniques: %v", techniques)

	allVariations := []advanced.Variation{}

	for _, domain := range domains {
		log.Infof("Generating variations for %s", domain)

		// Determine allowed TLDs for this domain: swap only if suffix has a single label
		allowedTLDs := computeAllowedTLDs(domain, threatAnalysisOpts.TLDs)

		generator := advanced.NewVariationGenerator(domain, advanced.GeneratorOptions{
			Techniques:    techniques,
			MaxVariations: threatAnalysisOpts.MaxVariations,
			TargetTLDs:    allowedTLDs,
		})

		variations := generator.GenerateAll()
		log.Infof("Generated %d variations for %s", len(variations), domain)
		allVariations = append(allVariations, variations...)
	}

	total := len(allVariations)
	log.Infof("Total variations to check: %s", tools.FormatInt(total))

	if threatAnalysisOpts.WarnLimits && total >= threatAnalysisOpts.MaxVariations*len(domains) {
		log.Warnf("Maximum variation limit reached (%d per domain). Some variations may have been excluded.", threatAnalysisOpts.MaxVariations)
		log.Warnf("To analyze more variations, increase --max-variations parameter (current: %d)", threatAnalysisOpts.MaxVariations)
	}

	return allVariations
}

func executeAnalysis(allVariations []advanced.Variation) {
	runnerLogger := slog.New(log.Logger)
	runnerInstance, err := runner.NewRunner(runnerLogger, *opts, advancedWriters)
	if err != nil {
		log.Error("Failed to create runner", "err", err)
		os.Exit(2)
	}

	conn, err := database.Connection(opts.Writer.CtrlDbURI, true, false)
	if err != nil {
		log.Error("Error establishing connection with Database", "err", err)
		os.Exit(2)
	}

	total := len(allVariations)
	feedTargetsToRunner(runnerInstance, conn, allVariations)

	status := runnerInstance.Run(total)
	runnerInstance.Close()

	if status.Skiped > 0 {
		log.Warnf("%d variations were skipped because they were already analyzed. Use the --force parameter to reanalyze them.", status.Skiped)
		ascii.ClearLine()
	}

	log.Info("Threat analysis completed")
}

func feedTargetsToRunner(runnerInstance *runner.Runner, conn *gorm.DB, allVariations []advanced.Variation) {
	go func() {
		defer close(runnerInstance.Targets)
		ascii.HideCursor()

		for _, variation := range allVariations {
			host := variation.Domain
			if !strings.HasSuffix(host, ".") {
				host = host + "."
			}
			shouldProcess := forceCheck

			if !forceCheck {
				if response := conn.Raw("SELECT count(id) as count from results WHERE failed = 0 AND fqdn = ?", host); response != nil {
					var cnt int
					_ = response.Row().Scan(&cnt)
					shouldProcess = (cnt == 0)
					if cnt > 0 {
						log.Debug("[Host already checked]", "fqdn", host)
					}
				}
			}

			if shouldProcess {
				runnerInstance.Targets <- host
			} else {
				runnerInstance.AddSkiped()
			}
		}
	}()
}

func runThreatAnalysis(cmd *cobra.Command, args []string) {
	validateDNSConnectivity()

	log.Debug("starting domain threat analysis")

	domains := loadAndValidateDomains()
	log.Infof("Analyzing %d domain(s) for threats", len(domains))

	allVariations := generateVariations(domains)
	if len(allVariations) == 0 {
		log.Warn("No variations generated")
		return
	}

	// Optionally emit all generated candidates to writers before probing
	if threatAnalysisOpts.EmitCandidates {
		now := time.Now()
		for _, w := range advancedWriters {
			for _, v := range allVariations {
				fq := &models.FQDNData{FQDN: v.Domain, Source: "Generated", ProbedAt: now}
				_ = w.WriteFqdn(fq)
			}
		}
	}

	executeAnalysis(allVariations)
}

func getEnabledTechniques() []string {
	var techniques []string

	if threatAnalysisOpts.AllTechniques {
		techniques = []string{"typosquatting", "bitsquatting", "homographic", "insertion", "deletion", "transposition", "tld_variation", "suffix_impersonation"}
		if threatAnalysisOpts.BrandCombo {
			techniques = append(techniques, "subdomain_pattern")
		}
		return techniques
	}

	if threatAnalysisOpts.Typosquatting {
		techniques = append(techniques, "typosquatting")
	}
	if threatAnalysisOpts.Bitsquatting {
		techniques = append(techniques, "bitsquatting")
	}
	if threatAnalysisOpts.Homographic {
		techniques = append(techniques, "homographic")
	}

	// Se nenhuma foi selecionada especificamente, habilitar conjunto padrão (sem brand-combo)
	if len(techniques) == 0 {
		techniques = []string{"typosquatting", "bitsquatting", "homographic", "insertion", "deletion", "transposition", "tld_variation", "suffix_impersonation"}
		if threatAnalysisOpts.BrandCombo {
			techniques = append(techniques, "subdomain_pattern")
		}
		return techniques
	}

	return techniques
}

func init() {
	rootCmd.AddCommand(threatAnalysisCmd)

	// Flags principais
	threatAnalysisCmd.Flags().StringVarP(&opts.DnsSuffix, "domain", "d", "", "Target domain for analysis")
	threatAnalysisCmd.Flags().StringVarP(&fileOptions.DnsSuffixFile, "domain-list", "L", "", "File with domains to analyze (one per line)")

	// Técnicas específicas
	threatAnalysisCmd.Flags().BoolVar(&threatAnalysisOpts.Typosquatting, "typosquatting", false, "Enable typosquatting detection (keyboard adjacency errors)")
	threatAnalysisCmd.Flags().BoolVar(&threatAnalysisOpts.Bitsquatting, "bitsquatting", false, "Enable bitsquatting detection (single bit-flip errors)")
	threatAnalysisCmd.Flags().BoolVar(&threatAnalysisOpts.Homographic, "homographic", false, "Enable homographic attacks detection (similar character substitution)")
	threatAnalysisCmd.Flags().BoolVar(&threatAnalysisOpts.AllTechniques, "all-techniques", false, "Enable all available detection techniques")

	// Configurações avançadas
	threatAnalysisCmd.Flags().IntVar(&threatAnalysisOpts.MaxVariations, "max-variations", 1000, "Maximum variations to generate per domain (10-50000)")
	threatAnalysisCmd.Flags().StringSliceVar(&threatAnalysisOpts.TLDs, "target-tlds", []string{"com", "net", "org", "co", "io", "tk", "ml", "ga", "cf", "com.br", "net.br", "org.br"}, "Target TLDs for variations (comma-separated)")
	threatAnalysisCmd.Flags().BoolVar(&threatAnalysisOpts.WarnLimits, "warn-limits", true, "Show warnings when variation limits are reached")
	threatAnalysisCmd.Flags().BoolVar(&threatAnalysisOpts.BrandCombo, "brand-combo", false, "Enable brand-combo patterns (prefix/suffix combos) in addition to eTLD+1 variations")
	threatAnalysisCmd.Flags().BoolVar(&threatAnalysisOpts.SpanLast3, "span-last3", false, "Operate over the last 3 labels: mutate the 3rd-from-right label and keep the last 2 as suffix")
	threatAnalysisCmd.Flags().BoolVar(&threatAnalysisOpts.EmitCandidates, "emit-candidates", false, "Write all generated candidates to outputs before probing (includes NX domains)")
	threatAnalysisCmd.Flags().StringVar(&threatAnalysisOpts.FocusSuffix, "focus-suffix", "", "Emphasize suffix-specific techniques (e.g., gov.br)")

	// Flags de compatibilidade
	threatAnalysisCmd.Flags().BoolVarP(&fileOptions.IgnoreNonexistent, "ignore-nonexistent", "I", false, "Ignore nonexistent domains during analysis")
	threatAnalysisCmd.Flags().BoolVarP(&opts.Quick, "quick", "Q", false, "Request only A records (faster)")

	// Marcar flags mutuamente exclusivas
	threatAnalysisCmd.MarkFlagsMutuallyExclusive("domain", "domain-list")

	// Adicionar validações
	threatAnalysisCmd.Flags().SetInterspersed(false) // Manter ordem das flags
}

// computeAllowedTLDs decides which TLDs should be used for variation generation
// respecting the Public Suffix List. If the suffix has multiple labels (e.g., gov.br),
// we restrict to that exact suffix. Otherwise we allow the provided target list.
func computeAllowedTLDs(domain string, targets []string) []string {
	base := strings.TrimSuffix(strings.ToLower(domain), ".")
	if base == "" {
		return targets
	}
	suffix, _ := publicsuffix.PublicSuffix(base)
	// Build union: keep original suffix and user-provided targets
	set := map[string]struct{}{}
	if suffix != "" {
		set[suffix] = struct{}{}
	}
	for _, t := range targets {
		tt := strings.Trim(strings.ToLower(t), ". ")
		if tt != "" {
			set[tt] = struct{}{}
		}
	}
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	return out
}
