package readers

import (
	"bufio"
	"errors"
	//"fmt"
	"net/url"
	"os"

	//"strconv"
	"strings"

	"github.com/helviojunior/enumdns/internal/tools"
	"github.com/helviojunior/enumdns/pkg/log"
)

// FileReader is a reader that expects a file with targets that
// is newline delimited.
type FileReader struct {
	Options *FileReaderOptions
}

// FileReaderOptions are options for the file reader
type FileReaderOptions struct {
	DnsSuffixFile     string
	HostFile          string
	DnsServer         string
	IgnoreNonexistent bool
	ProxyUri          *url.URL

	// AllowParentSOA resolves the real parent zone apex for a name that is not a
	// zone apex itself (e.g. www.example.com -> example.com) instead of treating
	// it as a non-existent suffix.
	AllowParentSOA bool

	// AddParentSOA resolves the parent zone apex and keeps BOTH the apex and the
	// informed name in the output list (used by crtsh to crawl both).
	AddParentSOA bool

	// RefusePublicSuffix rejects entries that resolve to a public suffix (e.g.
	// com.br, co.uk), preventing enumeration of registry-operated zones.
	RefusePublicSuffix bool

	// AllowTLD explicitly permits public suffixes / TLDs, overriding
	// RefusePublicSuffix and the public-suffix guard in the parent SOA walk.
	AllowTLD bool
}

// NewFileReader prepares a new file reader
func NewFileReader(opts *FileReaderOptions) *FileReader {
	return &FileReader{
		Options: opts,
	}
}

// Read from a file.
func (fr *FileReader) ReadDnsList(outList *[]string) error {

	var file *os.File
	var err error

	file, err = os.Open(fr.Options.DnsSuffixFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		candidate := scanner.Text()
		if candidate == "" {
			continue
		}

		// crtsh mode: resolve the zone apex and keep BOTH the apex (real SOA) and
		// the informed name in the list, so crt.sh is crawled for both.
		if fr.Options.AddParentSOA {
			apex, err := tools.GetZoneApexSuffix(fr.Options.DnsServer, candidate, fr.Options.ProxyUri, fr.Options.AllowTLD)
			if err != nil {
				if !fr.Options.IgnoreNonexistent {
					return err
				}
				log.Warnf("DNS suffix (%s) does not exists: %s", candidate, err.Error())
				continue
			}

			cand := strings.Trim(strings.ToLower(candidate), ". ")
			log.Infof("Resolved parent SOA for '%s': %s", cand, strings.Trim(apex, ". "))
			for _, n := range []string{apex, cand + "."} {
				if strings.Trim(n, ". ") == "" {
					continue
				}
				if !tools.SliceHasStr(*outList, n) {
					*outList = append(*outList, n)
				}
			}
			continue
		}

		//Check if DNS exists
		s, err := tools.GetValidDnsSuffix(fr.Options.DnsServer, candidate, fr.Options.ProxyUri)
		if err != nil && fr.Options.AllowParentSOA {
			// The name itself is not a zone apex; resolve the real parent SOA.
			var apex string
			apex, err = tools.GetZoneApexSuffix(fr.Options.DnsServer, candidate, fr.Options.ProxyUri, fr.Options.AllowTLD)
			if err == nil {
				log.Infof("Resolved parent SOA for '%s': %s", candidate, strings.Trim(apex, ". "))
				s = apex
			}
		}
		if err != nil {
			if !fr.Options.IgnoreNonexistent {
				return err
			}

			log.Warnf("DNS suffix (%s) does not exists: %s", candidate, err.Error())
		}

		if s == "" {
			continue
		}

		if fr.Options.RefusePublicSuffix && !fr.Options.AllowTLD && tools.IsPublicSuffix(s) {
			if !fr.Options.IgnoreNonexistent {
				return errors.New("refusing to enumerate public suffix '" + strings.Trim(s, ". ") + "'")
			}
			log.Warnf("Refusing to enumerate public suffix: %s", strings.Trim(s, ". "))
			continue
		}

		if !tools.SliceHasStr(*outList, s) {
			*outList = append(*outList, s)
		}

	}

	return scanner.Err()
}

func (fr *FileReader) ReadWordList(outList *[]string) error {
	return fr.readFileList(fr.Options.HostFile, outList)
}

// Read from a file.
func (fr *FileReader) readFileList(fileName string, outList *[]string) error {

	var file *os.File
	var err error

	file, err = os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		candidate := scanner.Text()
		if candidate == "" {
			continue
		}

		*outList = append(*outList, strings.ToLower(candidate))
	}

	return scanner.Err()
}
