package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/helviojunior/enumdns/internal/ascii"
	"github.com/helviojunior/enumdns/pkg/database"
	"github.com/helviojunior/enumdns/pkg/log"
	"github.com/helviojunior/enumdns/pkg/models"
	"github.com/helviojunior/enumdns/pkg/writers"
	"github.com/spf13/cobra"
	"gorm.io/gorm/clause"
)

var rptWriters = []writers.Writer{}
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Work with enumdns reports",
	Long: ascii.LogoHelp(ascii.Markdown(`
# report

Work with enumdns reports.
`)),
}

func init() {
	rootCmd.AddCommand(reportCmd)
}

func convertFromDbTo(from string, writers []writers.Writer) error {
	log.Info("starting conversion...")

	var results = []*models.Result{}
	conn, err := database.Connection(fmt.Sprintf("sqlite:///%s", from), true, false)
	if err != nil {
		return err
	}

	if err := conn.Model(&models.Result{}).Preload(clause.Associations).Where("`exists` = ?", 1).Find(&results).Error; err != nil {
		return err
	}

	for _, result := range results {
		for _, w := range writers {
			if err := w.Write(result); err != nil {
				return err
			}
		}
	}

	log.Info("converted from a database", "rows", len(results))
	return nil
}

func openJsonlFile(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func processJsonLine(line []byte, writers []writers.Writer) error {
	var result models.Result
	if err := json.Unmarshal(line, &result); err != nil {
		log.Error("could not unmarshal JSON line", "err", err)
		return nil // Continue processing other lines
	}

	for _, w := range writers {
		if err := w.Write(&result); err != nil {
			return err
		}
	}
	return nil
}

func processJsonlFile(reader *bufio.Reader, writers []writers.Writer) (int, error) {
	count := 0

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				if len(line) == 0 {
					break // End of file
				}
				// Handle the last line without '\n'
			} else {
				return count, err
			}
		}

		if err := processJsonLine(line, writers); err != nil {
			return count, err
		}

		count++

		if err == io.EOF {
			break
		}
	}

	return count, nil
}

func convertFromJsonlTo(from string, writers []writers.Writer) error {
	if len(writers) == 0 {
		log.Warn("no writers have been configured. to persist probe results, add writers using --write-* flags")
	}

	log.Info("starting conversion...")

	file, err := openJsonlFile(from)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	count, err := processJsonlFile(reader, writers)
	if err != nil {
		return err
	}

	log.Info("converted from a JSON Lines file", "rows", count)
	return nil
}
