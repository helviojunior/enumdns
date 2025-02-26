package models

import (
	"time"
	"encoding/json"
	//"fmt"
	"strings"

	"github.com/helviojunior/enumdns/internal/islazy"
)

// Result is a github.com/helviojunior/enumdnsenumdns result
type Result struct {
	ID uint `json:"id" gorm:"primarykey"`

	TestId                string    `json:"test_id"`
	FQDN                  string    `json:"fqdn" gorm:"index:idx_exists"`
	RType                 string    `json:"result_type"`
	IPv4                  string    `json:"ipv4,omitempty"`
	IPv6                  string    `json:"ipv6,omitempty"`
	Target                string    `json:"target,omitempty"`
	Ptr                   string    `json:"ptr,omitempty"`
	CloudProduct          string    `json:"cloud_product,omitempty"`
	ProbedAt              time.Time `json:"probed_at"`

	Exists       		  bool   	`json:"exists"`

	// Failed flag set if the result should be considered failed
	Failed       bool   `json:"failed" gorm:"index:idx_exists"`
	FailedReason string `json:"failed_reason,omitempty"`

}

/* Custom Marshaller for Result */
func (result Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		FQDN                  string    `json:"fqdn"`
		RType                 string    `json:"result_type"`
		IPv4                  string    `json:"ipv4,omitempty"`
		IPv6                  string    `json:"ipv6,omitempty"`
		Target                string    `json:"target,omitempty"`
		Ptr                   string    `json:"ptr,omitempty"`
		CloudProduct          string    `json:"cloud_product,omitempty"`
		ProbedAt              string    `json:"probed_at"`

	}{
		FQDN 				: strings.Trim(strings.ToLower(result.FQDN), ". "),
		RType 				: strings.ToUpper(result.RType),
		ProbedAt    		: result.ProbedAt.Format(time.RFC3339),
		IPv4 				: result.IPv4,
		IPv6 				: result.IPv6,
		Target 				: strings.Trim(strings.ToLower(result.Target), ". "),
		Ptr 				: strings.Trim(strings.ToLower(result.Ptr), ". "),
		CloudProduct 		: result.CloudProduct,
	})
}

func (result Result) Clone() *Result {
	return &Result{
		TestId 				: result.TestId,
		FQDN 				: result.FQDN,
		RType 				: result.RType,
		IPv4 				: result.IPv4,
		IPv6 				: result.IPv6,
		Target 				: result.Target,
		Ptr 				: result.Ptr,
		CloudProduct 		: result.CloudProduct,
		ProbedAt 			: result.ProbedAt,
		Exists 				: result.Exists,
		Failed 				: result.Failed,
		FailedReason 		: result.FailedReason,
	}
}

func (result Result) String() string {
	r := strings.Trim(strings.ToLower(result.FQDN), ". ") + ": "
	switch result.RType {
	case "A":
		r += result.IPv4
	case "AAAA":
		r += result.IPv6
	case "CNAME":
		r += strings.Trim(strings.ToLower(result.Target), ". ")
	case "PTR":
		r += strings.Trim(strings.ToLower(result.Ptr), ". ") + " -> "
		if result.IPv6 != "" {
			r += result.IPv6
		}else{
			r += result.IPv4
		}
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
	if result.CloudProduct != "" {
		r += " (Cloud = " + result.CloudProduct + ")"
	}
	return r
}

func (result Result) GetHash() string {
	b_data := []byte(result.String())
	return islazy.GetHash(b_data)
}

func SliceHasResult(s []*Result, r *Result) bool {
    for _, a := range s {
    	if a.FQDN != r.FQDN || a.RType != r.RType || a.Ptr != r.Ptr {
    		continue
    	}
    	switch a.RType {
    	case "A":
    		if a.IPv4 == r.IPv4 {
    			return true
    		}
    	case "AAAA":
    		if a.IPv6 == r.IPv6 {
    			return true
    		}
    	case "CNAME":
    		if a.Target == r.Target {
    			return true
    		}
    	}
    }
    return false
}