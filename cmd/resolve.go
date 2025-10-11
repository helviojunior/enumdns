package cmd

import (
	"errors"
	"log/slog"

	//"os"
	//"fmt"

	"github.com/bob-reis/enumdns/internal/ascii"
	"github.com/bob-reis/enumdns/internal/tools"
	"github.com/bob-reis/enumdns/pkg/log"
	"github.com/bob-reis/enumdns/pkg/runner"

	//"github.com/bob-reis/enumdns/pkg/database"
	"github.com/bob-reis/enumdns/pkg/writers"
	//"github.com/bob-reis/enumdns/pkg/readers"
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

		if opts.PrivateDns {
			log.Warnf("DNS server: %s (private)", fileOptions.DnsServer)
		} else {
			log.Warn("DNS server: " + fileOptions.DnsServer)
		}

		// An slog-capable logger to use with drivers and runners
		logger := slog.New(log.Logger)

		// Configure writers using the common function
		if err := ConfigureWriters(&resolveWriters); err != nil {
			return err
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
