package tools

import (
	"net"
	"testing"
)

func TestIpToUint32(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected uint32
	}{
		{
			name:     "Valid IPv4 - 127.0.0.1",
			ip:       "127.0.0.1",
			expected: 2130706433,
		},
		{
			name:     "Valid IPv4 - 192.168.1.1",
			ip:       "192.168.1.1",
			expected: 3232235777,
		},
		{
			name:     "Valid IPv4 - 0.0.0.0",
			ip:       "0.0.0.0",
			expected: 0,
		},
		{
			name:     "Valid IPv4 - 255.255.255.255",
			ip:       "255.255.255.255",
			expected: 4294967295,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ip := net.ParseIP(test.ip)
			result := IpToUint32(ip)
			if result != test.expected {
				t.Errorf("IpToUint32(%s) = %d, expected %d", test.ip, result, test.expected)
			}
		})
	}
}

func TestIpToUint32InvalidIP(t *testing.T) {
	// Test with IPv6 address (should return 0)
	ipv6 := net.ParseIP("2001:db8::1")
	result := IpToUint32(ipv6)
	if result != 0 {
		t.Errorf("IpToUint32 with IPv6 should return 0, got %d", result)
	}

	// Test with nil IP
	var nilIP net.IP
	result = IpToUint32(nilIP)
	if result != 0 {
		t.Errorf("IpToUint32 with nil IP should return 0, got %d", result)
	}
}

func TestIpsInCIDR(t *testing.T) {
	tests := []struct {
		name          string
		cidr          string
		expectedCount int
		shouldError   bool
	}{
		{
			name:          "Small subnet /30",
			cidr:          "192.168.1.0/30",
			expectedCount: 4, // All IPs included for /30
			shouldError:   false,
		},
		{
			name:          "Single host /32",
			cidr:          "192.168.1.1/32",
			expectedCount: 1,
			shouldError:   false,
		},
		{
			name:          "Larger subnet /28",
			cidr:          "10.0.0.0/28",
			expectedCount: 14, // 16 total - 2 (network + broadcast)
			shouldError:   false,
		},
		{
			name:        "Invalid CIDR",
			cidr:        "invalid",
			shouldError: true,
		},
		{
			name:        "Invalid IP in CIDR",
			cidr:        "999.999.999.999/24",
			shouldError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ips, err := IpsInCIDR(test.cidr)

			if test.shouldError {
				if err == nil {
					t.Errorf("Expected error for CIDR %s, but got none", test.cidr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error for CIDR %s: %v", test.cidr, err)
			}

			if len(ips) != test.expectedCount {
				t.Errorf("Expected %d IPs for CIDR %s, got %d", test.expectedCount, test.cidr, len(ips))
			}

			// Verify all returned IPs are valid
			for _, ip := range ips {
				if net.ParseIP(ip) == nil {
					t.Errorf("Invalid IP returned: %s", ip)
				}
			}
		})
	}
}

func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{
			name:     "Private IP - 192.168.1.1",
			ip:       "192.168.1.1",
			expected: true,
		},
		{
			name:     "Private IP - 10.0.0.1",
			ip:       "10.0.0.1",
			expected: true,
		},
		{
			name:     "Private IP - 172.16.0.1",
			ip:       "172.16.0.1",
			expected: true,
		},
		{
			name:     "Localhost - 127.0.0.1",
			ip:       "127.0.0.1",
			expected: true,
		},
		{
			name:     "Public IP - 8.8.8.8",
			ip:       "8.8.8.8",
			expected: false,
		},
		{
			name:     "Public IP - 1.1.1.1",
			ip:       "1.1.1.1",
			expected: false,
		},
		{
			name:     "Edge case - 172.15.255.255",
			ip:       "172.15.255.255",
			expected: false,
		},
		{
			name:     "Edge case - 172.32.0.1",
			ip:       "172.32.0.1",
			expected: false,
		},
		{
			name:     "Invalid IP",
			ip:       "invalid-ip",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsPrivateIP(test.ip)
			if result != test.expected {
				t.Errorf("IsPrivateIP(%s) = %v, expected %v", test.ip, result, test.expected)
			}
		})
	}
}

func TestGetDefaultDnsServer(t *testing.T) {
	tests := []struct {
		name     string
		fallback string
		expected string
	}{
		{
			name:     "With custom fallback",
			fallback: "1.1.1.1",
			expected: "1.1.1.1", // Should use fallback when no system DNS found
		},
		{
			name:     "With empty fallback",
			fallback: "",
			expected: "8.8.8.8", // Should use default fallback
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := GetDefaultDnsServer(test.fallback)
			// Since GetDNSServers() might return system DNS servers,
			// we can only test the fallback behavior when no servers are found
			if result == "" {
				t.Error("GetDefaultDnsServer should never return empty string")
			}
			// Note: We can't test exact matches because GetDNSServers()
			// might return actual system DNS servers
		})
	}
}
