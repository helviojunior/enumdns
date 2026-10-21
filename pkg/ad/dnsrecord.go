// Package ad implements an ADIDNS-style dump of the DNS zones and records that
// are stored inside Active Directory (the "MicrosoftDNS" containers of the
// DomainDnsZones / ForestDnsZones application partitions). It connects over
// LDAP, enumerates every dnsZone and its dnsNode children and decodes the
// binary `dnsRecord` attribute (MS-DNSP DNS_RECORD) into usable values.
package ad

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
)

// MS-DNSP record type codes (the subset relevant to a zone dump). See
// [MS-DNSP] 2.2.2.1.1 (DNS_RECORD_TYPE).
const (
	typeZERO  uint16 = 0x0000
	typeA     uint16 = 0x0001
	typeNS    uint16 = 0x0002
	typeCNAME uint16 = 0x0005
	typeSOA   uint16 = 0x0006
	typePTR   uint16 = 0x000C
	typeMX    uint16 = 0x000F
	typeTXT   uint16 = 0x0010
	typeAAAA  uint16 = 0x001C
	typeSRV   uint16 = 0x0021
)

// typeNames maps the known record type codes to their presentation name.
var typeNames = map[uint16]string{
	typeZERO:  "ZERO",
	typeA:     "A",
	typeNS:    "NS",
	typeCNAME: "CNAME",
	typeSOA:   "SOA",
	typePTR:   "PTR",
	typeMX:    "MX",
	typeTXT:   "TXT",
	typeAAAA:  "AAAA",
	typeSRV:   "SRV",
}

// DNSRecord is a decoded value of the Active Directory `dnsRecord` attribute
// (MS-DNSP DNS_RECORD). Only the fields relevant for an ADIDNS-style dump are
// surfaced.
type DNSRecord struct {
	Type     uint16
	TypeName string
	TTL      uint32
	Serial   uint32
	Rank     uint8

	// Parsed value, which field is populated depends on Type.
	IPv4    string   // A
	IPv6    string   // AAAA
	Target  string   // NS, CNAME, PTR, MX (exchange), SRV (target)
	Strings []string // TXT

	// SOA-specific.
	SOASerial    uint32
	SOARefresh   uint32
	SOARetry     uint32
	SOAExpire    uint32
	SOAMinTTL    uint32
	SOAPrimaryNS string
	SOAAdmin     string

	// SRV-specific.
	Priority uint16
	Weight   uint16
	Port     uint16

	// MX-specific.
	Preference uint16
}

// ParseDNSRecord decodes a single raw value of the AD `dnsRecord` attribute.
// The 24-byte DNS_RECORD header is fixed; only the Data portion is type
// dependent. Unknown record types return a *DNSRecord with just the header
// populated so the caller may still report them.
func ParseDNSRecord(blob []byte) (*DNSRecord, error) {
	if len(blob) < 24 {
		return nil, errors.New("dnsRecord too short")
	}

	dataLen := int(binary.LittleEndian.Uint16(blob[0:2]))
	rType := binary.LittleEndian.Uint16(blob[2:4])

	rec := &DNSRecord{
		Type:     rType,
		TypeName: typeNames[rType],
		Rank:     blob[5],
		Serial:   binary.LittleEndian.Uint32(blob[8:12]),
		// TtlSeconds is stored big-endian, unlike the rest of the header.
		TTL: binary.BigEndian.Uint32(blob[12:16]),
	}
	if rec.TypeName == "" {
		rec.TypeName = fmt.Sprintf("TYPE%d", rType)
	}

	end := 24 + dataLen
	if end > len(blob) {
		end = len(blob)
	}
	data := blob[24:end]

	switch rType {
	case typeA:
		if len(data) >= 4 {
			rec.IPv4 = net.IP(data[0:4]).String()
		}
	case typeAAAA:
		if len(data) >= 16 {
			rec.IPv6 = net.IP(data[0:16]).String()
		}
	case typeNS, typeCNAME, typePTR:
		rec.Target, _ = parseCountName(data)
	case typeMX:
		if len(data) >= 2 {
			rec.Preference = binary.BigEndian.Uint16(data[0:2])
			rec.Target, _ = parseCountName(data[2:])
		}
	case typeSRV:
		if len(data) >= 6 {
			rec.Priority = binary.BigEndian.Uint16(data[0:2])
			rec.Weight = binary.BigEndian.Uint16(data[2:4])
			rec.Port = binary.BigEndian.Uint16(data[4:6])
			rec.Target, _ = parseCountName(data[6:])
		}
	case typeTXT:
		rec.Strings = parseTxtStrings(data)
	case typeSOA:
		if len(data) >= 20 {
			rec.SOASerial = binary.BigEndian.Uint32(data[0:4])
			rec.SOARefresh = binary.BigEndian.Uint32(data[4:8])
			rec.SOARetry = binary.BigEndian.Uint32(data[8:12])
			rec.SOAExpire = binary.BigEndian.Uint32(data[12:16])
			rec.SOAMinTTL = binary.BigEndian.Uint32(data[16:20])
			primary, consumed := parseCountName(data[20:])
			rec.SOAPrimaryNS = primary
			rec.SOAAdmin, _ = parseCountName(data[20+consumed:])
		}
	}

	return rec, nil
}

// parseCountName decodes a DNS_COUNT_NAME (MS-DNSP 2.2.2.2.1): a single length
// byte, a label-count byte, then each label prefixed by its own length byte. It
// returns the dotted name and the number of bytes the structure occupied so a
// caller (SOA) can locate the record that follows it.
func parseCountName(data []byte) (string, int) {
	if len(data) < 2 {
		return "", len(data)
	}

	nameLen := int(data[0])
	labelCount := int(data[1])
	raw := data[2:]

	labels := make([]string, 0, labelCount)
	ind := 0
	for i := 0; i < labelCount && ind < len(raw); i++ {
		l := int(raw[ind])
		ind++
		if ind+l > len(raw) {
			break
		}
		labels = append(labels, string(raw[ind:ind+l]))
		ind += l
	}

	// The structure always occupies 2 header bytes plus nameLen raw bytes.
	return strings.Join(labels, "."), 2 + nameLen
}

// parseTxtStrings decodes a DNS_RPC_RECORD_STRING (MS-DNSP 2.2.2.2.4): a
// sequence of length-prefixed character strings.
func parseTxtStrings(data []byte) []string {
	out := []string{}
	ind := 0
	for ind < len(data) {
		l := int(data[ind])
		ind++
		if ind+l > len(data) {
			break
		}
		out = append(out, string(data[ind:ind+l]))
		ind += l
	}
	return out
}
