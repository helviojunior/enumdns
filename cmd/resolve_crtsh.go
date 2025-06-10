package cmd

import (
    "errors"
    "log/slog"
    "os"
    "time"

    "github.com/helviojunior/enumdns/internal/ascii"
    "github.com/helviojunior/enumdns/internal/tools"
    "github.com/helviojunior/enumdns/pkg/log"
    "github.com/helviojunior/enumdns/pkg/runner"
    "github.com/helviojunior/enumdns/pkg/database"
    "github.com/helviojunior/enumdns/pkg/writers"
    "github.com/helviojunior/enumdns/pkg/readers"
    resolver "github.com/helviojunior/gopathresolver"
    "github.com/spf13/cobra"
)

var fqdnOutFile = ""
var resolveCrtshWriters = []writers.Writer{}
var resolveCrtshCmd = &cobra.Command{
    Use:   "crtsh",
    Short: "Perform resolve roperations",
    Long: ascii.LogoHelp(ascii.Markdown(`
# resolve crtsh

Perform cert.sh crawler + resolve enumeration.

By default, enumdns will only show information regarding the resolver process. 
However, that is only half the fun! You can add multiple _writers_ that will 
collect information such as response codes, content, and more. You can specify 
multiple writers using the _--writer-*_ flags (see --help).
`)),
    Example: `
   - enumdns resolve crtsh -d helviojunior.com.br -o enumdns.txt
   - enumdns resolve crtsh -L domains.txt --write-jsonl
   - enumdns resolve crtsh -L domains.txt --write-db`,
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        var err error
        // Annoying quirk, but because I'm overriding PersistentPreRun
        // here which overrides the parent it seems.
        // So we need to explicitly call the parent's one now.
        if err := resolveCmd.PersistentPreRunE(cmd, args); err != nil {
            return err
        }

        // An slog-capable logger to use with drivers and runners
        logger := slog.New(log.Logger)

        // Get the runner up. Basically, all of the subcommands will use this.
        bruteRunner, err = runner.NewRunner(logger, *opts, resolveWriters)
        if err != nil {
            return err
        }


        return nil
    },
    PreRunE: func(cmd *cobra.Command, args []string) error {
        var err error

        if opts.DnsSuffix == "" && fileOptions.DnsSuffixFile == "" {
            return errors.New("a DNS suffix or DNS suffix file must be specified")
        }

        if fileOptions.DnsSuffixFile != "" {
            if !tools.FileExists(fileOptions.DnsSuffixFile) {
                return errors.New("DNS suffix file is not readable")
            }
        }

        if fqdnOutFile != "" {
            fqdnOutFile, err = resolver.ResolveFullPath(fqdnOutFile)
            if err != nil {
                return err
            }
        }

        if err := resolveCmd.PreRunE(cmd, args); err != nil {
            return err
        }

        return nil

    },
    Run: func(cmd *cobra.Command, args []string) {

        crtshOpts := &readers.CrtShReaderOptions{
            Timeout      : 300 * time.Second,
            ProxyUri     : opts.Proxy,
        }

        dnsSuffix := []string{}
        hostWordList := []string{}
        fqdnList := []string{}
        reader := readers.NewFileReader(fileOptions)
        crtShReader := readers.NewCrtShReader(crtshOpts)
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

        if len(dnsSuffix) == 0 {
            log.Error("DNS suffix list is empty")
            os.Exit(2)
        }

        log.Debug("starting https://crt.sh crawler")
        for _, d := range dnsSuffix {
            log.Debugf("Reading dns prefix from Crt.sh to %s", d)
            if err := crtShReader.ReadFromCrtsh(d, &hostWordList, &fqdnList); err != nil {
                log.Error("error getting data from Crt.sh", "err", err)
                os.Exit(2)
            }
        }

        if len(hostWordList) == 0 {
            log.Error("DNS host list is empty")
            os.Exit(2)
        }

        total = len(dnsSuffix) * len(hostWordList)

        if fqdnOutFile != "" {
            file, err := os.OpenFile(fqdnOutFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
            if err != nil {
                log.Error("Error writting FQDN file", "err", err)
                os.Exit(2)
            }
            defer file.Close()

            for _, s := range fqdnList {
                if _, err := file.WriteString(s + "\r\n"); err != nil {
                    log.Error("Error writting FQDN file file", "line", s, "err", err)
                    os.Exit(2)
                }
            }

            log.Infof("FQDN list file saved at %s", fqdnOutFile)
        }

        log.Infof("Enumerating %s DNS hosts", tools.FormatInt(total))

        //Check DNS connectivity
        _, err := tools.GetValidDnsSuffix(fileOptions.DnsServer, "google.com.", opts.Proxy)
        if err != nil {
            log.Error("Error checking DNS connectivity", "err", err)
            os.Exit(2)
        }

        // Check runned items
        conn, _ := database.Connection(controlDb, true, false)

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
    resolveCmd.AddCommand(resolveCrtshCmd)

    resolveCrtshCmd.Flags().StringVarP(&opts.DnsSuffix, "dns-suffix", "d", "", "Single DNS suffix. (ex: helviojunior.com.br)")
    resolveCrtshCmd.Flags().StringVarP(&fileOptions.DnsSuffixFile, "dns-list", "L", "", "File containing a list of DNS suffix")
    resolveCrtshCmd.Flags().StringVar(&fqdnOutFile, "fqdn-out", "", "Output file to save requested FQDN")
}