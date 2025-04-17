package writers

import (
	"fmt"
	"os"

	"github.com/helviojunior/enumdns/pkg/models"
	logger "github.com/helviojunior/enumdns/pkg/log"
)

// StdoutWriter is a Stdout writer
type StdoutWriter struct {
	WriteAll  bool
}

// NewStdoutWriter initialises a stdout writer
func NewStdoutWriter() (*StdoutWriter, error) {
	return &StdoutWriter{
		WriteAll: false,
	}, nil
}

// Write results to stdout
func (s *StdoutWriter) Write(result *models.Result) error {
	fmt.Fprintf(os.Stderr, "                                                                               \r")
	if result.Failed {
		logger.Errorf("[%s] FQDN=%s", result.FailedReason, result.FQDN)
		return nil
	}

	if !result.Exists {
		return nil
	}

	if s.WriteAll {
		switch result.RType {
		case "A", "AAAA":
			logger.Infof("%s", result.String())
		case "SOA":
			if result.FQDN != result.Target {
				logger.Infof("%s", result.String())
			}
		default:
			logger.Infof("%s", result.String())		
		}
	}else{
		switch result.RType {
		case "A", "AAAA":
			logger.Infof("%s", result.String())
		default:
			logger.Debug(result.String())		
		} 
	}
	
	return nil
}

func (s *StdoutWriter) Finish() error {
    return nil
}