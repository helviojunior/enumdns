package cmd

import (
    "errors"
    "path/filepath"
    "strings"

    "github.com/helviojunior/enumdns/internal/ascii"
    "github.com/helviojunior/enumdns/internal/islazy"
    "github.com/helviojunior/enumdns/pkg/log"
    "github.com/helviojunior/enumdns/pkg/writers"
    "github.com/spf13/cobra"
    
)

var elkCmdExtensions = []string{".sqlite3", ".db", ".jsonl"}
var elkCmdFlags = struct {
    fromFile string
    fromExt string
    elasticURI string
}{}
var elkCmd = &cobra.Command{
    Use:   "elastic",
    Short: "Sync from local SQLite or JSON Lines file formats to Elastic",
    Long: ascii.LogoHelp(ascii.Markdown(`
# report elastic

Sync from local SQLite or JSON Lines file formats to Elastic.

A --from-file and --elasticsearch-uri must be specified.`)),
    Example: ascii.Markdown(`
- enumdns report elastic --from-file enumdns.sqlite3 --elasticsearch-uri http://localhost:9200/enumdns
- enumdns report elastic --from-file enumdns.jsonl --elasticsearch-uri http://localhost:9200/enumdns`),
    PreRunE: func(cmd *cobra.Command, args []string) error {
        if elkCmdFlags.fromFile == "" {
            return errors.New("from file not set")
        }

        elkCmdFlags.fromFile = strings.Replace(elkCmdFlags.fromFile, "~", opts.Writer.UserPath, 2)

        elkCmdFlags.fromExt = strings.ToLower(filepath.Ext(elkCmdFlags.fromFile))

        if elkCmdFlags.fromExt == "" {
            return errors.New("source file must have extension")
        }

        if !islazy.SliceHasStr(elkCmdExtensions, elkCmdFlags.fromExt) {
            return errors.New("unsupported from file type")
        }

        return nil
    },
    Run: func(cmd *cobra.Command, args []string) {
        var writer writers.Writer
        var err error

        log.Info("Checking Elasticsearch indexes...")
        writer, err = writers.NewElasticWriter(elkCmdFlags.elasticURI)
        if err != nil {
            log.Error("could not get a elastic writer up", "err", err)
            return
        }

        if elkCmdFlags.fromExt == ".sqlite3" || elkCmdFlags.fromExt == ".db" {
            if err := convertFromDbTo(elkCmdFlags.fromFile, writer); err != nil {
                log.Error("failed to convert from SQLite", "err", err)
                return
            }
        } else if elkCmdFlags.fromExt == ".jsonl" {
            if err := convertFromJsonlTo(elkCmdFlags.fromFile, writer); err != nil {
                log.Error("failed to convert from JSON Lines", "err", err)
                return
            } 
        }
    },
}

func init() {
    reportCmd.AddCommand(elkCmd)

    elkCmd.Flags().StringVar(&elkCmdFlags.fromFile, "from-file", "~/.enumdns.db", "The file to convert from. Use .sqlite3 for conversion from SQLite, and .jsonl for conversion from JSON Lines")
    elkCmd.Flags().StringVar(&elkCmdFlags.elasticURI, "elasticsearch-uri", "http://localhost:9200/enumdns", "The elastic search URI to use. (e.g., http://user:pass@host:9200/index)")

}
