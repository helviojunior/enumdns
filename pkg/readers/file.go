package readers

import (
	"bufio"
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

		//Check if DNS exists
		s, err := tools.GetValidDnsSuffix(fr.Options.DnsServer, candidate, fr.Options.ProxyUri)
		if err != nil && fr.Options.AllowParentSOA {
			// The name itself is not a zone apex; resolve the real parent SOA.
			var apex string
			apex, err = tools.GetZoneApexSuffix(fr.Options.DnsServer, candidate, fr.Options.ProxyUri)
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
