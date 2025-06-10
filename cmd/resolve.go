package cmd

import (
    "errors"
    "log/slog"
    //"os"
    //"fmt"

    "github.com/helviojunior/enumdns/internal/ascii"
    "github.com/helviojunior/enumdns/internal/tools"
    "github.com/helviojunior/enumdns/pkg/log"
    "github.com/helviojunior/enumdns/pkg/runner"
    //"github.com/helviojunior/enumdns/pkg/database"
    "github.com/helviojunior/enumdns/pkg/writers"
    //"github.com/helviojunior/enumdns/pkg/readers"
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
        w, err := writers.NewDbWriter(controlDb, false)
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

        return nil
    },
    PreRunE: func(cmd *cobra.Command, args []string) error {
        //Check DNS connectivity
        _, err := tools.GetValidDnsSuffix(fileOptions.DnsServer, "google.com.", opts.Proxy)
        if err != nil {
            return errors.New("Error checking DNS connectivity: " + err.Error())
        }

        return nil
    },
}

func init() {
    rootCmd.AddCommand(resolveCmd)

    resolveCmd.PersistentFlags().BoolVarP(&fileOptions.IgnoreNonexistent, "IgnoreNonexistent", "I", false, "Ignore Nonexistent DNS suffix.")

}