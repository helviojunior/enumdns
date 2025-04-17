package cmd

import (
    "errors"
    "log/slog"
    "os"
    "fmt"

    "github.com/helviojunior/enumdns/internal/ascii"
    "github.com/helviojunior/enumdns/internal/tools"
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
   - enumdns brute -d helviojunior.com.br -w /tmp/wordlist.txt -o enumdns.txt
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
        if opts.DnsSuffix == "" && fileOptions.DnsSuffixFile == "" {
            return errors.New("a DNS suffix or DNS suffix file must be specified")
        }

        if fileOptions.DnsSuffixFile != "" {
            if !tools.FileExists(fileOptions.DnsSuffixFile) {
                return errors.New("DNS suffix file is not readable")
            }
        }

        if fileOptions.HostFile == "" {
            return errors.New("a wordlist file must be specified")
        }

        if !tools.FileExists(fileOptions.HostFile) {
            return errors.New("wordlist file is not readable")
        }

        return nil
    },
    Run: func(cmd *cobra.Command, args []string) {

        //Check DNS connectivity
        _, err := tools.GetValidDnsSuffix(fileOptions.DnsServer, "google.com.", opts.Proxy)
        if err != nil {
            log.Error("Error checking DNS connectivity", "err", err)
            os.Exit(2)
        }

        log.Debug("starting DNS brute-force")

        dnsSuffix := []string{}
        hostWordList := []string{}
        reader := readers.NewFileReader(fileOptions)
        total := 0

        if fileOptions.DnsSuffixFile != "" {
            log.Debugf("Reading dns suffix file: %s", fileOptions.DnsSuffixFile)
            if err := reader.ReadDnsList(&dnsSuffix); err != nil {
                log.Error("error in reader.Read", "err", err)
                log.Warn("If you are facing error related to 'SOA not found for domain' you can ignore it with -I option")
                os.Exit(2)
            }
        }else{
            //Check if DNS exists
            s, err := tools.GetValidDnsSuffix(fileOptions.DnsServer, opts.DnsSuffix, opts.Proxy)
            if err != nil {
                log.Error("invalid dns suffix", "suffix", opts.DnsSuffix, "err", err)
                os.Exit(2)
            }
            dnsSuffix = append(dnsSuffix, s)
        }
        log.Debugf("Loaded %s DNS suffix(es)", tools.FormatInt(len(dnsSuffix)))

        log.Debugf("Reading dns word list file: %s", fileOptions.HostFile)
        if err := reader.ReadWordList(&hostWordList); err != nil {
            log.Error("error in reader.Read", "err", err)
            os.Exit(2)
        }
        total = len(dnsSuffix) * len(hostWordList)

        if len(dnsSuffix) == 0 {
            log.Error("DNS suffix list is empty")
            os.Exit(2)
        }

        log.Infof("Enumerating %s DNS hosts", tools.FormatInt(total))

        // Check runned items
        conn, _ := database.Connection("sqlite:///" + opts.Writer.UserPath +"/.enumdns.db", true, false)

        go func() {
            defer close(bruteRunner.Targets)

            ascii.HideCursor()
            for _, s := range dnsSuffix {
                bruteRunner.Targets <- s
                for _, h := range hostWordList {

                    i := true
                    host := h + "." + s
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

                    if i || forceCheck{
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
    
    bruteCmd.Flags().StringVarP(&opts.DnsSuffix, "dns-suffix", "d", "", "Single DNS suffix. (ex: helviojunior.com.br)")
    bruteCmd.Flags().StringVarP(&fileOptions.DnsSuffixFile, "dns-list", "L", "", "File containing a list of DNS suffix")
    bruteCmd.Flags().StringVarP(&fileOptions.HostFile, "word-list", "w", "", "File containing a list of DNS hosts")
    
    bruteCmd.Flags().BoolVarP(&fileOptions.IgnoreNonexistent, "IgnoreNonexistent", "I", false, "Ignore Nonexistent DNS suffix. Used only with --dns-list option.")

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
    bruteCmd.Flags().StringVar(&opts.Writer.ELasticURI, "write-elasticsearch-uri", "http://localhost:9200/enumdns", "The elastic search URI to use. (e.g., http://user:pass@host:9200/index)")

}