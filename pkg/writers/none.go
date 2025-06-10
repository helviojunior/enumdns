package writers

import (
	"github.com/helviojunior/enumdns/pkg/models"
)

// NoneWriter is a None writer
type NoneWriter struct {
}

// NewNoneWriter initialises a none writer
func NewNoneWriter() (*NoneWriter, error) {
	return &NoneWriter{}, nil
}

// Write does nothing
func (s *NoneWriter) Write(result *models.Result) error {
	return nil
}

func (s *NoneWriter)WriteFqdn(result *models.FQDNData) error {
	return nil
}

func (s *NoneWriter) Finish() error {
    return nil
}