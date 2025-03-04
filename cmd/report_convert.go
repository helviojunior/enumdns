package cmd

import (

    "errors"
    "fmt"
    "path/filepath"
    "strings"

    "github.com/helviojunior/enumdns/internal/ascii"
    "github.com/helviojunior/enumdns/internal/islazy"
    "github.com/helviojunior/enumdns/pkg/log"
    "github.com/helviojunior/enumdns/pkg/writers"
    "github.com/spf13/cobra"
)

var conversionCmdExtensions = []string{".sqlite3", ".db", ".jsonl"}
var convertCmdFlags = struct {
    fromFile string
    toFile   string

    fromExt string
    toExt   string
}{}
var convertCmd = &cobra.Command{
    Use:   "convert",
    Short: "Convert between SQLite and JSON Lines file formats",
    Long: ascii.LogoHelp(ascii.Markdown(`
# report convert

Convert between SQLite and JSON Lines file formats.

A --from-file and --to-file must be specified. The extension used for the
specified filenames will be used to determine the conversion direction and
target.`)),
    Example: ascii.Markdown(`
- gowitness report convert --from-file gowitness.sqlite3 --to-file data.jsonl
- gowitness report convert --from-file gowitness.jsonl --to-file db.sqlite3`),
    PreRunE: func(cmd *cobra.Command, args []string) error {
        var err error
        
        if convertCmdFlags.fromFile == "" {
            return errors.New("from file not set")
        }
        if convertCmdFlags.toFile == "" {
            return errors.New("to file not set")
        }

        convertCmdFlags.fromFile, err = islazy.ResolveFullPath(convertCmdFlags.fromFile)
        if err != nil {
            return err
        }

        convertCmdFlags.toFile, err = islazy.ResolveFullPath(convertCmdFlags.toFile)
        if err != nil {
            return err
        }

        convertCmdFlags.fromExt = strings.ToLower(filepath.Ext(convertCmdFlags.fromFile))
        convertCmdFlags.toExt = strings.ToLower(filepath.Ext(convertCmdFlags.toFile))

        if convertCmdFlags.fromExt == "" || convertCmdFlags.toExt == "" {
            return errors.New("source and destination files must have extensions")
        }

        if convertCmdFlags.fromExt == convertCmdFlags.toExt {
            return errors.New("👀 source and destination file types must be different")
        }

        if convertCmdFlags.fromFile == convertCmdFlags.toFile {
            return errors.New("source and destination files cannot be the same")
        }

        if !islazy.SliceHasStr(conversionCmdExtensions, convertCmdFlags.fromExt) {
            return errors.New("unsupported from file type")
        }
        if !islazy.SliceHasStr(conversionCmdExtensions, convertCmdFlags.toExt) {
            return errors.New("unsupported to file type")
        }

        return nil
    },
    Run: func(cmd *cobra.Command, args []string) {
        var writer writers.Writer
        var err error
        if convertCmdFlags.toExt == ".sqlite3" || convertCmdFlags.toExt == ".db" {
            writer, err = writers.NewDbWriter(fmt.Sprintf("sqlite:///%s", convertCmdFlags.toFile), false)
            if err != nil {
                log.Error("could not get a database writer up", "err", err)
                return
            }
            if err := convertFromJsonlTo(convertCmdFlags.fromFile, writer); err != nil {
                log.Error("failed to convert to SQLite", "err", err)
                return
            }
        } else if convertCmdFlags.toExt == ".jsonl" {
            toFile, err := islazy.CreateFileWithDir(convertCmdFlags.toFile)
            if err != nil {
                log.Error("could not create target file", "err", err)
                return
            }
            writer, err = writers.NewJsonWriter(toFile)
            if err != nil {
                log.Error("could not get a JSON writer up", "err", err)
                return
            }
            if err := convertFromDbTo(convertCmdFlags.fromFile, writer); err != nil {
                log.Error("failed to convert to JSON Lines", "err", err)
                return
            }
        }
    },
}

func init() {
    reportCmd.AddCommand(convertCmd)

    convertCmd.Flags().StringVar(&convertCmdFlags.fromFile, "from-file", "", "The file to convert from")
    convertCmd.Flags().StringVar(&convertCmdFlags.toFile, "to-file", "", "The file to convert to. Use .sqlite3 for conversion to SQLite, and .jsonl for conversion to JSON Lines")
}
