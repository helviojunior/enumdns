package cmd

import (
    "errors"
    //"log/slog"
    "os"
    "strings"

    "github.com/helviojunior/enumdns/internal/ascii"
    "github.com/helviojunior/enumdns/internal/tools"
    "github.com/helviojunior/enumdns/pkg/log"
    "github.com/helviojunior/enumdns/pkg/runner"
    "github.com/helviojunior/enumdns/pkg/database"
    "github.com/helviojunior/enumdns/pkg/writers"
    "github.com/helviojunior/enumdns/pkg/readers"
    "github.com/spf13/cobra"
)

var resolveFileRunner *runner.Runner

var resolveFileWriters = []writers.Writer{}
var resolveFileCmd = &cobra.Command{
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
        // Annoying quirk, but because I'm overriding PersistentPreRun
        // here which overrides the parent it seems.
        // So we need to explicitly call the parent's one now.
        if err := resolveCmd.PersistentPreRunE(cmd, args); err != nil {
            return err
        }

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

        if err := resolveCmd.PreRunE(cmd, args); err != nil {
            return err
        }

        return nil
    },
    Run: func(cmd *cobra.Command, args []string) {

        log.Debug("starting DNS resolver with file list")

        hostWordList := []string{}
        reader := readers.NewFileReader(fileOptions)
        total := 0

        log.Debugf("Reading dns hosts list file: %s", fileOptions.HostFile)
        if err := reader.ReadWordList(&hostWordList); err != nil {
            log.Error("error in reader.Read", "err", err)
            os.Exit(2)
        }
        total = len(hostWordList)

        if len(hostWordList) == 0 {
            log.Error("DNS host list is empty")
            os.Exit(2)
        }

        log.Infof("Enumerating %s DNS hosts", tools.FormatInt(total))

        // Check runned items
        conn, _ := database.Connection("sqlite:///" + opts.Writer.UserPath +"/.enumdns.db", true, false)

        go func() {
            defer close(resolveFileRunner.Targets)

            ascii.HideCursor()

            for _, h := range hostWordList {

                i := true
                host := strings.Trim(h, ". ") + "."
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
                    resolveFileRunner.Targets <- host
                }else{
                    resolveFileRunner.AddSkiped()
                }
            }
        
        
        }()

        resolveFileRunner.Run(total)
        resolveFileRunner.Close()

    },
}

func init() {
    resolveCmd.AddCommand(resolveFileCmd)

    resolveFileCmd.Flags().StringVarP(&fileOptions.HostFile, "host-list", "L", "", "File containing a list of DNS hosts")
}