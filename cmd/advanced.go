package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/bob-reis/enumdns/internal/ascii"
	"github.com/bob-reis/enumdns/internal/tools"
	"github.com/bob-reis/enumdns/pkg/advanced"
	"github.com/bob-reis/enumdns/pkg/database"
	"github.com/bob-reis/enumdns/pkg/log"
	"github.com/bob-reis/enumdns/pkg/readers"
	"github.com/bob-reis/enumdns/pkg/runner"
	"github.com/bob-reis/enumdns/pkg/writers"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// Opções específicas do módulo threat-analysis
var threatAnalysisOpts = struct {
	Typosquatting bool
	Bitsquatting  bool
	Homographic   bool
	AllTechniques bool
	MaxVariations int
	TLDs          []string
	ShowProgress  bool
	WarnLimits    bool
}{
	MaxVariations: 1000,
	TLDs:          []string{"com", "net", "org", "co", "io"},
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
   - enumdns threat-analysis -L domains.txt --bitsquatting --write-jsonl`,
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

		// Configurar writers (mesmo padrão dos outros comandos)
		// Writer de controle (obrigatório)
		w, err := writers.NewDbWriter(opts.Writer.CtrlDbURI, opts.Writer.DbDebug)
		if err != nil {
			return err
		}
		advancedWriters = append(advancedWriters, w)

		// Writer de stdout
		if !opts.Logging.Silence {
			w, err := writers.NewStdoutWriter()
			if err != nil {
				return err
			}
			advancedWriters = append(advancedWriters, w)
		}

		// Writers opcionais
		if opts.Writer.Text {
			w, err := writers.NewTextWriter(opts.Writer.TextFile)
			if err != nil {
				return err
			}
			advancedWriters = append(advancedWriters, w)
		}

		if opts.Writer.Jsonl {
			w, err := writers.NewJsonWriter(opts.Writer.JsonlFile)
			if err != nil {
				return err
			}
			advancedWriters = append(advancedWriters, w)
		}

		if opts.Writer.Db {
			w, err := writers.NewDbWriter(opts.Writer.DbURI, opts.Writer.DbDebug)
			if err != nil {
				return err
			}
			advancedWriters = append(advancedWriters, w)
		}

		if opts.Writer.Csv {
			w, err := writers.NewCsvWriter(opts.Writer.CsvFile)
			if err != nil {
				return err
			}
			advancedWriters = append(advancedWriters, w)
		}

		if opts.Writer.ELastic {
			w, err := writers.NewElasticWriter(opts.Writer.ELasticURI)
			if err != nil {
				return err
			}
			advancedWriters = append(advancedWriters, w)
		}

		if opts.Writer.None {
			w, err := writers.NewNoneWriter()
			if err != nil {
				return err
			}
			advancedWriters = append(advancedWriters, w)
		}

		if len(advancedWriters) == 0 {
			log.Warn("no writers have been configured. to persist probe results, add writers using --write-* flags")
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

		generator := advanced.NewVariationGenerator(domain, advanced.GeneratorOptions{
			Techniques:    techniques,
			MaxVariations: threatAnalysisOpts.MaxVariations,
			TargetTLDs:    threatAnalysisOpts.TLDs,
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

	executeAnalysis(allVariations)
}

func getEnabledTechniques() []string {
	var techniques []string

	if threatAnalysisOpts.AllTechniques {
		return []string{"typosquatting", "bitsquatting", "homographic", "insertion", "deletion", "transposition", "tld_variation", "subdomain_pattern"}
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

	// Se nenhuma foi selecionada especificamente, habilitar todas
	if len(techniques) == 0 {
		return []string{"typosquatting", "bitsquatting", "homographic", "insertion", "deletion", "transposition", "tld_variation", "subdomain_pattern"}
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
	threatAnalysisCmd.Flags().StringSliceVar(&threatAnalysisOpts.TLDs, "target-tlds", []string{"com", "net", "org", "co", "io", "tk", "ml", "ga", "cf"}, "Target TLDs for variations (comma-separated)")
	threatAnalysisCmd.Flags().BoolVar(&threatAnalysisOpts.WarnLimits, "warn-limits", true, "Show warnings when variation limits are reached")

	// Flags de compatibilidade
	threatAnalysisCmd.Flags().BoolVarP(&fileOptions.IgnoreNonexistent, "ignore-nonexistent", "I", false, "Ignore nonexistent domains during analysis")
	threatAnalysisCmd.Flags().BoolVarP(&opts.Quick, "quick", "Q", false, "Request only A records (faster)")

	// Marcar flags mutuamente exclusivas
	threatAnalysisCmd.MarkFlagsMutuallyExclusive("domain", "domain-list")

	// Adicionar validações
	threatAnalysisCmd.Flags().SetInterspersed(false) // Manter ordem das flags
}
