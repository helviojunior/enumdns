package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/helviojunior/enumdns/internal/ascii"
	"github.com/helviojunior/enumdns/internal/tools"
	"github.com/helviojunior/enumdns/pkg/advanced"
	"github.com/helviojunior/enumdns/pkg/database"
	"github.com/helviojunior/enumdns/pkg/log"
	"github.com/helviojunior/enumdns/pkg/readers"
	"github.com/helviojunior/enumdns/pkg/runner"
	"github.com/helviojunior/enumdns/pkg/writers"
	"github.com/spf13/cobra"
)

// Opções específicas do módulo advanced
var advancedOpts = struct {
	Typosquatting bool
	Bitsquatting  bool
	Homographic   bool
	AllTechniques bool
	MaxVariations int
	TLDs          []string
}{
	MaxVariations: 1000,
	TLDs:          []string{"com", "net", "org", "co", "io"},
}

var advancedWriters = []writers.Writer{}

var advancedCmd = &cobra.Command{
	Use:   "advanced",
	Short: "Advanced domain analysis for threats and typosquatting",
	Long: ascii.LogoHelp(ascii.Markdown(`
# advanced

Perform advanced domain analysis including typosquatting detection,
homographic attacks, and suspicious domain discovery.

This module generates domain variations using multiple techniques and 
verifies their existence to identify potential threats.
`)),
	Example: `
   - enumdns advanced -d example.com --typosquatting -o threats.txt
   - enumdns advanced -d example.com --all-techniques --write-db
   - enumdns advanced -L domains.txt --bitsquatting --write-jsonl`,
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
		if opts.Logging.Silence != true {
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
		if !advancedOpts.Typosquatting && !advancedOpts.Bitsquatting &&
			!advancedOpts.Homographic && !advancedOpts.AllTechniques {
			log.Warn("No specific techniques selected, enabling all techniques")
			advancedOpts.AllTechniques = true
		}

		return nil
	},
	Run: runAdvancedAnalysis,
}

func runAdvancedAnalysis(cmd *cobra.Command, args []string) {
	// Verificar conectividade DNS
	_, err := tools.GetValidDnsSuffix(fileOptions.DnsServer, "google.com.", opts.Proxy)
	if err != nil {
		log.Error("Error checking DNS connectivity", "err", err)
		os.Exit(2)
	}

	log.Debug("starting advanced domain analysis")

	domains := []string{}

	// Determinar domínios alvo
	if fileOptions.DnsSuffixFile != "" {
		log.Debugf("Reading domain list file: %s", fileOptions.DnsSuffixFile)
		reader := readers.NewFileReader(fileOptions)
		if err := reader.ReadDnsList(&domains); err != nil {
			log.Error("error reading domain list", "err", err)
			os.Exit(2)
		}
	} else {
		// Validar domínio único
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

	log.Infof("Analyzing %d domain(s) for threats", len(domains))

	// Configurar técnicas habilitadas
	techniques := getEnabledTechniques()
	log.Infof("Enabled techniques: %v", techniques)

	allVariations := []advanced.Variation{}

	// Gerar variações para cada domínio
	for _, domain := range domains {
		log.Infof("Generating variations for %s", domain)

		generator := advanced.NewVariationGenerator(domain, advanced.GeneratorOptions{
			Techniques:    techniques,
			MaxVariations: advancedOpts.MaxVariations,
			TargetTLDs:    advancedOpts.TLDs,
		})

		variations := generator.GenerateAll()
		log.Infof("Generated %d variations for %s", len(variations), domain)

		allVariations = append(allVariations, variations...)
	}

	total := len(allVariations)
	log.Infof("Total variations to check: %s", tools.FormatInt(total))

	if total == 0 {
		log.Warn("No variations generated")
		return
	}

	// Configurar runner
	runnerLogger := slog.New(log.Logger)
	runnerInstance, err := runner.NewRunner(runnerLogger, *opts, advancedWriters)
	if err != nil {
		log.Error("Failed to create runner", "err", err)
		os.Exit(2)
	}

	// Verificar itens já analisados (controle de duplicatas)
	conn, err := database.Connection(opts.Writer.CtrlDbURI, true, false)
	if err != nil {
		log.Error("Error establishing connection with Database", "err", err)
		os.Exit(2)
	}

	// Alimentar o canal com as variações
	go func() {
		defer close(runnerInstance.Targets)

		ascii.HideCursor()
		for _, variation := range allVariations {
			i := true
			host := variation.Domain

			if !forceCheck {
				response := conn.Raw("SELECT count(id) as count from results WHERE failed = 0 AND fqdn = ?", host)
				if response != nil {
					var cnt int
					_ = response.Row().Scan(&cnt)
					i = (cnt == 0)
					if cnt > 0 {
						log.Debug("[Host already checked]", "fqdn", host)
					}
				}
			}

			if i || forceCheck {
				runnerInstance.Targets <- host
			} else {
				runnerInstance.AddSkiped()
			}
		}
	}()

	// Executar análise
	status := runnerInstance.Run(total)
	runnerInstance.Close()

	if status.Skiped > 0 {
		log.Warnf("%d variations were skipped because they were already analyzed. Use the --force parameter to reanalyze them.", status.Skiped)
		ascii.ClearLine()
	}

	log.Info("Advanced analysis completed")
}

func getEnabledTechniques() []string {
	var techniques []string

	if advancedOpts.AllTechniques {
		return []string{"typosquatting", "bitsquatting", "homographic", "insertion", "deletion", "transposition", "tld_variation", "subdomain_pattern"}
	}

	if advancedOpts.Typosquatting {
		techniques = append(techniques, "typosquatting")
	}
	if advancedOpts.Bitsquatting {
		techniques = append(techniques, "bitsquatting")
	}
	if advancedOpts.Homographic {
		techniques = append(techniques, "homographic")
	}

	// Se nenhuma foi selecionada especificamente, habilitar todas
	if len(techniques) == 0 {
		return []string{"typosquatting", "bitsquatting", "homographic", "insertion", "deletion", "transposition", "tld_variation", "subdomain_pattern"}
	}

	return techniques
}

func init() {
	rootCmd.AddCommand(advancedCmd)

	// Flags principais
	advancedCmd.Flags().StringVarP(&opts.DnsSuffix, "domain", "d", "", "Target domain for analysis")
	advancedCmd.Flags().StringVarP(&fileOptions.DnsSuffixFile, "domain-list", "L", "", "File with domains to analyze (one per line)")

	// Técnicas específicas
	advancedCmd.Flags().BoolVar(&advancedOpts.Typosquatting, "typosquatting", false, "Enable typosquatting detection (keyboard adjacency errors)")
	advancedCmd.Flags().BoolVar(&advancedOpts.Bitsquatting, "bitsquatting", false, "Enable bitsquatting detection (single bit-flip errors)")
	advancedCmd.Flags().BoolVar(&advancedOpts.Homographic, "homographic", false, "Enable homographic attacks detection (similar character substitution)")
	advancedCmd.Flags().BoolVar(&advancedOpts.AllTechniques, "all-techniques", false, "Enable all available detection techniques")

	// Configurações avançadas
	advancedCmd.Flags().IntVar(&advancedOpts.MaxVariations, "max-variations", 1000, "Maximum variations to generate per domain (1-10000)")
	advancedCmd.Flags().StringSliceVar(&advancedOpts.TLDs, "target-tlds", []string{"com", "net", "org", "co", "io", "tk", "ml", "ga", "cf"}, "Target TLDs for variations (comma-separated)")

	// Flags de compatibilidade
	advancedCmd.Flags().BoolVarP(&fileOptions.IgnoreNonexistent, "ignore-nonexistent", "I", false, "Ignore nonexistent domains during analysis")
	advancedCmd.Flags().BoolVarP(&opts.Quick, "quick", "Q", false, "Request only A records (faster)")

	// Marcar flags mutuamente exclusivas
	advancedCmd.MarkFlagsMutuallyExclusive("domain", "domain-list")

	// Adicionar validações
	advancedCmd.Flags().SetInterspersed(false) // Manter ordem das flags
}
