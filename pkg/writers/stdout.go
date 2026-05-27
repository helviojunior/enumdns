package writers

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"golang.org/x/term"

	logger "github.com/helviojunior/enumdns/pkg/log"
	"github.com/helviojunior/enumdns/pkg/models"
)

// StdoutWriter is a Stdout writer
type StdoutWriter struct {
	WriteAll   bool
	IsTerminal bool
	seenSOA    map[string]struct{}
	soaMutex   sync.Mutex
}

// NewStdoutWriter initialises a stdout writer
func NewStdoutWriter() (*StdoutWriter, error) {
	return &StdoutWriter{
		WriteAll:   false,
		IsTerminal: term.IsTerminal(int(os.Stdin.Fd())),
		seenSOA:    make(map[string]struct{}),
	}, nil
}

// Write results to stdout
func (s *StdoutWriter) Write(result *models.Result) error {
	if s.IsTerminal {
		fmt.Fprintf(os.Stderr, "                                                                               \r")
	}
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
	} else {
		switch result.RType {
		case "A", "AAAA":
			logger.Infof("%s", result.String())
		default:
			logger.Debug(result.String())
		}
	}

	return nil
}

func (s *StdoutWriter) WriteFqdn(result *models.FQDNData) error {
	return nil
}

func (s *StdoutWriter) WriteSOA(soa *models.SOA) error {
	if soa == nil || strings.TrimSpace(soa.Name) == "" {
		return nil
	}

	s.soaMutex.Lock()
	// The cached SOA is persisted once per host of the zone; print it only once.
	key := soa.GetHash()
	if _, ok := s.seenSOA[key]; ok {
		s.soaMutex.Unlock()
		return nil
	}
	s.seenSOA[key] = struct{}{}
	s.soaMutex.Unlock()

	if s.IsTerminal {
		fmt.Fprintf(os.Stderr, "                                                                               \r")
	}

	name := strings.Trim(strings.ToLower(soa.Name), ". ")
	line := name + ": SOA " + soa.FormatValue() + soa.FormatSuffix()
	if s.WriteAll {
		logger.Infof("%s", line)
	} else {
		logger.Debug(line)
	}

	return nil
}

func (s *StdoutWriter) Finish() error {
	return nil
}
