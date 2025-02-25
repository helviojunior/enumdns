package readers

import (
	"bufio"
	//"fmt"
	//"net/url"
	"os"
	//"strconv"
	//"strings"

	"github.com/helviojunior/enumdns/internal/islazy"
)

// FileReader is a reader that expects a file with targets that
// is newline delimited.
type FileReader struct {
	Options *FileReaderOptions
}

// FileReaderOptions are options for the file reader
type FileReaderOptions struct {
	DnsSufixFile    string
	HostFile		string
	DnsServer 		string
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

	file, err = os.Open(fr.Options.DnsSufixFile)
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
		s, err := islazy.GetValidDnsSufix(fr.Options.DnsServer, candidate)
		if err != nil {
			return err
		}

		*outList = append(*outList, s)

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

		*outList = append(*outList, candidate)
	}

	return scanner.Err()
}
