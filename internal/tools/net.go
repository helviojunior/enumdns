package tools

import (
	"errors"
	"encoding/binary"
	"net"
	"strings"
	"net/url"
	//"context"
	"github.com/miekg/dns"

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
		// Exclude network and broadcast addresses in larger CIDR ranges
		if !(i&0xFF == 255 || i&0xFF == 0) || ipnet.Mask[3] >= 30 {
			binary.BigEndian.PutUint32(ip, i)
			ips = append(ips, ip.String())
		}
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
	m.Question[0] = dns.Question{suffix, dns.TypeSOA, dns.ClassINET}

	c := new(internal.SocksClient)
	in, err := c.Exchange(m, proxyUri, dnsServer); 
	if err != nil {
		return "", err
	}else{
		
		for _, ans1 := range in.Answer {
			if _, ok := ans1.(*dns.SOA); ok {
				i = true
			}
		}
		
	}

	if i == false {
		return "", errors.New("SOA not found for domain '"+ suffix + "'")
	}

	return suffix, nil

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
