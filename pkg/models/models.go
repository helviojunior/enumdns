package models

import (
	"time"
	"encoding/json"
	"strings"
	//"math/big"
	"net"

	"fmt"
	"crypto/sha1"
	"encoding/hex"

	"github.com/helviojunior/enumdns/internal/tools"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ASN struct {
    Number int64 `gorm:"primarykey;column:asn" json:"asn"`
    RIRName               string    `gorm:"column:rir_name" json:"rir_name"`
    CountryCode           string    `gorm:"column:country_code" json:"country_code"`
    Org                   string    `gorm:"column:org" json:"org"`
}

func (*ASN) TableName() string {
    return "asn"
}

func (asn *ASN) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{
		//Columns:   cols,
		Columns:   []clause.Column{{Name: "asn"}},
		UpdateAll: true,
	})
	return nil
}

type ASNIpDelegate struct {
    ID uint `json:"id" gorm:"primarykey"`

    Hash                  string    `gorm:"column:hash;index:,unique;" json:"hash"`
    RIRName               string    `gorm:"column:rir_name" json:"rir_name"`
    CountryCode           string    `gorm:"column:country_code" json:"country_code"`
    Subnet                string    `gorm:"column:subnet" json:"subnet"` 
    IntIPv4				  int64     `gorm:"column:int_ipv4" json:"int_ipv4"` 
    Addresses             int       `gorm:"column:addresses" json:"addresses"`
    Date                  string    `gorm:"column:date" json:"date"`
    ASN                   int64     `gorm:"column:asn" json:"asn"`
    Status                string    `gorm:"column:status" json:"status"`
}

func (*ASNIpDelegate) TableName() string {
    return "asn_subnets"
}

func (subnet *ASNIpDelegate) BeforeCreate(tx *gorm.DB) (err error) {
	_calcHash(&subnet.Hash, subnet.Subnet, fmt.Sprintf("%d", subnet.Addresses))

	tx.Statement.AddClause(clause.OnConflict{
		//Columns:   cols,
		Columns:   []clause.Column{{Name: "hash"}},
		UpdateAll: true,
	})
	return nil
}

// Result is a github.com/helviojunior/enumdnsenumdns result
type Result struct {
	ID uint `json:"id" gorm:"primarykey"`

	TestId                string    `gorm:"column:test_id"`
	Hash				  string    `gorm:"column:hash;index:,unique;"`
	FQDN                  string    `gorm:"column:fqdn"`
	RType                 string    `gorm:"column:result_type"`
	IPv4                  string    `gorm:"column:ipv4"`
	IPv6                  string    `gorm:"column:ipv6"`
	ASN                   int64     `gorm:"column:asn"`
	Target                string    `gorm:"column:target"`
	Ptr                   string    `gorm:"column:ptr"`
	Txt                   string    `gorm:"column:txt"`
	CloudProduct          string    `gorm:"column:cloud_product"`
	SaaSProduct           string    `gorm:"column:saas_product"`
	Datacenter            string    `gorm:"column:datacenter"`
	ProbedAt              time.Time `gorm:"column:probed_at"`

	DC      	 		  bool   	`gorm:"column:dc"`
	GC  	       		  bool   	`gorm:"column:gc"`

	Exists       		  bool   	`gorm:"column:exists"`

	// Failed flag set if the result should be considered failed
	Failed       bool   `gorm:"column:failed" gorm:"index:idx_exists"`
	FailedReason string `gorm:"column:failed_reason"`

}

func (*Result) TableName() string {
    return "results"
}

func (result *Result) BeforeCreate(tx *gorm.DB) (err error) {
	_calcHash(&result.Hash, result.String())
	asn := result.GetASN(tx)

	if asn != nil {
		result.ASN = asn.ASN
	}

	tx.Statement.AddClause(clause.OnConflict{
		//Columns:   cols,
		Columns:   []clause.Column{{Name: "hash"}},
		UpdateAll: true,
	})
	return nil
}

func (result *Result) GetASN(tx *gorm.DB) *ASNIpDelegate {
	var asn *ASNIpDelegate
	if result.IPv4 != "" {
		ip := net.ParseIP(result.IPv4)
		if ip == nil {
			return nil
		}

		err := tx.Model(&ASNIpDelegate{}).
		    Where("int_ipv4 != 0 AND int_ipv4 <= ?", tools.IpToUint32(ip)).
		    Order("int_ipv4 DESC"). // Optional: get the closest (largest) match <= ip
		    First(&asn).Error

		if err != nil {
			return nil
		}

		_, subnet, err := net.ParseCIDR(asn.Subnet)
        if err != nil {
            return nil
        }

         if !subnet.Contains(ip) {
            return nil
        }
		
	}else if result.IPv6 != "" {
		/*
		ip := net.ParseIP(result.IPv6)
		if ip == nil {
			return nil
		}

		iIp := ip.To16()
		if iIp == nil {
			return nil
		}
		err := tx.Model(&ASNIpDelegate{}).
		    Where("int_ipv6 <= ?", iIp).
		    Order("int_ipv6 DESC"). // Optional: get the closest (largest) match <= ip
		    First(&asn).Error

		if err != nil {
			return nil
		}

		_, subnet, err := net.ParseCIDR(asn.Subnet)
        if err != nil {
            return nil
        }

         if !subnet.Contains(ip) {
            return nil
        }*/
		
	}

	return asn
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
		Txt                   string    `json:"txt,omitempty"`
		CloudProduct          string    `json:"cloud_product,omitempty"`
		SaaSProduct           string    `json:"saas_product,omitempty"`
		Datacenter            string    `json:"datacenter,omitempty"`
		DC      	 		  bool   	`json:"dc"`
		GC  	       		  bool   	`json:"gc"`
		ProbedAt              string    `json:"probed_at"`

	}{
		FQDN 				: strings.Trim(strings.ToLower(result.FQDN), ". "),
		RType 				: strings.ToUpper(result.RType),
		ProbedAt    		: result.ProbedAt.Format(time.RFC3339),
		IPv4 				: result.IPv4,
		IPv6 				: result.IPv6,
		Target 				: strings.Trim(strings.ToLower(result.Target), ". "),
		Ptr 				: strings.Trim(strings.ToLower(result.Ptr), ". "),
		Txt 				: result.Txt,
		DC 					: result.DC,
		GC 					: result.GC,
		CloudProduct 		: result.CloudProduct,
		SaaSProduct 		: result.SaaSProduct,
		Datacenter 		    : result.Datacenter,
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
		Txt 				: result.Txt,
		DC 					: result.DC,
		GC 					: result.GC,
		CloudProduct 		: result.CloudProduct,
		SaaSProduct 		: result.SaaSProduct,
		Datacenter 		    : result.Datacenter,
		ProbedAt 			: result.ProbedAt,
		Exists 				: result.Exists,
		Failed 				: result.Failed,
		FailedReason 		: result.FailedReason,
	}
}

type FQDNData struct {
	ID uint `json:"id" gorm:"primarykey"`

	Hash				  string    `gorm:"column:hash;index:,unique;"`
	FQDN                  string    `gorm:"column:fqdn"`
	Source                string    `gorm:"column:source"`
	ProbedAt              time.Time `gorm:"column:probed_at"`
}

func (*FQDNData) TableName() string {
    return "fqdn_results"
}

func (fqdn *FQDNData) BeforeCreate(tx *gorm.DB) (err error) {
	_calcHash(&fqdn.Hash, fqdn.FQDN)

	tx.Statement.AddClause(clause.OnConflict{
		//Columns:   cols,
		Columns:   []clause.Column{{Name: "hash"}},
		DoNothing: true,
	})
	return nil
}

func (result Result) ToFqdn() *FQDNData {

	if !result.Exists {
		return nil
	}

	return &FQDNData{
		FQDN 				: strings.Trim(strings.ToLower(result.FQDN), ". "),
		Source 				: "Enum",
		ProbedAt 			: result.ProbedAt,
	}
}

func (result Result) Equal(r1 Result) bool {
	if result.RType != r1.RType {
		return false
	}
	if result.FQDN != r1.FQDN {
		return false
	}
	switch result.RType {
	case "A":
		return result.IPv4 == r1.IPv4
	case "AAAA":
		return result.IPv6 == r1.IPv6
	case "CNAME", "SRV", "NS", "SOA":
		return strings.Trim(strings.ToLower(result.Target), ". ") == strings.Trim(strings.ToLower(r1.Target), ". ")
	case "PTR":
		r2 := strings.Trim(strings.ToLower(result.Ptr), ". ") == strings.Trim(strings.ToLower(r1.Ptr), ". ")
		if result.IPv6 != "" {
			return result.IPv6 == r1.IPv6 && r2
		}else{
			return result.IPv4 == r1.IPv4 && r2
		}
	default:
		if result.IPv6 != "" {
			return result.IPv6 == r1.IPv6
		}else if result.IPv4 != "" {
			return result.IPv4 == r1.IPv4
		}else if result.Target != "" {
			return strings.Trim(strings.ToLower(result.Target), ". ") == strings.Trim(strings.ToLower(r1.Target), ". ")
		}else if result.Ptr != "" {
			return strings.Trim(strings.ToLower(result.Ptr), ". ") == strings.Trim(strings.ToLower(r1.Ptr), ". ")
		}
	}

	return false
}

func (result Result) String() string {
	r := strings.Trim(strings.ToLower(result.FQDN), ". ") + ": "
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

func (result Result) GetHash() string {
	b_data := []byte(result.String())
	return tools.GetHash(b_data)
}

func (result Result) GetCompHash() string {
	r := ""
	switch result.RType {
	case "SOA":
		r += "000"
	case "SRV":
		r += "010"
	case "NS":
		r += "020"
	case "CNAME":
		r += "030"
	case "A":
		r += "040"
	case "AAAA":
		r += "050"
	case "PTR":
		r += "060"
	default:
		if !result.Exists {
			r += "990"
		}else{
			r += "900"
		}
	}

	r += result.String()
	return r
}

func SliceHasResult(s []*Result, r *Result) bool {
    for _, a := range s {
    	if r.Equal(*a) {
    		return true
    	}
    	/*
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
    	}*/
    }
    return false
}


func _calcHash(outValue *string, keyvals ...interface{}) {

	data := ""
	for _, v := range keyvals {
		if _, ok := v.(int); ok {
			data += fmt.Sprintf("%d,", v)
		}else{
			data += fmt.Sprintf("%s,", v)
		}
	}

	h := sha1.New()
	h.Write([]byte(data))

	*outValue = hex.EncodeToString(h.Sum(nil))

}