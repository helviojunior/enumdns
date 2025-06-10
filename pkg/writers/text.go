package writers

import (
	"time"
	"os"
	"strings"

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
	txt := "######################################\r\n## Date: " + time.Now().Format(time.RFC3339) + "\r\n\r\n"
	txt += "FQDN" + strings.Repeat(" ", 67)
	txt += "Type" + strings.Repeat(" ", 7)
	txt += "Value" + strings.Repeat(" ", 50)
	txt += "\r\n"
	txt += strings.Repeat("=", 70) + " "
	txt += strings.Repeat("=", 10) + " "
	txt += strings.Repeat("=", 50)
	txt += "\r\n"

	return txt
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

	if _, err := file.WriteString(t.formatResult(result) + "\r\n"); err != nil {
		return err
	}

	return nil
}

func (t *TextWriter) formatResult(result *models.Result) string {

	r := strings.Trim(strings.ToLower(result.FQDN), ". ")
	s := 71 - len(strings.Trim(strings.ToLower(result.FQDN), ". "))
	if s <= 0 {
		s = 1
	}
	r += strings.Repeat(" ", s)

	r += result.RType
	s = 11 - len(result.RType)
	if s <= 0 {
		s = 1
	}
	r += strings.Repeat(" ", s)

	switch result.RType {
	case "A":
		r += result.IPv4
	case "AAAA":
		r += result.IPv6
	case "CNAME", "SRV", "NS", "SOA", "MX":
		r += strings.Trim(strings.ToLower(result.Target), ". ")
	case "PTR":
		r += strings.Trim(strings.ToLower(result.Ptr), ". ") + " -> "
		if result.IPv6 != "" {
			r += result.IPv6
		}else{
			r += result.IPv4
		}
	case "TXT":
		r += result.Txt
	default:
		r = r + result.RType + " "
		if result.IPv6 != "" {
			r += result.IPv6
		}else if result.IPv4 != "" {
			r += result.IPv4
		}else if result.Target != "" {
			r += strings.Trim(strings.ToLower(result.Target), ". ")
		}else if result.Ptr != "" {
			r += result.Ptr
		}
	}
	if result.CloudProduct != "" || result.SaaSProduct != "" || result.Datacenter != "" {
		prod := ""
		
		if result.CloudProduct != "" {
			prod += "Cloud = " + result.CloudProduct
		}
		if result.SaaSProduct != "" {
			if prod != "" {
				prod += ", "
			}
			prod += "SaaS = " + result.SaaSProduct
		}
		if result.Datacenter != "" {
			if prod != "" {
				prod += ", "
			}
			prod += "Datacenter = " + result.Datacenter
		}

		r += " (" + prod + ")"
	}
	if result.DC || result.GC {
		ad := []string{}
		if result.GC {
			ad = append(ad, "GC")
		}
		if result.DC {
			ad = append(ad, "DC")
		}
		r += " (" + strings.Join(ad, ", ") + ")"
	}
	return r
}