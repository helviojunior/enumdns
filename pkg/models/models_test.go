package models

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestResultTableName(t *testing.T) {
	result := &Result{}
	if result.TableName() != "results" {
		t.Errorf("Expected table name 'results', got '%s'", result.TableName())
	}
}

func TestResultString(t *testing.T) {
	result := &Result{
		FQDN:  "example.com",
		RType: "A",
		IPv4:  "1.2.3.4",
	}

	str := result.String()
	if str == "" {
		t.Error("String() should not return empty string")
	}
}

func TestResultGetHash(t *testing.T) {
	result := &Result{
		FQDN:  "example.com",
		RType: "A",
		IPv4:  "1.2.3.4",
	}

	hash := result.GetHash()
	if hash == "" {
		t.Error("GetHash() should not return empty string")
	}

	// Test that same data produces same hash
	result2 := &Result{
		FQDN:  "example.com",
		RType: "A",
		IPv4:  "1.2.3.4",
	}

	hash2 := result2.GetHash()
	if hash != hash2 {
		t.Error("Same data should produce same hash")
	}
}

func TestResultClone(t *testing.T) {
	original := &Result{
		FQDN:         "example.com",
		RType:        "A",
		IPv4:         "1.2.3.4",
		CloudProduct: "AWS",
		ProbedAt:     time.Now(),
	}

	clone := original.Clone()

	if clone.FQDN != original.FQDN {
		t.Error("Clone should have same FQDN")
	}
	if clone.RType != original.RType {
		t.Error("Clone should have same RType")
	}
	if clone.IPv4 != original.IPv4 {
		t.Error("Clone should have same IPv4")
	}
	if clone.CloudProduct != original.CloudProduct {
		t.Error("Clone should have same CloudProduct")
	}

	// Test that modifying clone doesn't affect original
	clone.FQDN = "modified.com"
	if original.FQDN == "modified.com" {
		t.Error("Modifying clone should not affect original")
	}
}

func TestResultMarshalJSON(t *testing.T) {
	result := Result{
		FQDN:         "example.com",
		RType:        "A",
		IPv4:         "1.2.3.4",
		CloudProduct: "AWS",
		SaaSProduct:  "Gmail",
		Datacenter:   "US-East",
		DC:           true,
		GC:           false,
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var unmarshaled map[string]interface{}
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if unmarshaled["fqdn"] != "example.com" {
		t.Error("JSON should contain correct FQDN")
	}
	if unmarshaled["result_type"] != "A" {
		t.Error("JSON should contain correct result_type")
	}
	if unmarshaled["ipv4"] != "1.2.3.4" {
		t.Error("JSON should contain correct IPv4")
	}
}

func TestSliceHasResult(t *testing.T) {
	result1 := &Result{FQDN: "example.com", RType: "A", IPv4: "1.2.3.4"}
	result2 := &Result{FQDN: "test.com", RType: "A", IPv4: "5.6.7.8"}
	result3 := &Result{FQDN: "example.com", RType: "A", IPv4: "1.2.3.4"} // Same as result1

	slice := []*Result{result1, result2}

	if !SliceHasResult(slice, result3) {
		t.Error("SliceHasResult should return true for equivalent result")
	}

	result4 := &Result{FQDN: "notfound.com", RType: "A", IPv4: "9.9.9.9"}
	if SliceHasResult(slice, result4) {
		t.Error("SliceHasResult should return false for non-existent result")
	}
}

func TestFQDNDataTableName(t *testing.T) {
	fqdn := &FQDNData{}
	if fqdn.TableName() != "fqdn_results" {
		t.Errorf("Expected table name 'fqdn_results', got '%s'", fqdn.TableName())
	}
}

func TestASNTableName(t *testing.T) {
	asn := &ASN{}
	if asn.TableName() != "asn" {
		t.Errorf("Expected table name 'asn', got '%s'", asn.TableName())
	}
}

func TestASNIpDelegateTableName(t *testing.T) {
	delegate := &ASNIpDelegate{}
	if delegate.TableName() != "asn_subnets" {
		t.Errorf("Expected table name 'asn_subnets', got '%s'", delegate.TableName())
	}
}

func TestResultGetCompHash(t *testing.T) {
	result := &Result{
		FQDN:  "example.com",
		RType: "A",
		IPv4:  "1.2.3.4",
	}

	hash := result.GetCompHash()
	if hash == "" {
		t.Error("GetCompHash() should not return empty string")
	}

	// Test that same data produces same hash
	result2 := &Result{
		FQDN:  "example.com",
		RType: "A",
		IPv4:  "1.2.3.4",
	}

	hash2 := result2.GetCompHash()
	if hash != hash2 {
		t.Error("Same data should produce same hash")
	}
}

func TestASNListVariables(t *testing.T) {
	if len(ASNList) == 0 {
		t.Error("ASNList should not be empty")
	}

	if len(ASNDelegated) == 0 {
		t.Error("ASNDelegated should not be empty")
	}

	// Test that ASNList contains valid data
	for i, asn := range ASNList {
		if i >= 5 { // Test first 5 entries
			break
		}
		if asn.Number == 0 {
			t.Error("ASN Number should not be 0")
		}
		if asn.RIRName == "" {
			t.Error("ASN RIRName should not be empty")
		}
		if asn.CountryCode == "" {
			t.Error("ASN CountryCode should not be empty")
		}
	}
}

func TestResultEqual(t *testing.T) {
	tests := []struct {
		name     string
		result1  Result
		result2  Result
		expected bool
	}{
		{
			name:     "Equal A records",
			result1:  Result{FQDN: "example.com", RType: "A", IPv4: "1.2.3.4"},
			result2:  Result{FQDN: "example.com", RType: "A", IPv4: "1.2.3.4"},
			expected: true,
		},
		{
			name:     "Different A records",
			result1:  Result{FQDN: "example.com", RType: "A", IPv4: "1.2.3.4"},
			result2:  Result{FQDN: "example.com", RType: "A", IPv4: "1.2.3.5"},
			expected: false,
		},
		{
			name:     "Equal AAAA records",
			result1:  Result{FQDN: "example.com", RType: "AAAA", IPv6: "2001:db8::1"},
			result2:  Result{FQDN: "example.com", RType: "AAAA", IPv6: "2001:db8::1"},
			expected: true,
		},
		{
			name:     "Equal CNAME records",
			result1:  Result{FQDN: "example.com", RType: "CNAME", Target: "target.com"},
			result2:  Result{FQDN: "example.com", RType: "CNAME", Target: "target.com"},
			expected: true,
		},
		{
			name:     "Equal PTR records with IPv4",
			result1:  Result{FQDN: "1.3.2.1.in-addr.arpa", RType: "PTR", Ptr: "example.com", IPv4: "1.2.3.1"},
			result2:  Result{FQDN: "1.3.2.1.in-addr.arpa", RType: "PTR", Ptr: "example.com", IPv4: "1.2.3.1"},
			expected: true,
		},
		{
			name:     "Equal PTR records with IPv6",
			result1:  Result{FQDN: "example.com", RType: "PTR", Ptr: "example.com", IPv6: "2001:db8::1"},
			result2:  Result{FQDN: "example.com", RType: "PTR", Ptr: "example.com", IPv6: "2001:db8::1"},
			expected: true,
		},
		{
			name:     "Different FQDN",
			result1:  Result{FQDN: "example.com", RType: "A", IPv4: "1.2.3.4"},
			result2:  Result{FQDN: "test.com", RType: "A", IPv4: "1.2.3.4"},
			expected: false,
		},
		{
			name:     "Different RType",
			result1:  Result{FQDN: "example.com", RType: "A", IPv4: "1.2.3.4"},
			result2:  Result{FQDN: "example.com", RType: "AAAA", IPv6: "2001:db8::1"},
			expected: false,
		},
		{
			name:     "Default case with IPv6",
			result1:  Result{FQDN: "example.com", RType: "CUSTOM", IPv6: "2001:db8::1"},
			result2:  Result{FQDN: "example.com", RType: "CUSTOM", IPv6: "2001:db8::1"},
			expected: true,
		},
		{
			name:     "Default case with IPv4",
			result1:  Result{FQDN: "example.com", RType: "CUSTOM", IPv4: "1.2.3.4"},
			result2:  Result{FQDN: "example.com", RType: "CUSTOM", IPv4: "1.2.3.4"},
			expected: true,
		},
		{
			name:     "Default case with Target",
			result1:  Result{FQDN: "example.com", RType: "CUSTOM", Target: "target.com"},
			result2:  Result{FQDN: "example.com", RType: "CUSTOM", Target: "target.com"},
			expected: true,
		},
		{
			name:     "Default case with Ptr",
			result1:  Result{FQDN: "example.com", RType: "CUSTOM", Ptr: "ptr.com"},
			result2:  Result{FQDN: "example.com", RType: "CUSTOM", Ptr: "ptr.com"},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.result1.Equal(test.result2)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestResultStringFormatting(t *testing.T) {
	tests := []struct {
		name     string
		result   Result
		contains []string
	}{
		{
			name:     "A record",
			result:   Result{FQDN: "example.com", RType: "A", IPv4: "1.2.3.4"},
			contains: []string{"example.com:", "1.2.3.4"},
		},
		{
			name:     "AAAA record",
			result:   Result{FQDN: "example.com", RType: "AAAA", IPv6: "2001:db8::1"},
			contains: []string{"example.com:", "2001:db8::1"},
		},
		{
			name:     "CNAME record",
			result:   Result{FQDN: "example.com", RType: "CNAME", Target: "target.com"},
			contains: []string{"example.com:", "target.com"},
		},
		{
			name:     "PTR record",
			result:   Result{FQDN: "1.3.2.1.in-addr.arpa", RType: "PTR", Ptr: "example.com", IPv4: "1.2.3.1"},
			contains: []string{"1.3.2.1.in-addr.arpa:", "example.com", "1.2.3.1"},
		},
		{
			name:     "TXT record",
			result:   Result{FQDN: "example.com", RType: "TXT", Txt: "v=spf1 include:_spf.google.com ~all"},
			contains: []string{"example.com:", "v=spf1"},
		},
		{
			name:     "Record with Cloud product",
			result:   Result{FQDN: "example.com", RType: "A", IPv4: "1.2.3.4", CloudProduct: "AWS"},
			contains: []string{"example.com:", "1.2.3.4", "Cloud = AWS"},
		},
		{
			name:     "Record with SaaS product",
			result:   Result{FQDN: "example.com", RType: "A", IPv4: "1.2.3.4", SaaSProduct: "Google"},
			contains: []string{"example.com:", "1.2.3.4", "SaaS = Google"},
		},
		{
			name:     "Record with Datacenter",
			result:   Result{FQDN: "example.com", RType: "A", IPv4: "1.2.3.4", Datacenter: "US-East"},
			contains: []string{"example.com:", "1.2.3.4", "Datacenter = US-East"},
		},
		{
			name:     "Record with DC flag",
			result:   Result{FQDN: "example.com", RType: "A", IPv4: "1.2.3.4", DC: true},
			contains: []string{"example.com:", "1.2.3.4", "DC"},
		},
		{
			name:     "Record with GC flag",
			result:   Result{FQDN: "example.com", RType: "A", IPv4: "1.2.3.4", GC: true},
			contains: []string{"example.com:", "1.2.3.4", "GC"},
		},
		{
			name:     "Unknown record type with IPv4",
			result:   Result{FQDN: "example.com", RType: "CUSTOM", IPv4: "1.2.3.4"},
			contains: []string{"example.com:", "CUSTOM", "1.2.3.4"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			str := test.result.String()
			for _, contain := range test.contains {
				if !strings.Contains(str, contain) {
					t.Errorf("String() = %q, should contain %q", str, contain)
				}
			}
		})
	}
}

func TestResultGetCompHashDifferentTypes(t *testing.T) {
	tests := []struct {
		name     string
		result   Result
		contains string
	}{
		{
			name:     "SOA record",
			result:   Result{FQDN: "example.com", RType: "SOA"},
			contains: "000",
		},
		{
			name:     "SRV record",
			result:   Result{FQDN: "example.com", RType: "SRV"},
			contains: "010",
		},
		{
			name:     "NS record",
			result:   Result{FQDN: "example.com", RType: "NS"},
			contains: "020",
		},
		{
			name:     "CNAME record",
			result:   Result{FQDN: "example.com", RType: "CNAME"},
			contains: "030",
		},
		{
			name:     "A record",
			result:   Result{FQDN: "example.com", RType: "A"},
			contains: "040",
		},
		{
			name:     "AAAA record",
			result:   Result{FQDN: "example.com", RType: "AAAA"},
			contains: "050",
		},
		{
			name:     "PTR record",
			result:   Result{FQDN: "example.com", RType: "PTR"},
			contains: "060",
		},
		{
			name:     "Unknown record type (exists)",
			result:   Result{FQDN: "example.com", RType: "CUSTOM", Exists: true},
			contains: "900",
		},
		{
			name:     "Unknown record type (doesn't exist)",
			result:   Result{FQDN: "example.com", RType: "CUSTOM", Exists: false},
			contains: "990",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hash := test.result.GetCompHash()
			if !strings.HasPrefix(hash, test.contains) {
				t.Errorf("GetCompHash() = %q, should start with %q", hash, test.contains)
			}
		})
	}
}

func TestResultToFqdn(t *testing.T) {
	// Test with existing result
	result := Result{
		FQDN:     "example.com",
		Exists:   true,
		ProbedAt: time.Now(),
	}

	fqdn := result.ToFqdn()
	if fqdn == nil {
		t.Error("ToFqdn() should return a FQDNData when Exists is true")
		return
	}
	if fqdn.FQDN != "example.com" {
		t.Errorf("Expected FQDN 'example.com', got '%s'", fqdn.FQDN)
	}
	if fqdn.Source != "Enum" {
		t.Errorf("Expected Source 'Enum', got '%s'", fqdn.Source)
	}

	// Test with non-existing result
	resultNotExists := Result{
		FQDN:   "example.com",
		Exists: false,
	}

	fqdnNil := resultNotExists.ToFqdn()
	if fqdnNil != nil {
		t.Error("ToFqdn() should return nil when Exists is false")
	}
}
