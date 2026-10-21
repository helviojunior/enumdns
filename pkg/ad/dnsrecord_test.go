package ad

import (
	"encoding/binary"
	"testing"
)

// countName builds a DNS_COUNT_NAME buffer (as stored in AD, including the
// trailing root/null byte that the real directory keeps and counts in Length).
func countName(labels ...string) []byte {
	raw := []byte{}
	for _, l := range labels {
		raw = append(raw, byte(len(l)))
		raw = append(raw, []byte(l)...)
	}
	raw = append(raw, 0x00) // trailing root label, counted in Length

	out := []byte{byte(len(raw)), byte(len(labels))}
	return append(out, raw...)
}

// record wraps a type code and data payload in the 24-byte DNS_RECORD header.
func record(rType uint16, ttl uint32, data []byte) []byte {
	buf := make([]byte, 24)
	binary.LittleEndian.PutUint16(buf[0:2], uint16(len(data)))
	binary.LittleEndian.PutUint16(buf[2:4], rType)
	buf[4] = 5 // Version
	binary.BigEndian.PutUint32(buf[12:16], ttl)
	return append(buf, data...)
}

func TestParseA(t *testing.T) {
	rec, err := ParseDNSRecord(record(typeA, 600, []byte{192, 168, 1, 10}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.TypeName != "A" || rec.IPv4 != "192.168.1.10" {
		t.Fatalf("got type=%s ipv4=%s", rec.TypeName, rec.IPv4)
	}
	if rec.TTL != 600 {
		t.Fatalf("ttl got %d want 600", rec.TTL)
	}
}

func TestParseAAAA(t *testing.T) {
	data := []byte{0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}
	rec, err := ParseDNSRecord(record(typeAAAA, 300, data))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.TypeName != "AAAA" || rec.IPv6 != "2001:db8::1" {
		t.Fatalf("got type=%s ipv6=%s", rec.TypeName, rec.IPv6)
	}
}

func TestParseCNAME(t *testing.T) {
	rec, err := ParseDNSRecord(record(typeCNAME, 3600, countName("www", "corp", "local")))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.TypeName != "CNAME" || rec.Target != "www.corp.local" {
		t.Fatalf("got type=%s target=%s", rec.TypeName, rec.Target)
	}
}

func TestParseMX(t *testing.T) {
	data := []byte{0x00, 0x0a} // preference 10, big-endian
	data = append(data, countName("mail", "corp", "local")...)
	rec, err := ParseDNSRecord(record(typeMX, 3600, data))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Preference != 10 || rec.Target != "mail.corp.local" {
		t.Fatalf("got preference=%d target=%s", rec.Preference, rec.Target)
	}
}

func TestParseSRV(t *testing.T) {
	data := []byte{0x00, 0x00, 0x00, 0x64, 0x01, 0x85} // prio 0, weight 100, port 389
	data = append(data, countName("dc01", "corp", "local")...)
	rec, err := ParseDNSRecord(record(typeSRV, 600, data))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Priority != 0 || rec.Weight != 100 || rec.Port != 389 {
		t.Fatalf("got prio=%d weight=%d port=%d", rec.Priority, rec.Weight, rec.Port)
	}
	if rec.Target != "dc01.corp.local" {
		t.Fatalf("got target=%s", rec.Target)
	}
}

func TestParseTXT(t *testing.T) {
	data := []byte{}
	for _, s := range []string{"v=spf1 -all", "hello"} {
		data = append(data, byte(len(s)))
		data = append(data, []byte(s)...)
	}
	rec, err := ParseDNSRecord(record(typeTXT, 3600, data))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.Strings) != 2 || rec.Strings[0] != "v=spf1 -all" || rec.Strings[1] != "hello" {
		t.Fatalf("got strings=%v", rec.Strings)
	}
}

func TestParseSOA(t *testing.T) {
	data := make([]byte, 20)
	binary.BigEndian.PutUint32(data[0:4], 42)   // serial
	binary.BigEndian.PutUint32(data[4:8], 900)  // refresh
	binary.BigEndian.PutUint32(data[8:12], 600) // retry
	binary.BigEndian.PutUint32(data[12:16], 86400)
	binary.BigEndian.PutUint32(data[16:20], 3600)
	data = append(data, countName("dc01", "corp", "local")...)
	data = append(data, countName("hostmaster", "corp", "local")...)

	rec, err := ParseDNSRecord(record(typeSOA, 3600, data))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.SOASerial != 42 || rec.SOARefresh != 900 || rec.SOARetry != 600 ||
		rec.SOAExpire != 86400 || rec.SOAMinTTL != 3600 {
		t.Fatalf("soa timers wrong: %+v", rec)
	}
	if rec.SOAPrimaryNS != "dc01.corp.local" {
		t.Fatalf("primary ns got %q", rec.SOAPrimaryNS)
	}
	if rec.SOAAdmin != "hostmaster.corp.local" {
		t.Fatalf("admin got %q", rec.SOAAdmin)
	}
}

func TestParseUnknownType(t *testing.T) {
	rec, err := ParseDNSRecord(record(0x00FF, 60, []byte{0x01}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.TypeName != "TYPE255" {
		t.Fatalf("got %q", rec.TypeName)
	}
}

func TestParseTooShort(t *testing.T) {
	if _, err := ParseDNSRecord([]byte{0x00, 0x01}); err == nil {
		t.Fatal("expected error for short buffer")
	}
}

func TestBuildFQDN(t *testing.T) {
	cases := []struct{ name, zone, want string }{
		{"@", "corp.local", "corp.local"},
		{"", "corp.local", "corp.local"},
		{"www", "corp.local", "www.corp.local"},
		{"WWW", "Corp.Local", "www.corp.local"},
		{"www.corp.local", "corp.local", "www.corp.local"},
	}
	for _, c := range cases {
		if got := buildFQDN(c.name, c.zone); got != c.want {
			t.Errorf("buildFQDN(%q,%q)=%q want %q", c.name, c.zone, got, c.want)
		}
	}
}

func TestNtHash(t *testing.T) {
	cases := map[string]string{
		"aad3b435b51404eeaad3b435b51404ee:31d6cfe0d16ae931b73c59d7e0c089c0": "31d6cfe0d16ae931b73c59d7e0c089c0",
		"31d6cfe0d16ae931b73c59d7e0c089c0":                                  "31d6cfe0d16ae931b73c59d7e0c089c0",
	}
	for in, want := range cases {
		if got := ntHash(in); got != want {
			t.Errorf("ntHash(%q)=%q want %q", in, got, want)
		}
	}
}
