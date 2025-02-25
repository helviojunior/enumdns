package writers

import (
	"time"
	"os"

	"github.com/helviojunior/enumdns/pkg/models"
)

// StdoutWriter is a Stdout writer
type TextWriter struct {
	FilePath  string
	finalPath string
}

// NewStdoutWriter initialises a stdout writer
func NewTextWriter(destination string) (*TextWriter, error) {
	// open the file and write the CSV headers to it
	file, err := os.OpenFile(destination, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
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
	return "######################################\r\n## Date: " + time.Now().Format(time.RFC3339) + "\r\n\r\n"
}

func (t *TextWriter) Finish() error {
	file, err := os.OpenFile(t.finalPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
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

	file, err := os.OpenFile(t.finalPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(result.String() + "\r\n"); err != nil {
		return err
	}

	return nil
}
