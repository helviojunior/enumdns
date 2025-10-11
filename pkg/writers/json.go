package writers

import (
	"encoding/json"
	"os"

	"github.com/bob-reis/enumdns/internal/tools"
	"github.com/bob-reis/enumdns/pkg/models"
)

// JsonWriter is a JSON lines writer
type JsonWriter struct {
	FilePath string
}

// NewJsonWriter return a new Json lines writer
func NewJsonWriter(destination string) (*JsonWriter, error) {
	// check if the destination exists, if not, create it
	dst, err := tools.CreateFileWithDir(destination)
	if err != nil {
		return nil, err
	}

	return &JsonWriter{
		FilePath: dst,
	}, nil
}

// Write JSON lines to a file
func (jw *JsonWriter) Write(result *models.Result) error {

	if !result.Exists {
		return nil
	}

	j, err := json.Marshal(result)
	if err != nil {
		return err
	}

	// Open the file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(jw.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	// Append the JSON data as a new line
	if _, err := file.Write(append(j, '\n')); err != nil {
		return err
	}

	/*
		fqdn := result.ToFqdn()
		if fqdn != nil {
			jw.WriteFqdn(fqdn)
		}*/

	return nil
}

func (jw *JsonWriter) WriteFqdn(result *models.FQDNData) error {

	j, err := json.Marshal(result)
	if err != nil {
		return err
	}

	// Open the file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(jw.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	// Append the JSON data as a new line
	if _, err := file.Write(append(j, '\n')); err != nil {
		return err
	}

	return nil
}

func (jw *JsonWriter) Finish() error {
	return nil
}
