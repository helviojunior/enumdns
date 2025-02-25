package cmd

import (
    "errors"
    "log/slog"
    "os"
    "fmt"

    "github.com/helviojunior/enumdns/internal/ascii"
    "github.com/helviojunior/enumdns/internal/islazy"
    "github.com/helviojunior/enumdns/pkg/log"
    "github.com/helviojunior/enumdns/pkg/runner"
    "github.com/helviojunior/enumdns/pkg/database"
    "github.com/helviojunior/enumdns/pkg/writers"
    "github.com/helviojunior/enumdns/pkg/readers"
    "github.com/spf13/cobra"
)

var bruteRunner *runner.Runner

var bruteWriters = []writers.Writer{}
var bruteCmd = &cobra.Command{
    Use:   "brute",
    Short: "Perform brute-force enumeration",
    Long: ascii.LogoHelp(ascii.Markdown(`
# brute

Perform brute-force enumeration.

By default, enumdns will only show information regarding the brute-force process. 
However, that is only half the fun! You can add multiple _writers_ that will 
collect information such as response codes, content, and more. You can specify 
multiple writers using the _--writer-*_ flags (see --help).
`)),
    Example: `
   - enumdns brute -d helviojunior.com.br -w /tmp/wordlist.txt --write-jsonl
   - enumdns brute -D domains.txt -w /tmp/wordlist.txt --write-db`,
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        var err error

        // Annoying quirk, but because I'm overriding PersistentPreRun
        // here which overrides the parent it seems.
        // So we need to explicitly call the parent's one now.
        if err = rootCmd.PersistentPreRunE(cmd, args); err != nil {
            return err
        }

        // An slog-capable logger to use with drivers and runners
        logger := slog.New(log.Logger)

        // Configure writers that subcommand scanners will pass to
        // a runner instance.

        //The first one is the general writer (global user)
        w, err := writers.NewDbWriter("sqlite:///" + opts.Writer.UserPath +"/.enumdns.db", false)
        if err != nil {
            return err
        }
        bruteWriters = append(bruteWriters, w)

        //The second one is the STDOut
        if opts.Logging.Silence != true {
            w, err := writers.NewStdoutWriter()
            if err != nil {
                return err
            }
            bruteWriters = append(bruteWriters, w)
        }
    
        if opts.Writer.Text {
            w, err := writers.NewTextWriter(opts.Writer.TextFile)
            if err != nil {
                return err
            }
            bruteWriters = append(bruteWriters, w)
        }

        if opts.Writer.Jsonl {
            w, err := writers.NewJsonWriter(opts.Writer.JsonlFile)
            if err != nil {
                return err
            }
            bruteWriters = append(bruteWriters, w)
        }

        if opts.Writer.Db {
            w, err := writers.NewDbWriter(opts.Writer.DbURI, opts.Writer.DbDebug)
            if err != nil {
                return err
            }
            bruteWriters = append(bruteWriters, w)
        }

        if opts.Writer.Csv {
            w, err := writers.NewCsvWriter(opts.Writer.CsvFile)
            if err != nil {
                return err
            }
            bruteWriters = append(bruteWriters, w)
        }

        if opts.Writer.ELastic {
            w, err := writers.NewElasticWriter(opts.Writer.ELasticURI)
            if err != nil {
                return err
            }
            bruteWriters = append(bruteWriters, w)
        }

        if opts.Writer.None {
            w, err := writers.NewNoneWriter()
            if err != nil {
                return err
            }
            bruteWriters = append(bruteWriters, w)
        }

        if len(bruteWriters) == 0 {
            log.Warn("no writers have been configured. to persist probe results, add writers using --write-* flags")
        }

        // Get the runner up. Basically, all of the subcommands will use this.
        bruteRunner, err = runner.NewRunner(logger, *opts, bruteWriters)
        if err != nil {
            return err
        }

        fileOptions.DnsServer = opts.DnsServer + ":" + fmt.Sprintf("%d", opts.DnsPort)

        return nil
    },
    PreRunE: func(cmd *cobra.Command, args []string) error {
        if opts.DnsSufix == "" && fileOptions.DnsSufixFile == "" {
            return errors.New("a DNS sufix or DNS sufix file must be specified")
        }

        if fileOptions.DnsSufixFile != "" {
            if !islazy.FileExists(fileOptions.DnsSufixFile) {
                return errors.New("DNS sufix file is not readable")
            }
        }

        if fileOptions.HostFile == "" {
            return errors.New("a wordlist file must be specified")
        }

        if !islazy.FileExists(fileOptions.HostFile) {
            return errors.New("wordlist file is not readable")
        }

        return nil
    },
    Run: func(cmd *cobra.Command, args []string) {

        //Check DNS connectivity
        _, err := islazy.GetValidDnsSufix(fileOptions.DnsServer, "google.com.")
        if err != nil {
            log.Error("Error checking DNS connectivity", "err", err)
            os.Exit(2)
        }

        log.Debug("starting DNS brute-force")

        dnsSufix := []string{}
        hostWordList := []string{}
        reader := readers.NewFileReader(fileOptions)
        total := 0

        if fileOptions.DnsSufixFile != "" {
            log.Debugf("Reading dns sufix file: %s", fileOptions.DnsSufixFile)
            if err := reader.ReadDnsList(&dnsSufix); err != nil {
                log.Error("error in reader.Read", "err", err)
                os.Exit(2)
            }
        }else{
            //Check if DNS exists
            s, err := islazy.GetValidDnsSufix(fileOptions.DnsServer, opts.DnsSufix)
            if err != nil {
                log.Error("invalid dns sufix", "sufix", opts.DnsSufix, "err", err)
                os.Exit(2)
            }
            dnsSufix = append(dnsSufix, s)
        }
        log.Debugf("Loaded %s DNS sufix(es)", islazy.FormatInt(len(dnsSufix)))

        log.Debugf("Reading dns word list file: %s", fileOptions.HostFile)
        if err := reader.ReadWordList(&hostWordList); err != nil {
            log.Error("error in reader.Read", "err", err)
            os.Exit(2)
        }
        total = len(dnsSufix) * len(hostWordList)

        log.Infof("Enumerating %s DNS hosts", islazy.FormatInt(total))

        // Check runned items
        conn, _ := database.Connection("sqlite:///" + opts.Writer.UserPath +"/.enumdns.db", true, false)

        go func() {
            defer close(bruteRunner.Targets)
            for _, s := range dnsSufix {
                for _, h := range hostWordList {

                    i := true
                    host := h + "." + s
                    response := conn.Raw("SELECT count(id) as count from results WHERE failed = 0 AND fqdn = ?", host)
                    if response != nil {
                        var cnt int
                        _ = response.Row().Scan(&cnt)
                        i = (cnt == 0)
                        if cnt > 0 {
                            log.Debug("[Host already checked]", "fqdn", host)
                        }
                    }

                    if i {
                        bruteRunner.Targets <- host
                    }else{
                        bruteRunner.AddSkiped()
                    }
                }
            }
        
        }()

        bruteRunner.Run(total)
        bruteRunner.Close()

    },
}

func init() {
    rootCmd.AddCommand(bruteCmd)
    
    bruteCmd.Flags().StringVarP(&opts.DnsSufix, "dns-sufix", "d", "", "Single DNS sufix. (ex: helviojunior.com.br)")
    bruteCmd.Flags().StringVarP(&fileOptions.DnsSufixFile, "dns-list", "L", "", "File containing a list of DNS sufix")
    bruteCmd.Flags().StringVarP(&fileOptions.HostFile, "word-list", "w", "", "File containing a list of DNS hosts")
    
    bruteCmd.Flags().StringVarP(&opts.DnsServer, "server", "s", "8.8.8.8", "DNS Server")
    bruteCmd.Flags().IntVar(&opts.DnsPort, "port", 53, "DNS Server Port")
    bruteCmd.Flags().StringVarP(&opts.DnsProtocol, "protocol", "", "UDP", "DNS Server protocol (TCP/UDP)")
    
    // Logging control for subcommands
    bruteCmd.Flags().BoolVar(&opts.Logging.LogScanErrors, "log-scan-errors", false, "Log scan errors (timeouts, DNS errors, etc.) to stderr (warning: can be verbose!)")

    // "Threads" & other
    bruteCmd.Flags().IntVarP(&opts.Scan.Threads, "threads", "t", 16, "Number of concurrent threads (goroutines) to use")
    bruteCmd.Flags().IntVarP(&opts.Scan.Timeout, "timeout", "T", 60, "Number of seconds before considering a page timed out")
    bruteCmd.Flags().IntVar(&opts.Scan.Delay, "delay", 3, "Number of seconds delay between navigation and screenshotting")

    // Write options for scan subcommands
    bruteCmd.Flags().BoolVar(&opts.Writer.Db, "write-db", false, "Write results to a SQLite database")
    bruteCmd.Flags().StringVar(&opts.Writer.DbURI, "write-db-uri", "sqlite://enumdns.sqlite3", "The database URI to use. Supports SQLite, Postgres, and MySQL (e.g., postgres://user:pass@host:port/db)")
    bruteCmd.Flags().BoolVar(&opts.Writer.DbDebug, "write-db-enable-debug", false, "Enable database query debug logging (warning: verbose!)")
    bruteCmd.Flags().BoolVar(&opts.Writer.Csv, "write-csv", false, "Write results as CSV (has limited columns)")
    bruteCmd.Flags().StringVar(&opts.Writer.CsvFile, "write-csv-file", "enumdns.csv", "The file to write CSV rows to")
    bruteCmd.Flags().BoolVar(&opts.Writer.Jsonl, "write-jsonl", false, "Write results as JSON lines")
    bruteCmd.Flags().StringVar(&opts.Writer.JsonlFile, "write-jsonl-file", "enumdns.jsonl", "The file to write JSON lines to")
    bruteCmd.Flags().BoolVar(&opts.Writer.None, "write-none", false, "Use an empty writer to silence warnings")

    bruteCmd.Flags().BoolVar(&opts.Writer.ELastic, "write-elastic", false, "Write results to a SQLite database")
    bruteCmd.Flags().StringVar(&opts.Writer.ELasticURI, "write-elasticsearch-uri", "http://localhost:9200/intelparser", "The elastic search URI to use. (e.g., http://user:pass@host:9200/index)")

}