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

var resolveRunner *runner.Runner

var resolveWriters = []writers.Writer{}
var resolveCmd = &cobra.Command{
    Use:   "resolve",
    Short: "Perform resolve roperations",
    Long: ascii.LogoHelp(ascii.Markdown(`
# resolve

Perform resolver operations.
`)),
    Example: `
   - enumdns resolve bloodhound -L /tmp/bloodhound_computers.json -o enumdns.txt
   - enumdns resolve bloodhound -L /tmp/bloodhound_files.zip --write-jsonl
   - enumdns resolve bloodhound -L /tmp/bloodhound_computers.json --write-db

   - enumdns resolve file -L /tmp/host_list.txt -o enumdns.txt
   - enumdns resolve file -L /tmp/host_list.txt --write-jsonl
   - enumdns resolve file -L /tmp/host_list.txt --write-db`,
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
        resolveWriters = append(resolveWriters, w)

        //The second one is the STDOut
        if opts.Logging.Silence != true {
            w, err := writers.NewStdoutWriter()
            if err != nil {
                return err
            }
            resolveWriters = append(resolveWriters, w)
        }
    
        if opts.Writer.Text {
            w, err := writers.NewTextWriter(opts.Writer.TextFile)
            if err != nil {
                return err
            }
            resolveWriters = append(resolveWriters, w)
        }

        if opts.Writer.Jsonl {
            w, err := writers.NewJsonWriter(opts.Writer.JsonlFile)
            if err != nil {
                return err
            }
            resolveWriters = append(resolveWriters, w)
        }

        if opts.Writer.Db {
            w, err := writers.NewDbWriter(opts.Writer.DbURI, opts.Writer.DbDebug)
            if err != nil {
                return err
            }
            resolveWriters = append(resolveWriters, w)
        }

        if opts.Writer.Csv {
            w, err := writers.NewCsvWriter(opts.Writer.CsvFile)
            if err != nil {
                return err
            }
            resolveWriters = append(resolveWriters, w)
        }

        if opts.Writer.ELastic {
            w, err := writers.NewElasticWriter(opts.Writer.ELasticURI)
            if err != nil {
                return err
            }
            resolveWriters = append(resolveWriters, w)
        }

        if opts.Writer.None {
            w, err := writers.NewNoneWriter()
            if err != nil {
                return err
            }
            resolveWriters = append(resolveWriters, w)
        }

        if len(resolveWriters) == 0 {
            log.Warn("no writers have been configured. to persist probe results, add writers using --write-* flags")
        }

        // Get the runner up. Basically, all of the subcommands will use this.
        resolveRunner, err = runner.NewRunner(logger, *opts, resolveWriters)
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

        //Check DNS connectivity
        _, err := tools.GetValidDnsSuffix(fileOptions.DnsServer, "google.com.", opts.Proxy)
        if err != nil {
            return errors.New("Error checking DNS connectivity: " + err.Error())
        }

        return nil
    },
    Run: func(cmd *cobra.Command, args []string) {

        log.Debug("starting DNS resolve-force")

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
            defer close(resolveRunner.Targets)

            ascii.HideCursor()
            for _, s := range dnsSuffix {
                resolveRunner.Targets <- s
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
                        resolveRunner.Targets <- host
                    }else{
                        resolveRunner.AddSkiped()
                    }
                }
            }
        
        }()

        resolveRunner.Run(total)
        resolveRunner.Close()

    },
}

func init() {
    rootCmd.AddCommand(resolveCmd)

    resolveCmd.PersistentFlags().BoolVarP(&fileOptions.IgnoreNonexistent, "IgnoreNonexistent", "I", false, "Ignore Nonexistent DNS suffix. Used only with --dns-list option.")

    resolveCmd.PersistentFlags().StringVarP(&opts.DnsServer, "server", "s", "8.8.8.8", "DNS Server")
    resolveCmd.PersistentFlags().IntVar(&opts.DnsPort, "port", 53, "DNS Server Port")
    resolveCmd.PersistentFlags().StringVarP(&opts.DnsProtocol, "protocol", "", "UDP", "DNS Server protocol (TCP/UDP)")
    
    // Logging control for subcommands
    resolveCmd.PersistentFlags().BoolVar(&opts.Logging.LogScanErrors, "log-scan-errors", false, "Log scan errors (timeouts, DNS errors, etc.) to stderr (warning: can be verbose!)")

    // "Threads" & other
    resolveCmd.PersistentFlags().IntVarP(&opts.Scan.Threads, "threads", "t", 16, "Number of concurrent threads (goroutines) to use")
    resolveCmd.PersistentFlags().IntVarP(&opts.Scan.Timeout, "timeout", "T", 60, "Number of seconds before considering a page timed out")
    resolveCmd.PersistentFlags().IntVar(&opts.Scan.Delay, "delay", 3, "Number of seconds delay between navigation and screenshotting")

    // Write options for scan subcommands
    resolveCmd.PersistentFlags().BoolVar(&opts.Writer.Db, "write-db", false, "Write results to a SQLite database")
    resolveCmd.PersistentFlags().StringVar(&opts.Writer.DbURI, "write-db-uri", "sqlite://enumdns.sqlite3", "The database URI to use. Supports SQLite, Postgres, and MySQL (e.g., postgres://user:pass@host:port/db)")
    resolveCmd.PersistentFlags().BoolVar(&opts.Writer.DbDebug, "write-db-enable-debug", false, "Enable database query debug logging (warning: verbose!)")
    resolveCmd.PersistentFlags().BoolVar(&opts.Writer.Csv, "write-csv", false, "Write results as CSV (has limited columns)")
    resolveCmd.PersistentFlags().StringVar(&opts.Writer.CsvFile, "write-csv-file", "enumdns.csv", "The file to write CSV rows to")
    resolveCmd.PersistentFlags().BoolVar(&opts.Writer.Jsonl, "write-jsonl", false, "Write results as JSON lines")
    resolveCmd.PersistentFlags().StringVar(&opts.Writer.JsonlFile, "write-jsonl-file", "enumdns.jsonl", "The file to write JSON lines to")
    resolveCmd.PersistentFlags().BoolVar(&opts.Writer.None, "write-none", false, "Use an empty writer to silence warnings")

    resolveCmd.PersistentFlags().BoolVar(&opts.Writer.ELastic, "write-elastic", false, "Write results to a SQLite database")
    resolveCmd.PersistentFlags().StringVar(&opts.Writer.ELasticURI, "write-elasticsearch-uri", "http://localhost:9200/enumdns", "The elastic search URI to use. (e.g., http://user:pass@host:9200/index)")

}