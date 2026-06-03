package runner

import (
	"context"
	//"errors"
	"log/slog"
	//"net/url"
	//"net/mail"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"time"

	//"strconv"

	"golang.org/x/term"

	"github.com/helviojunior/enumdns/internal"
	"github.com/helviojunior/enumdns/internal/ascii"
	"github.com/helviojunior/enumdns/internal/tools"
	"github.com/helviojunior/enumdns/pkg/models"
	"github.com/helviojunior/enumdns/pkg/writers"
	"github.com/miekg/dns"
)

// Runner is a runner that probes web targets using a driver
type Runner struct {

	//Test id
	uid string

	// DNS FQDN to scan.
	Targets chan string

	//Status
	status *Status

	//Context
	ctx    context.Context
	cancel context.CancelFunc

	// writers are the result writers to use
	writers []writers.Writer

	// log handler
	log *slog.Logger

	// options for the Runner to consider
	options Options

	//search order
	searchOrder []uint16

	//DNS Server
	dnsServer string

	//In-memory SOA cache keyed by zone apex. Avoids re-querying the SOA for
	//every host of an already known zone (shared logic in soaCache).
	soa *soaCache
}

type Status struct {
	Total      int
	Complete   int
	Skiped     int
	Error      int
	Spin       string
	Running    bool
	IsTerminal bool
	log        *slog.Logger
}

func (st *Status) Print() {

	if st.IsTerminal {

		st.Spin = ascii.GetNextSpinner(st.Spin)

		fmt.Fprintf(os.Stderr, "%s\n    %s (%s/%s) failed: %s               \r\033[A",
			"                                                                        ",
			ascii.ColoredSpin(st.Spin), tools.FormatInt(st.Complete), tools.FormatInt(st.Total), tools.FormatInt(st.Error))

	} else {
		st.log.Info("STATUS", "Total", tools.FormatInt(st.Total), "Complete", tools.FormatInt(st.Complete), "Errors", tools.FormatInt(st.Error))
	}
}

func (run *Runner) GetLog() *slog.Logger {
	return run.log
}

func (run *Runner) AddSkiped() {
	run.status.Complete += 1
	run.status.Skiped += 1
}

func (st *Status) AddResult(result *models.Result) {
	st.Complete += 1
	if result.Failed {
		st.Error += 1
		return
	}
}

// New gets a new Runner ready for probing.
// It's up to the caller to call Close() on the runner
func NewRunner(logger *slog.Logger, opts Options, writers []writers.Writer) (*Runner, error) {
	ctx, cancel := context.WithCancel(context.Background())

	return &Runner{
		Targets:     make(chan string),
		uid:         fmt.Sprintf("%d", time.Now().UnixMilli()),
		ctx:         ctx,
		cancel:      cancel,
		log:         logger,
		writers:     writers,
		options:     opts,
		searchOrder: []uint16{dns.TypeCNAME, dns.TypeA, dns.TypeAAAA, dns.TypeANY},
		dnsServer:   opts.DnsServer + ":" + fmt.Sprintf("%d", opts.DnsPort),
		soa:         newSOACache(),
		status: &Status{
			Total:      0,
			Complete:   0,
			Error:      0,
			Skiped:     0,
			Spin:       "",
			Running:    true,
			IsTerminal: term.IsTerminal(int(os.Stdin.Fd())),
			log:        logger,
		},
	}, nil
}

func ContainsCloudProduct(s string) (bool, string, string) {
	s = strings.Trim(strings.ToLower(s), ". ")
	for prodName, identifiers := range products {
		for _, id := range identifiers {
			if strings.Contains(s, strings.ToLower(id)) {
				return true, prodName, id
			}
		}
	}
	return false, "", ""
}

func ContainsSaaS(s string) (bool, string, string) {
	s = strings.Trim(strings.ToLower(s), ". ")
	for prodName, identifiers := range saas_products {
		for _, id := range identifiers {
			if strings.Contains(s, strings.ToLower(id)) {
				return true, prodName, id
			}
		}
	}
	return false, "", ""
}

func ContainsDatacenter(s string) (bool, string, string) {
	s = strings.Trim(strings.ToLower(s), ". ")
	for prodName, identifiers := range datacenter {
		for _, id := range identifiers {
			if strings.Contains(s, strings.ToLower(id)) {
				return true, prodName, id
			}
		}
	}
	return false, "", ""
}

// runWriters takes a result and passes it to writers
func (run *Runner) runWriters(result *models.Result) error {
	for _, writer := range run.writers {
		if err := writer.Write(result); err != nil {
			run.log.Debug("Error at writer", "type", reflect.TypeOf(writer).Name(), "err", err)
			return err
		}
	}

	return nil
}

// runWritersSOA takes a SOA object and passes it to writers
func (run *Runner) runWritersSOA(soa *models.SOA) error {
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

// soaFromDNS builds a models.SOA from a DNS SOA record, normalizing names and
// enriching the Cloud/SaaS/Datacenter attributes from the primary nameserver.
func soaFromDNS(testId string, soa *dns.SOA) *models.SOA {
	if soa == nil {
		return nil
	}

	s := &models.SOA{
		TestId:    testId,
		Name:      strings.Trim(strings.ToLower(tools.UnescapeDNSName(soa.Hdr.Name)), ". "),
		PrimaryNS: strings.Trim(strings.ToLower(tools.UnescapeDNSName(soa.Ns)), ". "),
		Mbox:      strings.Trim(strings.ToLower(tools.UnescapeDNSName(soa.Mbox)), ". "),
		Serial:    soa.Serial,
		Refresh:   soa.Refresh,
		Retry:     soa.Retry,
		Expire:    soa.Expire,
		MinTTL:    soa.Minttl,
		ProbedAt:  time.Now(),
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

// linkResultToSOA points a resolved record to its zone SOA object, but only when
// the record's FQDN actually belongs to that zone (apex or sub-name). This avoids
// linking unrelated names such as PTR results or external CNAME targets.
func linkResultToSOA(res *models.Result, soa *models.SOA) {
	if res == nil || soa == nil || soa.Name == "" {
		return
	}
	fqdn := strings.Trim(strings.ToLower(res.FQDN), ". ")
	if fqdn == soa.Name || strings.HasSuffix(fqdn, "."+soa.Name) {
		res.SOA = soa.Name
	}
}

// soaCache is a concurrency-safe in-memory cache of SOA objects keyed by zone
// apex (normalized, no trailing dot). It is shared by every runner mode so a
// zone's SOA is queried at most once.
type soaCache struct {
	mu      sync.RWMutex
	entries map[string]*models.SOA
}

func newSOACache() *soaCache {
	return &soaCache{entries: map[string]*models.SOA{}}
}

// lookup walks the host and its parent domains (most specific first) and returns
// the first cached zone that matches. Returns nil on a miss.
func (c *soaCache) lookup(host string) *models.SOA {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.entries) == 0 {
		return nil
	}

	for _, candidate := range parentDomains(host) {
		if soa, ok := c.entries[candidate]; ok {
			return soa
		}
	}
	return nil
}

// store inserts the SOA keyed by its apex and returns the canonical cached entry
// (an existing one wins, keeping a single object per zone under concurrency).
func (c *soaCache) store(soa *models.SOA) *models.SOA {
	if soa == nil || soa.Name == "" {
		return soa
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if existing, ok := c.entries[soa.Name]; ok {
		return existing
	}
	c.entries[soa.Name] = soa
	return soa
}

// resolve returns the SOA object for the host's zone. It checks the cache
// (extracting parent domains) before querying the DNS server via queryFn, and
// caches the answer by zone apex.
func (c *soaCache) resolve(host string, queryFn func(string) *models.SOA) *models.SOA {
	host = strings.Trim(strings.ToLower(host), ". ")
	if host == "" {
		return nil
	}

	if soa := c.lookup(host); soa != nil {
		return soa
	}

	soa := queryFn(host)
	if soa == nil {
		return nil
	}

	return c.store(soa)
}

// parentDomains returns the host itself followed by each of its parent domains,
// from the most specific to the least specific (the TLD alone is not returned).
// e.g. "a.b.example.com" -> ["a.b.example.com", "b.example.com", "example.com"]
func parentDomains(host string) []string {
	host = strings.Trim(strings.ToLower(host), ". ")
	if host == "" {
		return nil
	}

	labels := strings.Split(host, ".")
	out := []string{}
	// Stop before the last label so we never consider a bare TLD as a zone apex.
	for i := 0; i < len(labels)-1; i++ {
		out = append(out, strings.Join(labels[i:], "."))
	}
	return out
}

// resolveSOA returns the SOA object for the zone the host belongs to, using the
// shared cache and the Runner's own DNS query.
func (run *Runner) resolveSOA(host string) *models.SOA {
	return run.soa.resolve(host, run.querySOA)
}

// querySOA performs the actual SOA DNS query for the given host. The SOA record
// may be returned either in the answer section (when the host is the apex) or in
// the authority section (when the host is a sub-name of the zone).
func (run *Runner) querySOA(host string) *models.SOA {
	logger := run.log.With("FQDN", host, "Type", "SOA")
	host = strings.Trim(strings.ToLower(host), ". ") + "."

	counter := 0
	for run.status.Running {
		m := new(dns.Msg)
		m.Id = dns.Id()
		m.RecursionDesired = true
		m.SetQuestion(host, dns.TypeSOA)

		c := new(internal.SocksClient)
		r, err := c.Exchange(m, run.options.Proxy, run.dnsServer)
		counter += 1

		if err != nil {
			logger.Debug("Error running SOA request, trying again...", "err", err)
			if counter >= 5 {
				return nil
			}
			n, _ := rand.Int(rand.Reader, big.NewInt(20))
			time.Sleep(time.Duration(n.Int64()) * time.Second)
			continue
		}

		// SOA may come in the answer or the authority section.
		for _, ans := range append(r.Answer, r.Ns...) {
			if soa, ok := ans.(*dns.SOA); ok {
				logger.Debug("SOA", "Name", soa.Hdr.Name, "MNAME", soa.Ns, "Serial", soa.Serial)
				return soaFromDNS(run.uid, soa)
			}
		}

		// Answered without an error but no SOA was present.
		return nil
	}

	return nil
}

func (run *Runner) Run(total int) Status {
	wg := sync.WaitGroup{}
	swg := sync.WaitGroup{}

	run.status.Total = total

	complementarySearchOrder := []uint16{}
	if run.options.Quick {
		for _, t := range run.searchOrder {
			if t != dns.TypeA {
				complementarySearchOrder = append(complementarySearchOrder, t)
			}
		}
		run.searchOrder = []uint16{dns.TypeA}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		run.status.Running = false
	}()

	if !run.options.Logging.Silence {
		swg.Add(1)
		go func() {
			defer swg.Done()
			for run.status.Running {
				select {
				case <-run.ctx.Done():
					return
				default:
					run.status.Print()
					if run.status.IsTerminal {
						time.Sleep(time.Duration(time.Second / 4))
					} else {
						time.Sleep(time.Duration(time.Second * 10))
					}
				}
			}
		}()
	}

	// will spawn Scan.Theads number of "workers" as goroutines
	for w := 0; w < run.options.Scan.Threads; w++ {
		wg.Add(1)

		// start a worker
		go func() {
			defer wg.Done()
			tools.RandSleep()
			for run.status.Running {
				select {
				case <-run.ctx.Done():
					return
				case host, ok := <-run.Targets:
					if !ok || !run.status.Running {
						return
					}
					logger := run.log.With("FQDN", host)

					results := run.Probe(host, run.searchOrder)
					if run.status.Running {

						// // May be it is not a root domain

						if len(results) == 0 {
							//Host not found
							run.status.Complete += 1
						} else {
							run.status.AddResult(results[0])
						}

						// Always resolve (cached) the SOA of the host's zone and persist
						// it through the registered writers, regardless of whether the
						// host itself resolved. linkResultToSOA below still only links
						// records that actually belong to the zone.
						soa := run.resolveSOA(host)
						if soa != nil {
							if err := run.runWritersSOA(soa); err != nil {
								logger.Error("failed to write SOA", "err", err)
							}
						}

						for _, res := range results {
							linkResultToSOA(res, soa)
							if err := run.runWriters(res); err != nil {
								logger.Error("failed to write result", "err", err)
							}
						}

						if len(results) >= 1 && results[0].Exists && len(complementarySearchOrder) > 0 {
							logger.Debug("Doing complementary search...")
							results := run.Probe(host, complementarySearchOrder)
							for _, res := range results {
								linkResultToSOA(res, soa)
								if err := run.runWriters(res); err != nil {
									logger.Error("failed to write result", "err", err)
								}
							}
						}
					}

					//We must put this to slow down the requests to prevent block from DNS Server
					if !run.options.PrivateDns {
						tools.RandSleep()
					}
				}
			}

		}()
	}

	wg.Wait()
	run.status.Running = false
	swg.Wait()

	return *run.status
}

func (run *Runner) Probe(host string, searchOrder []uint16) []*models.Result {
	logger := run.log.With("FQDN", host)
	resList := []*models.Result{}

	resultBase := &models.Result{
		TestId:   run.uid,
		FQDN:     host,
		ProbedAt: time.Now(),
		Exists:   true,
	}

	ips := []string{}

	for _, t := range searchOrder {
		tName := dns.Type(t).String()
		logger := run.log.With("FQDN", host, "Type", tName)
		resultBase.RType = tName

		good_to_go := false
		counter := 0
		for !good_to_go && run.status.Running {
			m := new(dns.Msg)
			m.Id = dns.Id()
			m.RecursionDesired = true

			//m.Question = make([]dns.Question, 1)
			//m.Question[0] = dns.Question{host, t, dns.ClassINET}
			m.SetQuestion(host, t)

			//r, err := dns.Exchange(m, run.dnsServer);
			c := new(internal.SocksClient)
			r, err := c.Exchange(m, run.options.Proxy, run.dnsServer)
			counter += 1
			good_to_go = (err == nil)

			if err != nil {
				logger.Debug("Error running DNS request, trying again...", "type", t, "err", err)
				// Use crypto/rand for secure random number generation
				n, _ := rand.Int(rand.Reader, big.NewInt(20))
				time.Sleep(time.Duration(n.Int64()) * time.Second)
			}

			if !good_to_go && counter >= 5 {
				resultBase.Exists = false
				resultBase.Failed = true
				resultBase.FailedReason = err.Error()
				return []*models.Result{resultBase}
			}

			if good_to_go {
				for _, ans := range r.Answer {
					cname, ok := ans.(*dns.CNAME)
					if ok {
						logger.Debug("CNAME", "Target", cname.Target)
						c1 := resultBase.Clone()
						c1.RType = "CNAME"
						c1.Target = cname.Target
						if !models.SliceHasResult(resList, c1) {
							cc, prodName, _ := ContainsCloudProduct(cname.Target)
							if cc {
								c1.CloudProduct = prodName
							}
							ss, saasName, _ := ContainsSaaS(cname.Target)
							if ss {
								c1.SaaSProduct = saasName
							}
							dc, dcName, _ := ContainsDatacenter(cname.Target)
							if dc {
								c1.Datacenter = dcName
							}
							resList = append(resList, c1)

							m1 := new(dns.Msg)
							m1.Id = dns.Id()
							m1.RecursionDesired = true
							m1.SetQuestion(strings.Trim(cname.Target, ". ")+".", dns.TypeANY)
							//r1, err := dns.Exchange(m1, run.dnsServer);
							r1, err := c.Exchange(m1, run.options.Proxy, run.dnsServer)
							if err != nil {
								good_to_go = false
							} else {
								for _, ans1 := range r1.Answer {
									a, ok := ans1.(*dns.A)
									if ok {
										logger.Debug("A", "IP", a.A.String())
										ips = append(ips, a.A.String())

										// With the same FQDN
										a1 := resultBase.Clone()
										a1.RType = "A"
										a1.IPv4 = a.A.String()
										a1.CloudProduct = c1.CloudProduct
										a1.SaaSProduct = c1.SaaSProduct
										a1.Datacenter = c1.Datacenter
										if !models.SliceHasResult(resList, a1) {
											resList = append(resList, a1)
										}

										// With CNAME fqdn
										a1 = resultBase.Clone()
										a1.FQDN = cname.Target
										a1.RType = "A"
										a1.IPv4 = a.A.String()
										a1.CloudProduct = c1.CloudProduct
										a1.SaaSProduct = c1.SaaSProduct
										a1.Datacenter = c1.Datacenter
										if !models.SliceHasResult(resList, a1) {
											resList = append(resList, a1)
										}
									}

									aaaa, ok := ans1.(*dns.AAAA)
									if ok {
										logger.Debug("AAAA", "IP", aaaa.AAAA.String())
										ips = append(ips, aaaa.AAAA.String())

										// With the same FQDN
										a2 := resultBase.Clone()
										a2.RType = "AAAA"
										a2.IPv6 = aaaa.AAAA.String()
										a2.CloudProduct = c1.CloudProduct
										a2.SaaSProduct = c1.SaaSProduct
										a2.Datacenter = c1.Datacenter
										if !models.SliceHasResult(resList, a2) {
											resList = append(resList, a2)
										}

										// With CNAME fqdn
										a2 = resultBase.Clone()
										a2.FQDN = cname.Target
										a2.RType = "AAAA"
										a2.IPv6 = aaaa.AAAA.String()
										a2.CloudProduct = c1.CloudProduct
										a2.SaaSProduct = c1.SaaSProduct
										a2.Datacenter = c1.Datacenter
										if !models.SliceHasResult(resList, a2) {
											resList = append(resList, a2)
										}
									}
								}
							}
						}
					}

					if t == dns.TypeA {
						a, ok := ans.(*dns.A)
						if ok {
							logger.Debug("A", "IP", a.A.String())
							ips = append(ips, a.A.String())
							a1 := resultBase.Clone()
							a1.RType = "A"
							a1.IPv4 = a.A.String()
							if !models.SliceHasResult(resList, a1) {
								resList = append(resList, a1)
							}
						}
					}

					if t == dns.TypeAAAA {
						aaaa, ok := ans.(*dns.AAAA)
						if ok {
							logger.Debug("AAAA", "IP", aaaa.AAAA.String())
							ips = append(ips, aaaa.AAAA.String())
							a2 := resultBase.Clone()
							a2.RType = "AAAA"
							a2.IPv6 = aaaa.AAAA.String()
							if !models.SliceHasResult(resList, a2) {
								resList = append(resList, a2)
							}
						}
					}

				}
			}
		}
	}

	for _, ip := range ips {

		if arpa, err := dns.ReverseAddr(ip); err == nil {

			m := new(dns.Msg)
			m.Id = dns.Id()
			m.RecursionDesired = true

			m.SetQuestion(arpa, dns.TypePTR)

			//r, err := dns.Exchange(m, run.dnsServer);
			c := new(internal.SocksClient)
			r, err := c.Exchange(m, run.options.Proxy, run.dnsServer)
			if err != nil {
				logger.Error("Error", "err", err)
			} else {
				for _, ans := range r.Answer {
					ptr, ok := ans.(*dns.PTR)
					if ok {
						logger.Debug("PTR", "PTR", arpa, "CNAME", ptr.Ptr)
						a2 := resultBase.Clone()
						a2.FQDN = ptr.Ptr
						a2.RType = "PTR"
						if strings.Contains(arpa, "ip6.arpa") {
							a2.IPv6 = ip
						} else {
							a2.IPv4 = ip
						}
						a2.Ptr = arpa

						cc, prodName, _ := ContainsCloudProduct(ptr.Ptr)
						ss, saasName, _ := ContainsSaaS(ptr.Ptr)
						dc, dcName, _ := ContainsDatacenter(ptr.Ptr)

						if cc {
							a2.CloudProduct = prodName
						}
						if ss {
							a2.SaaSProduct = saasName
						}
						if dc {
							a2.Datacenter = dcName
						}

						if !models.SliceHasResult(resList, a2) {
							resList = append(resList, a2)
						}

						for _, res := range resList {
							if res.RType != "PTR" && (res.IPv4 == ip || res.IPv6 == ip) {
								res.Ptr = ptr.Ptr
								if cc {
									res.CloudProduct = prodName
								}
								if ss {
									res.SaaSProduct = saasName
								}
								if dc {
									res.Datacenter = dcName
								}
							}
						}
					}
				}
			}
		}
	}

	if len(resList) == 0 {
		resultBase.Exists = false
		return []*models.Result{resultBase}
	}

	return resList
}

func (run *Runner) Close() {
	for _, writer := range run.writers {
		if err := writer.Finish(); err != nil {
			run.log.Error("failed to finish writer", "err", err)
		}
	}
}
