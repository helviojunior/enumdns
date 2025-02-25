package writers

import (
	"fmt"
	"os"

	"github.com/helviojunior/enumdns/pkg/models"
	logger "github.com/helviojunior/enumdns/pkg/log"
)

// StdoutWriter is a Stdout writer
type StdoutWriter struct {
}

// NewStdoutWriter initialises a stdout writer
func NewStdoutWriter() (*StdoutWriter, error) {
	return &StdoutWriter{}, nil
}

// Write results to stdout
func (s *StdoutWriter) Write(result *models.Result) error {
	fmt.Fprintf(os.Stderr, "                                                                               \r")
	if result.Failed {
		logger.Errorf("[%s] FQDN=%s", result.FailedReason, result.FQDN)
		return nil
	}

	switch result.RType {
	case "A", "AAAA":
		logger.Infof("%s", result.String())
	default:
		logger.Debug(result.String())
	} 
	
	return nil
}

func (s *StdoutWriter) Finish() error {
    return nil
}