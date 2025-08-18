package cmd

import (
	"errors"
	//"log/slog"
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bob-reis/enumdns/internal/ascii"
	"github.com/bob-reis/enumdns/internal/tools"
	"github.com/bob-reis/enumdns/pkg/database"
	"github.com/bob-reis/enumdns/pkg/log"
	"github.com/bob-reis/enumdns/pkg/runner"

	//"github.com/bob-reis/enumdns/pkg/readers"
	resolver "github.com/helviojunior/gopathresolver"
	"github.com/spf13/cobra"
)

var zipTempFolder = ""
var resolveBloodhoundExtensions = []string{".zip", ".json"}
var resolveBloodhoundCmd = &cobra.Command{
	Use:   "bloodhound",
	Short: "Perform resolve roperations",
	Long: ascii.LogoHelp(ascii.Markdown(`
# resolve bloodhound

Perform resolver operations.
`)),
	Example: `
   - enumdns resolve bloodhound -L /tmp/bloodhound_computers.json -o enumdns.txt
   - enumdns resolve bloodhound -L /tmp/bloodhound_files.zip --write-jsonl
   - enumdns resolve bloodhound -L /tmp/bloodhound_computers.json --write-db`,
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
		var err error
		if fileOptions.HostFile == "" {
			return errors.New("a Bloodhound file must be specified")
		}

		if !tools.FileExists(fileOptions.HostFile) {
			return errors.New("bloodhound file is not readable")
		}

		fileOptions.HostFile, err = resolver.ResolveFullPath(fileOptions.HostFile)
		if err != nil {
			return err
		}

		fromExt := strings.ToLower(filepath.Ext(fileOptions.HostFile))

		if fromExt == "" {
			return errors.New("bloodhound file must have extension")
		}

		if !tools.SliceHasStr(resolveBloodhoundExtensions, fromExt) {
			return errors.New("unsupported Bloodhound file type")
		}

		if err = resolveCmd.PreRunE(cmd, args); err != nil {
			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		log.Debug("starting DNS resolver with Bloodhound computers")

		hostWordList := []string{}
		domainList := []string{}
		total := 0

		log.Debugf("Reading Bloodhound file: %s", fileOptions.HostFile)
		fromExt := strings.ToLower(filepath.Ext(fileOptions.HostFile))
		if fromExt == ".zip" {
			fileOptions.HostFile, err = getComputersFile(fileOptions.HostFile)
			if err != nil {
				log.Error("error extracting zip file", "err", err)
				os.Exit(2)
			}
		}

		err = readComputerFile(fileOptions.HostFile, &hostWordList, &domainList)

		if zipTempFolder != "" {
			if err := tools.RemoveFolder(zipTempFolder); err != nil {
				log.Error("failed to remove temp folder", "err", err)
			}
		}
		if err != nil {
			log.Error("error reading json file", "err", err)
			os.Exit(2)
		}

		total = len(hostWordList)

		if len(hostWordList) == 0 {
			log.Error("DNS host list is empty")
			os.Exit(2)
		}

		log.Infof("Checking connection with %s domain(s)", tools.FormatInt(len(domainList)))
		for _, d := range domainList {
			_, err := tools.GetValidDnsSuffix(fileOptions.DnsServer, strings.Trim(d, ". ")+".", opts.Proxy)
			if err != nil {
				log.Error("Error checking DNS connectivity. Try to ise -s option to set the DC ip", "err", err)
				if !fileOptions.IgnoreNonexistent {
					os.Exit(2)
				}
			} else {
				log.Infof("%s: OK", strings.Trim(d, ". "))
			}
		}

		log.Warnf("Enumerating %s Domains", tools.FormatInt(len(domainList)))

		reconRunner, err := runner.NewRecon(resolveRunner.GetLog(), *opts, resolveWriters)
		if err == nil {

			go func() {
				defer close(reconRunner.Targets)

				ascii.HideCursor()

				for _, d := range domainList {
					reconRunner.Targets <- d
				}

			}()

			reconRunner.Run(total)
		}

		log.Warnf("Enumerating %s DNS hosts", tools.FormatInt(total))

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

		st := resolveRunner.Run(total)
		resolveRunner.Close()

		if st.Skiped > 0 {
			log.Warnf("%d hosts were skipped because they were already scanned. Use the --force parameter to rescan them.", st.Skiped)
			ascii.ClearLine()
		}

	},
}

func init() {
	resolveCmd.AddCommand(resolveBloodhoundCmd)

	resolveBloodhoundCmd.Flags().StringVarP(&fileOptions.HostFile, "bloodhound-file", "L", "", "Bloodhound outoput file (.zip or _computers.json")
}

func getComputersFile(filePath string) (string, error) {
	var mime string
	var dst string
	var err error
	fileName := filepath.Base(filePath)
	logger := log.With("file", fileName)

	logger.Debug("Checking file")
	if mime, err = tools.GetMimeType(filePath); err != nil {
		logger.Debug("Error getting mime type", "err", err)
		return "", err
	}

	logger.Debug("Mime type", "mime", mime)
	if mime != "application/zip" {
		return "", errors.New("invalid file type")
	}

	if zipTempFolder, err = tools.CreateDir(tools.TempFileName("", "intelparser_", "")); err != nil {
		logger.Debug("Error creating temp folder to extract zip file", "err", err)
		return "", err
	}

	if dst, err = tools.CreateDirFromFilename(zipTempFolder, filePath); err != nil {
		logger.Debug("Error creating temp folder to extract zip file", "err", err)
		return "", err
	}

	if err = tools.Unzip(filePath, dst); err != nil {
		logger.Debug("Error extracting zip file", "temp_folder", dst, "err", err)
		return "", err
	}

	entries, err := os.ReadDir(dst)
	if err != nil {
		logger.Debug("Error listing folder files", "temp_folder", dst, "err", err)
		return "", err
	}

	for _, e := range entries {
		logger.Debug(e.Name())
		if strings.Contains(strings.ToLower(e.Name()), "_computers.json") {
			return filepath.Join(dst, e.Name()), nil
		}
	}

	return "", errors.New("computer file not found")

}

func readComputerFile(fileName string, outList *[]string, domainList *[]string) error {

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	br := bufio.NewReader(f)
	r, _, err := br.ReadRune()
	if err != nil {
		return err
	}
	if r != '\uFEFF' {
		if err := br.UnreadRune(); err != nil {
			return err
		} // Not a BOM -- put the rune back
	}

	fileBytes, err := io.ReadAll(br)
	if err != nil {
		return err
	}

	data := &computerFileData{}
	err = json.Unmarshal(fileBytes, data)
	if err != nil {
		return err
	}

	for _, c := range data.Data {
		n := strings.ToLower(c.Properties.Name)
		if c.Properties.Enabled {
			d := strings.ToLower(c.Properties.Domain)
			if !tools.SliceHasStr(*domainList, d) {
				*domainList = append(*domainList, d)
			}

			*outList = append(*outList, n)
		} else {
			log.Debug("Computer disabled, ignoring.", "Name", n)
		}
	}

	return nil
}

type computerDataProperties struct {
	Name    string `json:"name"`
	Domain  string `json:"domain"`
	Enabled bool   `json:"enabled"`
}

type computerData struct {
	ObjectIdentifier string                 `json:"ObjectIdentifier"`
	Properties       computerDataProperties `json:"Properties"`
}

type computerFileData struct {
	Data []computerData `json:"data"`
}
