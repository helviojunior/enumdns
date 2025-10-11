package writers

import (
	"os"
	"strings"
	"time"

	"github.com/helviojunior/enumdns/pkg/models"
)

// StdoutWriter is a Stdout writer
type TextWriter struct {
	FilePath  string
	finalPath string
	seen      map[string]struct{}
	seenCand  map[string]struct{}
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
		seen:      make(map[string]struct{}),
		seenCand:  make(map[string]struct{}),
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

	// Deduplicate by composite hash (type+value)
	key := result.GetCompHash()
	if _, ok := t.seen[key]; ok {
		return nil
	}
	t.seen[key] = struct{}{}

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
	if result == nil || strings.TrimSpace(result.FQDN) == "" {
		return nil
	}

	fqdn := strings.Trim(strings.ToLower(result.FQDN), ". ")
	if _, ok := t.seenCand[fqdn]; ok {
		return nil
	}
	t.seenCand[fqdn] = struct{}{}

	file, err := os.OpenFile(t.finalPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	line := fqdn
	// pad to FQDN column width ~71
	s := 71 - len(fqdn)
	if s <= 0 {
		s = 1
	}
	line += strings.Repeat(" ", s)
	// mark as candidate
	rtype := "CANDIDATE"
	line += rtype
	s = 11 - len(rtype)
	if s <= 0 {
		s = 1
	}
	line += strings.Repeat(" ", s)
	line += "generated"

	if _, err := file.WriteString(line + "\r\n"); err != nil {
		return err
	}

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
