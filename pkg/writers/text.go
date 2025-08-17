package writers

import (
	"os"
	"strings"
	"time"

	"github.com/bob-reis/enumdns/pkg/models"
)

// StdoutWriter is a Stdout writer
type TextWriter struct {
	FilePath  string
	finalPath string
}

// NewStdoutWriter initialises a stdout writer
func NewTextWriter(destination string) (*TextWriter, error) {
	// open the file and write the CSV headers to it
	file, err := os.OpenFile(destination, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if _, err := file.WriteString(txtHeader()); err != nil {
		return nil, err
	}

	return &TextWriter{
		FilePath:  destination,
		finalPath: destination,
	}, nil
}

func txtHeader() string {
	txt := "######################################\r\n## Date: " + time.Now().Format(time.RFC3339) + "\r\n\r\n"
	txt += "FQDN" + strings.Repeat(" ", 67)
	txt += "Type" + strings.Repeat(" ", 7)
	txt += "Value" + strings.Repeat(" ", 50)
	txt += "\r\n"
	txt += strings.Repeat("=", 70) + " "
	txt += strings.Repeat("=", 10) + " "
	txt += strings.Repeat("=", 50)
	txt += "\r\n"

	return txt
}

func (t *TextWriter) Finish() error {
	file, err := os.OpenFile(t.finalPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString("\r\nFinished at: " + time.Now().Format(time.RFC3339) + "\r\n\r\n"); err != nil {
		return err
	}

	return nil
}

// Write results to stdout
func (t *TextWriter) Write(result *models.Result) error {

	if !result.Exists {
		return nil
	}

	file, err := os.OpenFile(t.finalPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(t.formatResult(result) + "\r\n"); err != nil {
		return err
	}

	return nil
}

func (t *TextWriter) WriteFqdn(result *models.FQDNData) error {
	return nil
}

func (t *TextWriter) formatResult(result *models.Result) string {
	fqdn := strings.Trim(strings.ToLower(result.FQDN), ". ")
	r := fqdn

	// Add spacing for FQDN column (71 chars)
	s := 71 - len(fqdn)
	if s <= 0 {
		s = 1
	}
	r += strings.Repeat(" ", s)

	// Add type column with spacing (11 chars)
	r += result.RType
	s = 11 - len(result.RType)
	if s <= 0 {
		s = 1
	}
	r += strings.Repeat(" ", s)

	// Add value using shared formatting logic
	r += result.FormatValue()

	// Add suffix using shared formatting logic
	r += result.FormatSuffix()

	return r
}
