package cmd

import (
	"errors"
	//"log/slog"
	"os"
	"strings"

	"github.com/bob-reis/enumdns/internal/ascii"
	"github.com/bob-reis/enumdns/internal/tools"
	"github.com/bob-reis/enumdns/pkg/log"

	//"github.com/bob-reis/enumdns/pkg/runner"
	"github.com/bob-reis/enumdns/pkg/database"
	"github.com/bob-reis/enumdns/pkg/readers"
	"github.com/spf13/cobra"
)

var resolveFileCmd = &cobra.Command{
	Use:   "file",
	Short: "Perform resolve roperations",
	Long: ascii.LogoHelp(ascii.Markdown(`
# resolve file

Perform resolver operations.
`)),
	Example: `
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
		if fileOptions.HostFile == "" {
			return errors.New("a hosts list file must be specified")
		}

		if !tools.FileExists(fileOptions.HostFile) {
			return errors.New("hosts list file is not readable")
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
		conn, _ := database.Connection(opts.Writer.CtrlDbURI, true, false)

		go func() {
			defer close(resolveRunner.Targets)

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

				if i || forceCheck {
					resolveRunner.Targets <- host
				} else {
					resolveRunner.AddSkiped()
				}
			}

		}()

		resolveRunner.Run(total)
		resolveRunner.Close()

	},
}

func init() {
	resolveCmd.AddCommand(resolveFileCmd)

	resolveFileCmd.Flags().StringVarP(&fileOptions.HostFile, "host-list", "L", "", "File containing a list of DNS hosts")
}
