package writers

import "github.com/helviojunior/enumdns/pkg/models"

// Writer is a results writer
type Writer interface {
	Write(*models.Result) error
	WriteFqdn(*models.FQDNData) error
	Finish() error
}
