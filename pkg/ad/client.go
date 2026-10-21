package ad

import (
	"crypto/tls"
	"errors"
	"log/slog"
	"net"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/go-ldap/ldap/v3/gssapi"
	"github.com/helviojunior/enumdns/internal"
)

// pageSize is the LDAP paging size used when dumping records; AD zones can hold
// tens of thousands of nodes and the server caps an unpaged search.
const pageSize = 1000

// Config holds everything needed to connect, authenticate and enumerate AD DNS.
type Config struct {
	Host  string // DC host/FQDN; used for the Kerberos SPN and TLS SNI (required)
	Addr  string // optional dial address (IP) when it differs from Host
	Port  int    // LDAP port (default 389, or 636 with LDAPS)
	LDAPS bool   // use TLS (LDAPS)

	// Authentication. Exactly one method is selected, in this order of
	// precedence: Kerberos -> NTLM hash -> user/password -> anonymous.
	Username string
	Password string
	Domain   string // AD domain / realm (NetBIOS or FQDN)

	NtlmHash string // NTLM hash for pass-the-hash ("NT" or "LM:NT")

	Kerberos   bool   // authenticate with a Kerberos ticket cache
	CCachePath string // ccache file (falls back to $KRB5CCNAME)
	Krb5Conf   string // krb5.conf path (falls back to /etc/krb5.conf)
	SPN        string // service principal override (default ldap/<host>)

	BaseDN string // optional: dump only this container instead of autodetecting

	Proxy   *url.URL
	Timeout time.Duration
}

// Client is a connected and bound LDAP session against a domain controller.
type Client struct {
	conn *ldap.Conn
	gss  *gssapi.Client
	cfg  Config
	log  *slog.Logger
}

// searchBase is an LDAP base DN plus the scope to search it with when looking
// for dnsZone objects.
type searchBase struct {
	dn    string
	scope int
}

// ZoneRecord is a single dnsNode (one DNS name) together with its decoded
// records. A node may carry several records (e.g. an apex with SOA + NS + A).
type ZoneRecord struct {
	Zone    string // zone apex, e.g. example.com
	Name    string // relative record name ("@" for the apex)
	FQDN    string // fully qualified name
	Records []*DNSRecord
}

// NewClient dials the domain controller, performs the configured bind and
// returns a ready-to-use session. The caller must Close it.
func NewClient(cfg Config, logger *slog.Logger) (*Client, error) {
	if cfg.Port == 0 {
		if cfg.LDAPS {
			cfg.Port = 636
		} else {
			cfg.Port = 389
		}
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}

	dialHost := cfg.Host
	if cfg.Addr != "" {
		dialHost = cfg.Addr
	}
	address := net.JoinHostPort(dialHost, strconv.Itoa(cfg.Port))

	rawConn, err := dial(cfg, address)
	if err != nil {
		return nil, errors.New("error connecting to " + address + ": " + err.Error())
	}

	var lconn *ldap.Conn
	if cfg.LDAPS {
		tconn := tls.Client(rawConn, &tls.Config{
			ServerName: cfg.Host,
			// AD DCs routinely present self-signed/internal-CA certificates; a
			// recon dump should not fail on an unverifiable chain.
			InsecureSkipVerify: true, //nolint:gosec
		})
		if err := tconn.Handshake(); err != nil {
			_ = rawConn.Close()
			return nil, errors.New("TLS handshake failed: " + err.Error())
		}
		lconn = ldap.NewConn(tconn, true)
	} else {
		lconn = ldap.NewConn(rawConn, false)
	}

	lconn.Start()
	lconn.SetTimeout(cfg.Timeout)

	c := &Client{conn: lconn, cfg: cfg, log: logger}
	if err := c.bind(); err != nil {
		c.Close()
		return nil, err
	}

	return c, nil
}

// Close releases the LDAP connection and any Kerberos context.
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
	if c.gss != nil {
		_ = c.gss.Close()
	}
}

// dial establishes the raw TCP connection, optionally through a SOCKS proxy
// (reusing the project's shared proxy helper).
func dial(cfg Config, address string) (net.Conn, error) {
	forward := &net.Dialer{Timeout: cfg.Timeout}
	if cfg.Proxy == nil {
		return forward.Dial("tcp", address)
	}

	d, err := internal.FromURL(cfg.Proxy, forward)
	if err != nil {
		return nil, errors.New("error connecting to proxy: " + err.Error())
	}
	return d.Dial("tcp", address)
}

// bind authenticates the session using the configured method.
func (c *Client) bind() error {
	switch {
	case c.cfg.Kerberos:
		return c.kerberosBind()

	case c.cfg.NtlmHash != "":
		if c.cfg.Username == "" {
			return errors.New("a username is required for NTLM hash authentication")
		}
		c.log.Debug("binding with NTLM hash (pass-the-hash)", "domain", c.cfg.Domain, "user", c.cfg.Username)
		return c.conn.NTLMBindWithHash(c.cfg.Domain, c.cfg.Username, ntHash(c.cfg.NtlmHash))

	case c.cfg.Username != "":
		if c.cfg.Domain != "" {
			// NTLM bind works over plain 389 without TLS and accepts either a
			// NetBIOS or FQDN domain, which is the most forgiving default.
			c.log.Debug("binding with NTLM password", "domain", c.cfg.Domain, "user", c.cfg.Username)
			return c.conn.NTLMBind(c.cfg.Domain, c.cfg.Username, c.cfg.Password)
		}
		c.log.Debug("binding with LDAP simple bind", "user", c.cfg.Username)
		return c.conn.Bind(c.cfg.Username, c.cfg.Password)

	default:
		c.log.Debug("binding anonymously")
		return c.conn.UnauthenticatedBind("")
	}
}

// kerberosBind authenticates using a Kerberos ticket cache (ccache) via
// GSSAPI/SASL.
func (c *Client) kerberosBind() error {
	ccache := c.cfg.CCachePath
	if ccache == "" {
		ccache = strings.TrimPrefix(os.Getenv("KRB5CCNAME"), "FILE:")
	}
	if ccache == "" {
		return errors.New("kerberos: no ccache provided (use --krb5-ccache or set KRB5CCNAME)")
	}

	krb5conf := c.cfg.Krb5Conf
	if krb5conf == "" {
		krb5conf = "/etc/krb5.conf"
	}

	gc, err := gssapi.NewClientFromCCache(ccache, krb5conf)
	if err != nil {
		return errors.New("kerberos: loading ccache failed: " + err.Error())
	}
	c.gss = gc

	spn := c.cfg.SPN
	if spn == "" {
		// The SPN must be the DC's FQDN; a bare IP in --dc will not work for
		// Kerberos, so point the user at --spn when it does.
		spn = "ldap/" + c.cfg.Host
	}

	c.log.Debug("binding with Kerberos", "ccache", ccache, "spn", spn)
	if err := c.conn.GSSAPIBind(gc, spn, ""); err != nil {
		return errors.New("kerberos bind failed (spn=" + spn + "): " + err.Error())
	}
	return nil
}

// Dump enumerates every DNS zone reachable in AD and returns the decoded
// records together with the sorted list of zone names that were found.
//
// It first looks in the standard ADIDNS locations (the classic MicrosoftDNS
// containers and the DNS partitions the DC advertises). Only when those yield
// nothing does it fall back to a more extensive search: the forest Partitions
// container, the domain partition subtree, and child probing for diagnostics.
func (c *Client) Dump() ([]string, []*ZoneRecord, error) {
	// A user-forced base short-circuits the tiered discovery.
	if c.cfg.BaseDN != "" {
		zones, records := c.enumerate([]searchBase{{dn: c.cfg.BaseDN, scope: ldap.ScopeWholeSubtree}}, true)
		return zones, records, nil
	}

	root, err := c.readRootDSE()
	if err != nil {
		return nil, nil, err
	}

	// Tier 1: the standard locations.
	zones, records := c.enumerate(c.standardBases(root), false)
	if len(zones) > 0 {
		return zones, records, nil
	}

	// Tier 2: extensive fallback.
	c.log.Debug("standard DNS locations returned nothing; running extensive search")
	extBases, dnsPartitions := c.extensiveBases(root)
	zones, records = c.enumerate(extBases, true)

	if len(zones) == 0 {
		if len(dnsPartitions) == 0 {
			c.log.Warn("this domain has no AD-integrated DNS partitions "+
				"(DomainDnsZones / ForestDnsZones were not found) — DNS is served "+
				"outside Active Directory (e.g. Infoblox/BIND), so there is nothing "+
				"to dump here",
				"partitions_advertised", root.namingContexts)
		} else {
			c.log.Warn("AD-integrated DNS partitions exist but no zones were readable; "+
				"the bind account most likely lacks read access to them",
				"dns_partitions", dnsPartitions)
		}
	}

	return zones, records, nil
}

// enumerate searches each base for dnsZone objects and dumps the records of
// every zone found. When probe is set, a base that yields no zone is probed for
// its children (a diagnostic reserved for the extensive fallback).
func (c *Client) enumerate(bases []searchBase, probe bool) ([]string, []*ZoneRecord) {
	tried := make([]string, 0, len(bases))
	for _, b := range bases {
		tried = append(tried, b.dn)
	}
	c.log.Debug("DNS search bases to enumerate", "bases", tried)

	zoneSet := map[string]bool{}
	zones := []string{}
	records := []*ZoneRecord{}

	for _, base := range bases {
		zreq := ldap.NewSearchRequest(
			base.dn, base.scope, ldap.NeverDerefAliases, 0, 0, false,
			"(objectClass=dnsZone)", []string{"name"}, nil,
		)
		zres, err := c.conn.SearchWithPaging(zreq, pageSize)
		if err != nil {
			// A base an account cannot list returns NoSuchObject (code 32) in AD;
			// keep that at debug (several bases are fallbacks that may not exist)
			// but surface any other failure (referral, size limit, access).
			if ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject) {
				c.log.Debug("search base not present or not listable", "base", base.dn)
			} else {
				c.log.Warn("zone listing failed", "base", base.dn, "err", err)
			}
			continue
		}

		c.log.Debug("zones found under base", "base", base.dn, "count", len(zres.Entries))
		if len(zres.Entries) == 0 {
			if probe {
				c.probeChildren(base.dn)
			}
			continue
		}

		for _, ze := range zres.Entries {
			zone := normalizeZone(firstNonEmpty(ze.GetEqualFoldAttributeValue("name"), firstRDNValue(ze.DN)))
			if zone == "" || isPseudoZone(zone) {
				continue
			}

			if !zoneSet[zone] {
				zoneSet[zone] = true
				zones = append(zones, zone)
			}

			recs, err := c.dumpZoneNodes(ze.DN, zone)
			if err != nil {
				c.log.Warn("could not read records for zone", "zone", zone, "err", err)
				continue
			}
			c.log.Debug("records dumped for zone", "zone", zone, "count", len(recs))
			records = append(records, recs...)
		}
	}

	sort.Strings(zones)
	c.log.Debug("enumeration complete", "zones", len(zones), "nodes", len(records))
	return zones, records
}

// probeChildren is a diagnostic: when a container holds no dnsZone object it
// lists whatever children it does have (with their objectClass) so the user can
// tell an empty/unreadable container from one with an unexpected layout.
func (c *Client) probeChildren(container string) {
	req := ldap.NewSearchRequest(
		container, ldap.ScopeSingleLevel, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)", []string{"objectClass"}, nil,
	)
	res, err := c.conn.SearchWithPaging(req, pageSize)
	if err != nil {
		c.log.Warn("container holds no dnsZone and could not be listed",
			"container", container, "err", err)
		return
	}

	c.log.Warn("container holds no dnsZone entries", "container", container, "children", len(res.Entries))
	for i, e := range res.Entries {
		if i >= 8 {
			c.log.Warn("  (more children omitted)")
			break
		}
		c.log.Warn("  child", "dn", e.DN, "objectClass", e.GetAttributeValues("objectClass"))
	}
}

// dumpZoneNodes reads every dnsNode under a zone and decodes its dnsRecord
// values. Tombstoned (deleted) nodes are skipped.
func (c *Client) dumpZoneNodes(zoneDN, zone string) ([]*ZoneRecord, error) {
	req := ldap.NewSearchRequest(
		zoneDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=dnsNode)", []string{"name", "dnsRecord", "dNSTombstoned"}, nil,
	)

	res, err := c.conn.SearchWithPaging(req, pageSize)
	if err != nil {
		// Some DCs reject the simple paged control combination; fall back to an
		// unpaged search before giving up on the zone.
		c.log.Debug("paged search failed, retrying unpaged", "zone", zone, "err", err)
		res, err = c.conn.Search(req)
		if err != nil {
			return nil, err
		}
	}

	out := []*ZoneRecord{}
	for _, e := range res.Entries {
		if strings.EqualFold(e.GetEqualFoldAttributeValue("dNSTombstoned"), "TRUE") {
			continue
		}

		name := firstNonEmpty(e.GetEqualFoldAttributeValue("name"), firstRDNValue(e.DN))
		zr := &ZoneRecord{
			Zone: zone,
			Name: name,
			FQDN: buildFQDN(name, zone),
		}

		for _, raw := range e.GetEqualFoldRawAttributeValues("dnsRecord") {
			rec, err := ParseDNSRecord(raw)
			if err != nil {
				c.log.Debug("dnsRecord parse failed", "fqdn", zr.FQDN, "err", err)
				continue
			}
			if rec == nil || rec.Type == typeZERO {
				continue
			}
			zr.Records = append(zr.Records, rec)
		}

		if len(zr.Records) > 0 {
			out = append(out, zr)
		}
	}

	return out, nil
}

// rootDSE holds the naming-context pointers read from the server's RootDSE.
type rootDSE struct {
	defaultNC      string
	rootNC         string
	configNC       string
	namingContexts []string
}

// readRootDSE reads the naming-context attributes from the RootDSE.
func (c *Client) readRootDSE() (*rootDSE, error) {
	req := ldap.NewSearchRequest(
		"", ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)",
		[]string{"defaultNamingContext", "rootDomainNamingContext", "configurationNamingContext", "namingContexts"}, nil,
	)
	res, err := c.conn.Search(req)
	if err != nil {
		return nil, errors.New("RootDSE query failed: " + err.Error())
	}
	if len(res.Entries) == 0 {
		return nil, errors.New("RootDSE query returned no entry")
	}

	e := res.Entries[0]
	root := &rootDSE{
		defaultNC:      e.GetEqualFoldAttributeValue("defaultNamingContext"),
		rootNC:         e.GetEqualFoldAttributeValue("rootDomainNamingContext"),
		configNC:       e.GetEqualFoldAttributeValue("configurationNamingContext"),
		namingContexts: e.GetEqualFoldAttributeValues("namingContexts"),
	}

	c.log.Debug("RootDSE",
		"defaultNamingContext", root.defaultNC,
		"rootDomainNamingContext", root.rootNC,
		"configurationNamingContext", root.configNC,
		"namingContexts", root.namingContexts)

	return root, nil
}

// standardBases returns the standard ADIDNS locations: the classic
// MicrosoftDNS containers plus any DNS partition the DC itself advertises.
func (c *Client) standardBases(root *rootDSE) []searchBase {
	bases := []searchBase{}
	seen := map[string]bool{}
	add := func(dn string, scope int) {
		if dn == "" {
			return
		}
		k := strings.ToLower(dn)
		if seen[k] {
			return
		}
		seen[k] = true
		bases = append(bases, searchBase{dn: dn, scope: scope})
	}

	// DNS partitions the DC advertises, searched as a subtree so we do not
	// depend on the exact CN=MicrosoftDNS child RDN.
	for _, nc := range root.namingContexts {
		if isDNSPartition(nc) {
			add(nc, ldap.ScopeWholeSubtree)
		}
	}

	// The classic MicrosoftDNS containers (single level).
	if root.defaultNC != "" {
		add("CN=MicrosoftDNS,DC=DomainDnsZones,"+root.defaultNC, ldap.ScopeSingleLevel)
		add("CN=MicrosoftDNS,CN=System,"+root.defaultNC, ldap.ScopeSingleLevel)
	}
	if root.rootNC != "" {
		add("CN=MicrosoftDNS,DC=ForestDnsZones,"+root.rootNC, ldap.ScopeSingleLevel)
	}

	return bases
}

// extensiveBases returns the fallback locations used only when the standard
// search finds nothing: DNS partitions defined forest-wide (the Partitions
// container) and, as a last resort, the whole domain partition. It also returns
// the set of DNS application partitions discovered (for diagnostics).
func (c *Client) extensiveBases(root *rootDSE) ([]searchBase, []string) {
	// Union of DNS partitions the DC hosts and those defined forest-wide. These
	// objects are readable regardless of ACLs, so this reliably tells an
	// AD-integrated-DNS domain apart from one served by Infoblox/BIND.
	dnsPartitions := []string{}
	dnsSet := map[string]bool{}
	addDNS := func(nc string) {
		if !isDNSPartition(nc) {
			return
		}
		k := strings.ToLower(nc)
		if !dnsSet[k] {
			dnsSet[k] = true
			dnsPartitions = append(dnsPartitions, nc)
		}
	}
	for _, nc := range root.namingContexts {
		addDNS(nc)
	}
	for _, nc := range c.partitionNCs(root.configNC) {
		addDNS(nc)
	}

	bases := []searchBase{}
	seen := map[string]bool{}
	add := func(dn string, scope int) {
		if dn == "" {
			return
		}
		k := strings.ToLower(dn)
		if seen[k] {
			return
		}
		seen[k] = true
		bases = append(bases, searchBase{dn: dn, scope: scope})
	}

	for _, nc := range dnsPartitions {
		add(nc, ldap.ScopeWholeSubtree)
	}

	// Last resort: the domain partition itself. The dnsZone filter keeps the
	// result small even though the server walks a wider subtree.
	add(root.defaultNC, ldap.ScopeWholeSubtree)

	return bases, dnsPartitions
}

// isDNSPartition reports whether a naming context is an AD DNS application
// partition (DomainDnsZones / ForestDnsZones).
func isDNSPartition(nc string) bool {
	l := strings.ToLower(nc)
	return strings.HasPrefix(l, "dc=domaindnszones,") || strings.HasPrefix(l, "dc=forestdnszones,")
}

// partitionNCs returns the nCName of every crossRef in the forest's Partitions
// container (Configuration NC). This lists all naming contexts defined in the
// forest, including DNS partitions not hosted by the connected DC.
func (c *Client) partitionNCs(configNC string) []string {
	if configNC == "" {
		return nil
	}

	req := ldap.NewSearchRequest(
		"CN=Partitions,"+configNC, ldap.ScopeSingleLevel, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=crossRef)", []string{"nCName"}, nil,
	)
	res, err := c.conn.SearchWithPaging(req, pageSize)
	if err != nil {
		c.log.Debug("partitions enumeration failed", "config", configNC, "err", err)
		return nil
	}

	out := []string{}
	for _, e := range res.Entries {
		if nc := e.GetEqualFoldAttributeValue("nCName"); nc != "" {
			out = append(out, nc)
		}
	}
	return out
}

// ntHash extracts the NT hash from either a bare NT hash or an "LM:NT" pair.
func ntHash(h string) string {
	h = strings.TrimSpace(h)
	if i := strings.LastIndex(h, ":"); i >= 0 {
		return h[i+1:]
	}
	return h
}

// buildFQDN joins a relative record name with its zone apex.
func buildFQDN(name, zone string) string {
	name = strings.Trim(strings.ToLower(name), ". ")
	zone = strings.Trim(strings.ToLower(zone), ". ")
	if name == "" || name == "@" || name == zone {
		return zone
	}
	if strings.HasSuffix(name, "."+zone) {
		return name
	}
	return name + "." + zone
}

// firstRDNValue returns the value of the first (leftmost) RDN of a DN, e.g.
// "DC=www,DC=example.com,..." -> "www".
func firstRDNValue(dn string) string {
	first := strings.SplitN(dn, ",", 2)[0]
	if i := strings.Index(first, "="); i >= 0 {
		return first[i+1:]
	}
	return first
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func normalizeZone(zone string) string {
	return strings.Trim(strings.ToLower(zone), ". ")
}

// isPseudoZone filters the internal helper zones that carry no useful records
// for a recon dump.
func isPseudoZone(zone string) bool {
	switch zone {
	case "rootdnsservers", "..trustanchors":
		return true
	}
	return strings.HasPrefix(zone, "..")
}
