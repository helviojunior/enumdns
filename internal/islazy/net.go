package islazy

import (
	"errors"
	"encoding/binary"
	"net"
	"strings"
	//"fmt"
	"github.com/miekg/dns"
)

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

func GetValidDnsSufix(dnsServer string, sufix string) (string, error) {
	sufix = strings.Trim(sufix, ". ")
	if sufix == "" {
		return "", errors.New("empty sufix string")
	}

	sufix = sufix + "."
	i := false

    m := new(dns.Msg)
    m.Id = dns.Id()
	m.RecursionDesired = true

	m.Question = make([]dns.Question, 1)
	m.Question[0] = dns.Question{sufix, dns.TypeSOA, dns.ClassINET}

	in, err := dns.Exchange(m, dnsServer); 
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
		return "", errors.New("SOA not found for domain '"+ sufix + "'")
	}

	return sufix, nil

}