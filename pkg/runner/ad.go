package runner

import (
	"fmt"
	"log/slog"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/helviojunior/enumdns/internal/tools"
	"github.com/helviojunior/enumdns/pkg/ad"
	"github.com/helviojunior/enumdns/pkg/database"
	"github.com/helviojunior/enumdns/pkg/log"
	"github.com/helviojunior/enumdns/pkg/models"
	"github.com/helviojunior/enumdns/pkg/writers"
	"gorm.io/gorm"
)

// adServer is a candidate LDAP endpoint discovered (or forced) for the dump.
type adServer struct {
	Host string // FQDN, used for the Kerberos SPN
	Addr string // dial address (IP) when known; empty means dial Host
}

// AdRunner dumps the DNS zones and records stored in Active Directory (ADIDNS
// style) over LDAP and persists them through the configured writers.
type AdRunner struct {
	uid     string
	log     *slog.Logger
	options Options
	writers []writers.Writer
	cfg     ad.Config
}

// NewAD builds an AD DNS dump runner from the global options. The LDAP server
// host is resolved at Run time (Global Catalog discovery or --target-ip).
func NewAD(logger *slog.Logger, opts Options, w []writers.Writer) (*AdRunner, error) {
	cfg := ad.Config{
		Port:       opts.AD.Port,
		LDAPS:      opts.AD.LDAPS,
		Username:   opts.AD.Username,
		Password:   opts.AD.Password,
		Domain:     opts.AD.Domain,
		NtlmHash:   opts.AD.NtlmHash,
		Kerberos:   opts.AD.Kerberos,
		CCachePath: opts.AD.CCachePath,
		Krb5Conf:   opts.AD.Krb5Conf,
		SPN:        opts.AD.SPN,
		BaseDN:     opts.AD.BaseDN,
		Proxy:      opts.Proxy,
		Timeout:    time.Duration(opts.Scan.Timeout) * time.Second,
	}

	return &AdRunner{
		uid:     fmt.Sprintf("%d", time.Now().UnixMilli()),
		log:     logger,
		options: opts,
		writers: w,
		cfg:     cfg,
	}, nil
}

// Run resolves the target LDAP server(s), dumps every AD-integrated DNS zone
// and writes the results. When no --target-ip is given it discovers the Global
// Catalog servers via DNS and tries each one, merging the zones found. It
// returns the number of distinct zones and records written.
func (run *AdRunner) Run() (zoneCount int, recordCount int, err error) {
	candidates := run.candidates()
	if len(candidates) == 0 {
		return 0, 0, fmt.Errorf("no LDAP server to query: set --target-ip or point -s at the AD DNS so Global Catalogs can be discovered")
	}

	seenZones := map[string]bool{}
	zones := []string{}
	connected := false
	var lastErr error

	for _, srv := range candidates {
		cfg := run.cfg
		cfg.Host = srv.Host
		cfg.Addr = srv.Addr

		run.log.Debug("trying LDAP server", "server", srv.Host, "addr", srv.Addr)
		client, err := ad.NewClient(cfg, run.log)
		if err != nil {
			run.log.Warn("could not connect/authenticate", "server", srv.Host, "err", err)
			lastErr = err
			continue
		}
		connected = true

		zlist, records, err := client.Dump()
		client.Close()
		if err != nil {
			run.log.Warn("dump failed", "server", srv.Host, "err", err)
			lastErr = err
			continue
		}

		// Persist records for zones not already written from an earlier server.
		newZones := map[string]bool{}
		for _, zr := range records {
			z := strings.Trim(strings.ToLower(zr.Zone), ". ")
			if seenZones[z] {
				continue
			}
			newZones[z] = true

			results, soa := run.toResults(zr)
			if soa != nil {
				if err := run.runWritersSOA(soa); err != nil {
					run.log.Error("failed to write SOA", "zone", zr.Zone, "err", err)
				}
			}
			for _, res := range results {
				if err := run.runWriters(res); err != nil {
					run.log.Error("failed to write result", "fqdn", res.FQDN, "err", err)
				}
				recordCount++
			}
		}

		for z := range newZones {
			seenZones[z] = true
			zones = append(zones, z)
		}

		if len(zlist) > 0 {
			log.Infof("%s: %s zone(s)", srv.Host, tools.FormatInt(len(zlist)))
		}
	}

	if !connected {
		if lastErr != nil {
			return 0, 0, lastErr
		}
		return 0, 0, fmt.Errorf("could not connect to any LDAP server")
	}

	return len(zones), recordCount, nil
}

// candidates returns the ordered list of LDAP servers to try. When --target-ip
// is set it is used exclusively; otherwise the domain controllers / Global
// Catalogs are identified by the recon module.
func (run *AdRunner) candidates() []adServer {
	if ip := strings.TrimSpace(run.options.AD.TargetIP); ip != "" {
		return []adServer{{Host: ip, Addr: ip}}
	}
	return run.discoverAdServers()
}

// discoverAdServers runs a recon pass over the AD domain (reusing the exact
// domain-controller / Global-Catalog identification the recon module already
// performs) and then reads the identified servers back from the control DB.
func (run *AdRunner) discoverAdServers() []adServer {
	domain := strings.Trim(strings.ToLower(run.options.AD.Domain), ". ")
	if domain == "" {
		run.log.Warn("no domain provided (-d); cannot discover domain controllers")
		return nil
	}

	log.Infof("Running recon on %s to identify domain controllers / global catalogs", domain)

	rec, err := NewRecon(run.log, run.options, run.writers)
	if err != nil {
		run.log.Warn("could not start recon", "err", err)
		return nil
	}
	go func() {
		defer close(rec.Targets)
		rec.Targets <- strings.Trim(domain, ". ") + "."
	}()
	// Do not call rec.Close(): the writers are shared with this runner and are
	// finalised once, by AdRunner.Close(), after the dump completes.
	rec.Run(1)

	servers := run.adServersFromDB()
	if len(servers) == 0 {
		run.log.Warn("recon did not identify any domain controller / global catalog", "domain", domain)
		return nil
	}

	names := make([]string, 0, len(servers))
	for _, s := range servers {
		names = append(names, s.Host)
	}
	log.Infof("Identified %s AD server(s): %s", tools.FormatInt(len(servers)), strings.Join(names, ", "))
	return servers
}

// adServersFromDB reads the SRV targets that recon flagged as Global Catalog
// (preferred) or Domain Controller from the control database, resolving each to
// an IP (also recorded by recon) so the LDAP dial does not depend on the local
// resolver.
func (run *AdRunner) adServersFromDB() []adServer {
	conn, err := database.Connection(run.options.Writer.CtrlDbURI, true, false)
	if err != nil {
		run.log.Warn("could not open control DB for AD discovery", "err", err)
		return nil
	}

	pickSRVTargets := func(gcOnly bool) []string {
		var targets []string
		q := conn.Model(&models.Result{}).Where("result_type = ? AND target != ?", "SRV", "")
		if gcOnly {
			q = q.Where("gc = ?", true)
		} else {
			q = q.Where("dc = ?", true)
		}
		q.Distinct().Pluck("target", &targets)
		return targets
	}

	seen := map[string]bool{}
	servers := []adServer{}
	add := func(host string) {
		host = strings.Trim(strings.ToLower(host), ". ")
		if host == "" || seen[host] {
			return
		}
		seen[host] = true
		servers = append(servers, adServer{Host: host, Addr: run.ipFromDB(conn, host)})
	}

	// Global Catalogs first, then any remaining domain controllers.
	for _, t := range pickSRVTargets(true) {
		add(t)
	}
	for _, t := range pickSRVTargets(false) {
		add(t)
	}

	return servers
}

// ipFromDB returns an IPv4 recon resolved for the given host (empty if none).
func (run *AdRunner) ipFromDB(conn *gorm.DB, fqdn string) string {
	var ips []string
	conn.Model(&models.Result{}).
		Where("LOWER(fqdn) = ? AND result_type = ? AND ipv4 != ?", fqdn, "A", "").
		Limit(1).Pluck("ipv4", &ips)
	if len(ips) > 0 {
		return ips[0]
	}
	return ""
}

// toResults converts a dnsNode (one DNS name with its records) into the hosts
// table records plus, when present, the zone SOA object.
func (run *AdRunner) toResults(zr *ad.ZoneRecord) ([]*models.Result, *models.SOA) {
	out := []*models.Result{}
	var soa *models.SOA
	now := time.Now()
	zone := strings.Trim(strings.ToLower(zr.Zone), ". ")

	for _, rec := range zr.Records {
		if rec.Type == 0x0006 { // SOA handled separately via the SOA object.
			soa = run.soaFromRecord(zr, rec, now)
			continue
		}

		res := &models.Result{
			TestId:   run.uid,
			FQDN:     zr.FQDN,
			RType:    rec.TypeName,
			SOA:      zone,
			ProbedAt: now,
			Exists:   true,
		}

		switch rec.TypeName {
		case "A":
			res.IPv4 = rec.IPv4
		case "AAAA":
			res.IPv6 = rec.IPv6
		case "NS", "CNAME", "MX", "SRV":
			res.Target = strings.Trim(strings.ToLower(rec.Target), ". ")
			run.enrich(res, res.Target)
		case "PTR":
			// Mirror the recon PTR layout: FQDN is the pointed-to host, Ptr is
			// the reverse (arpa) name and the address is derived from it.
			res.FQDN = strings.Trim(strings.ToLower(rec.Target), ". ")
			res.Ptr = zr.FQDN
			if v4, v6 := ptrToIP(zr.FQDN); v4 != "" {
				res.IPv4 = v4
			} else if v6 != "" {
				res.IPv6 = v6
			}
			run.enrich(res, res.FQDN)
		case "TXT":
			res.Txt = strings.Join(rec.Strings, " ")
		default:
			// Unknown / unsupported type: skip rather than emit an empty row.
			continue
		}

		out = append(out, res)
	}

	return out, soa
}

// soaFromRecord builds a models.SOA from a decoded AD SOA record.
func (run *AdRunner) soaFromRecord(zr *ad.ZoneRecord, rec *ad.DNSRecord, now time.Time) *models.SOA {
	s := &models.SOA{
		TestId:    run.uid,
		Name:      strings.Trim(strings.ToLower(zr.Zone), ". "),
		PrimaryNS: strings.Trim(strings.ToLower(rec.SOAPrimaryNS), ". "),
		Mbox:      strings.Trim(strings.ToLower(rec.SOAAdmin), ". "),
		Serial:    rec.SOASerial,
		Refresh:   rec.SOARefresh,
		Retry:     rec.SOARetry,
		Expire:    rec.SOAExpire,
		MinTTL:    rec.SOAMinTTL,
		ProbedAt:  now,
	}

	if cc, prodName, _ := ContainsCloudProduct(s.PrimaryNS); cc {
		s.CloudProduct = prodName
	}
	if ss, saasName, _ := ContainsSaaS(s.PrimaryNS); ss {
		s.SaaSProduct = saasName
	}
	if dc, dcName, _ := ContainsDatacenter(s.PrimaryNS); dc {
		s.Datacenter = dcName
	}

	return s
}

// enrich tags a result with the cloud / SaaS / datacenter it points at, exactly
// as the DNS recon path does.
func (run *AdRunner) enrich(res *models.Result, target string) {
	if target == "" {
		return
	}
	if cc, prodName, _ := ContainsCloudProduct(target); cc {
		res.CloudProduct = prodName
	}
	if ss, saasName, _ := ContainsSaaS(target); ss {
		res.SaaSProduct = saasName
	}
	if dc, dcName, _ := ContainsDatacenter(target); dc {
		res.Datacenter = dcName
	}
}

// runWriters passes a result to every writer.
func (run *AdRunner) runWriters(result *models.Result) error {
	for _, writer := range run.writers {
		if err := writer.Write(result); err != nil {
			run.log.Debug("Error at writer", "type", reflect.TypeOf(writer).Name(), "err", err)
			return err
		}
	}
	return nil
}

// runWritersSOA passes a SOA object to every writer.
func (run *AdRunner) runWritersSOA(soa *models.SOA) error {
	if soa == nil {
		return nil
	}
	for _, writer := range run.writers {
		if err := writer.WriteSOA(soa); err != nil {
			run.log.Debug("Error at writer", "type", reflect.TypeOf(writer).Name(), "err", err)
			return err
		}
	}
	return nil
}

// Close finalises every writer.
func (run *AdRunner) Close() {
	for _, writer := range run.writers {
		if err := writer.Finish(); err != nil {
			run.log.Error("failed to finish writer", "err", err)
		}
	}
}

// ptrToIP converts a reverse DNS name (in-addr.arpa / ip6.arpa) back to the
// address it represents. Returns (ipv4, ipv6); both empty when it is not a
// recognised reverse name.
func ptrToIP(arpa string) (string, string) {
	arpa = strings.Trim(strings.ToLower(arpa), ". ")

	if suffix := ".in-addr.arpa"; strings.HasSuffix(arpa, suffix) {
		labels := strings.Split(strings.TrimSuffix(arpa, suffix), ".")
		// Reverse the octets.
		for i, j := 0, len(labels)-1; i < j; i, j = i+1, j-1 {
			labels[i], labels[j] = labels[j], labels[i]
		}
		if len(labels) == 4 {
			if ip := net.ParseIP(strings.Join(labels, ".")); ip != nil {
				return ip.String(), ""
			}
		}
		return "", ""
	}

	if suffix := ".ip6.arpa"; strings.HasSuffix(arpa, suffix) {
		nibbles := strings.Split(strings.TrimSuffix(arpa, suffix), ".")
		if len(nibbles) != 32 {
			return "", ""
		}
		// Reverse the nibbles, then group into 8 hextets.
		for i, j := 0, len(nibbles)-1; i < j; i, j = i+1, j-1 {
			nibbles[i], nibbles[j] = nibbles[j], nibbles[i]
		}
		var b strings.Builder
		for i, n := range nibbles {
			if i > 0 && i%4 == 0 {
				b.WriteByte(':')
			}
			b.WriteString(n)
		}
		if ip := net.ParseIP(b.String()); ip != nil {
			return "", ip.String()
		}
	}

	return "", ""
}
