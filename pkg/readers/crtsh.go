package readers

import (
	"fmt"
	"net/url"
	"errors"
	"strings"

	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/helviojunior/enumdns/internal/tools"
	"github.com/helviojunior/enumdns/pkg/log"
)

type CrtShReader struct {
	Options *CrtShReaderOptions
}

type CrtShReaderOptions struct {
	Timeout         time.Duration
	ProxyUri 		*url.URL
}

func NewCrtShReader(opts *CrtShReaderOptions) *CrtShReader {
	if opts.Timeout <= (60 * time.Second) {
		opts.Timeout = (60 * time.Second)
	}
	return &CrtShReader{
		Options: opts,
	}
}

// Read from a https://crt.sh.
func (crtr *CrtShReader) ReadFromCrtsh(domain string, outList *[]string, fqdnList *[]string) error {
	domain = strings.Trim(strings.ToLower(domain), ".")
	crtUrl := fmt.Sprintf("https://crt.sh/?CN=%s", domain)

	client := &http.Client{
		Timeout:   crtr.Options.Timeout,
	}

	if crtr.Options.ProxyUri != nil {
		// Create transport with proxy
		transport := &http.Transport{
			Proxy: http.ProxyURL(crtr.Options.ProxyUri),
		}

		// Create HTTP client with timeout and transport
		client = &http.Client{
			Timeout:   crtr.Options.Timeout,
			Transport: transport,
		}
	}

	resp, err := crtr.fetchWithRetry(client, crtUrl) // client.Get(crtUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	*outList = append(*outList, domain)

	// Extract all <TD>...</TD> content
	re := regexp.MustCompile(`(?i)<TD>(.*?)</TD>`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	for _, match := range matches {
		candidate := strings.ToLower(strings.TrimSpace(strings.Replace(match[1], "*.", "", -1)))
		if candidate != "" && !strings.Contains(candidate, "white-space:normal") {
			// Check if it is a valid FQDN
			_, err := url.Parse(fmt.Sprintf("https://%s/", candidate))
        	if err != nil {
	        	log.Debug("Invalid host", "host", candidate, "err", err)
	        }else{
	        	candidate = strings.Trim(candidate, ".")
	        	candidate = strings.Replace(candidate, fmt.Sprintf(".%s", domain), "", -1)
	        	candidate = strings.Replace(candidate, domain, "", -1)
	        	if candidate != "" {
		        	if !tools.SliceHasStr(*outList, candidate) {
		        		log.Debug("Match", "domain", domain, "host", candidate)
			        	*outList = append(*outList, candidate)
			        }
			        fqdn := fmt.Sprintf("%s.%s", candidate, domain)
			        if !tools.SliceHasStr(*fqdnList, fqdn) {
			        	*fqdnList = append(*fqdnList, fqdn)
			        }
			    }
	        }
		}
	}

	return nil
}

func (crtr *CrtShReader) fetchWithRetry(client *http.Client, crtUrl string) (*http.Response, error) {
	var resp *http.Response
	var err error

	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		resp, err = client.Get(crtUrl)
		if err == nil {
			return resp, nil
		}
		if i < maxRetries-1 {
			time.Sleep(time.Second * time.Duration(10 * i))
		}
	}

	return nil, errors.New("failed after 3 retries: " + err.Error())
}