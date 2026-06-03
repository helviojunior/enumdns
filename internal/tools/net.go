package tools

import (
	"encoding/binary"
	"errors"
	"net"
	"net/url"
	"strings"

	//"context"
	"github.com/miekg/dns"
	"golang.org/x/net/publicsuffix"

	"github.com/helviojunior/enumdns/internal"
)

var privateNets = []string{
	"192.168.0.0/16",
	"10.0.0.0/8",
	"172.16.0.0/12",
	"127.0.0.0/8",
}

func IpToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	if ip == nil {
		return 0
	}
	return binary.BigEndian.Uint32(ip)
}

// IpsInCIDR returns a list of usable IP addresses in a given CIDR block
// excluding network and broadcast addresses for CIDRs larger than /31.
func IpsInCIDR(cidr string) ([]string, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	mask := binary.BigEndian.Uint32(ipnet.Mask)
	start := binary.BigEndian.Uint32(ipnet.IP)
	end := (start & mask) | (mask ^ 0xFFFFFFFF)

	var ips []string
	ip := make(net.IP, 4) // Preallocate buffer

	// Iterate over the range of IPs
	for i := start; i <= end; i++ {
		// For networks smaller than /30 (i.e., /29, /28, etc.), exclude network and broadcast
		ones, _ := ipnet.Mask.Size()
		if ones < 30 && (i == start || i == end) {
			continue // Skip network and broadcast addresses
		}
		binary.BigEndian.PutUint32(ip, i)
		ips = append(ips, ip.String())
	}

	return ips, nil
}

func GetValidDnsSuffix(dnsServer string, suffix string, proxyUri *url.URL) (string, error) {
	suffix = strings.Trim(suffix, ". ")
	if suffix == "" {
		return "", errors.New("empty suffix string")
	}

	suffix = strings.ToLower(suffix) + "."
	i := false

	m := new(dns.Msg)
	m.Id = dns.Id()
	m.RecursionDesired = true

	m.Question = make([]dns.Question, 1)
	m.Question[0] = dns.Question{Name: suffix, Qtype: dns.TypeSOA, Qclass: dns.ClassINET}

	c := new(internal.SocksClient)
	in, err := c.Exchange(m, proxyUri, dnsServer)
	if err != nil {
		return "", err
	} else {

		for _, ans1 := range in.Answer {
			if _, ok := ans1.(*dns.SOA); ok {
				i = true
			}
		}

	}

	if !i {
		return "", errors.New("SOA not found for domain '" + suffix + "'")
	}

	return suffix, nil

}

// IsPublicSuffix reports whether name is itself a public suffix (an entry on the
// Public Suffix List, such as "com", "com.br" or "co.uk") — a registry-operated
// zone with no registrable label, which must never be enumerated as a target.
func IsPublicSuffix(name string) bool {
	name = strings.Trim(strings.ToLower(name), ". ")
	if name == "" {
		return false
	}
	ps, _ := publicsuffix.PublicSuffix(name)
	return ps == name
}

// IsTLD reports whether name is a known, ICANN-managed public suffix (a real
// TLD such as "com", "br", "gov", "cloud" or a multi-label one like "gov.br").
// Unlike IsPublicSuffix it consults the ICANN flag, so an arbitrary unlisted
// single label (e.g. "portal", "intranet") is NOT reported as a TLD — the public
// suffix algorithm otherwise treats any unknown rightmost label as a suffix.
func IsTLD(name string) bool {
	name = strings.Trim(strings.ToLower(name), ". ")
	if name == "" {
		return false
	}
	ps, icann := publicsuffix.PublicSuffix(name)
	return icann && ps == name
}

// GetZoneApexSuffix resolves the authoritative zone apex for a DNS name that is
// not necessarily a zone apex itself (e.g. www.example.com). It walks the name
// and each of its parent domains, from most to least specific, and returns the
// first candidate that is itself a zone apex (i.e. a SOA whose owner name equals
// the candidate). Requiring the SOA owner to match the candidate avoids honoring
// a SOA returned via a CNAME chase (e.g. www -> *.github.io), which would point
// at the target's zone instead of the queried name's own zone.
// Returns the apex with a trailing dot, matching GetValidDnsSuffix.
//
// When allowPublicSuffix is false the walk stops (and fails) before reaching a
// public suffix such as com.br or co.uk, so registry-operated zones are never
// enumerated. Set allowPublicSuffix to true to explicitly permit them.
func GetZoneApexSuffix(dnsServer string, suffix string, proxyUri *url.URL, allowPublicSuffix bool) (string, error) {
	suffix = strings.Trim(strings.ToLower(suffix), ". ")
	if suffix == "" {
		return "", errors.New("empty suffix string")
	}

	labels := strings.Split(suffix, ".")
	// Walk from the most specific name up to (but not including) the bare TLD.
	for i := 0; i < len(labels)-1; i++ {
		candidate := strings.Join(labels[i:], ".")

		// Unless explicitly allowed, never walk up to (or accept) a public suffix
		// such as com.br or co.uk: those are registry-operated zones and
		// enumerating them is out of scope. Stop here, as every shorter parent is
		// also a public suffix.
		if !allowPublicSuffix && IsPublicSuffix(candidate) {
			return "", errors.New("refusing to enumerate public suffix '" + candidate + "' (no registrable parent zone found for '" + suffix + "')")
		}

		m := new(dns.Msg)
		m.Id = dns.Id()
		m.RecursionDesired = true
		m.Question = []dns.Question{{Name: candidate + ".", Qtype: dns.TypeSOA, Qclass: dns.ClassINET}}

		c := new(internal.SocksClient)
		in, err := c.Exchange(m, proxyUri, dnsServer)
		if err != nil {
			return "", err
		}

		// SOA may come in the answer (apex) or the authority (sub-name) section.
		// Only accept it when its owner name is the candidate itself, meaning the
		// candidate is a real zone apex.
		for _, ans := range append(in.Answer, in.Ns...) {
			if soa, ok := ans.(*dns.SOA); ok {
				if strings.Trim(strings.ToLower(soa.Hdr.Name), ". ") == candidate {
					return candidate + ".", nil
				}
			}
		}
	}

	return "", errors.New("SOA not found for domain '" + suffix + ".'")
}

func IsPrivateIP(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	for _, netip := range privateNets {
		_, subnet, _ := net.ParseCIDR(netip)
		if subnet.Contains(ip) {
			return true
		}
	}

	return false
}

func GetDefaultDnsServer(fallback string) string {
	if fallback == "" {
		fallback = "8.8.8.8"
	}

	srv := GetDNSServers()
	if len(srv) == 0 {
		return fallback
	}

	return srv[0].Addr().String()
}
