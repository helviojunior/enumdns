package runner

import (
	"context"
	//"errors"
	"log/slog"
	//"net/url"
	//"net/mail"
	"os"
	"fmt"
	"sync"
	"time"
	"strings"
	"math/rand/v2"
	"os/signal"
    "syscall"
    //"strconv"
    "errors"

	"github.com/helviojunior/enumdns/internal"
	//"github.com/helviojunior/enumdns/internal/ascii"
	//"github.com/helviojunior/enumdns/internal/tools"
	"github.com/helviojunior/enumdns/pkg/models"
	"github.com/helviojunior/enumdns/pkg/writers"
	"github.com/miekg/dns"
)

// Runner is a runner that probes web targets using a driver
type Recon struct {
	
	//Test id
	uid string

	// DNS FQDN to scan.
	Targets chan string

	//Context
	ctx    context.Context
	cancel context.CancelFunc

	// writers are the result writers to use
	writers []writers.Writer

	// log handler
	log *slog.Logger

	// options for the Recon to consider
	options Options

	//search order
	searchOrder []uint16

	//DNS Server
	dnsServer string

	//Running
	Running bool
}

// New gets a new Recon ready for probing.
// It's up to the caller to call Close() on the runner
func NewRecon(logger *slog.Logger, opts Options, writers []writers.Writer) (*Recon, error) {
	ctx, cancel := context.WithCancel(context.Background())

	return &Recon{
		Targets:      make(chan string),
		uid: fmt.Sprintf("%d", time.Now().UnixMilli()),
		ctx:        ctx,
		cancel:     cancel,
		log:        logger,
		writers:    writers,
		options:    opts,
		searchOrder: []uint16{ dns.TypeSOA, dns.TypeCNAME, dns.TypeA, dns.TypeAAAA, dns.TypeANY, dns.TypeMX, dns.TypeTXT, dns.TypeNS, dns.TypeSRV },
		dnsServer: opts.DnsServer + ":" + fmt.Sprintf("%d", opts.DnsPort),
		Running: true,
	}, nil
}

// runWriters takes a result and passes it to writers
func (run *Recon) runWriters(result *models.Result) error {
	for _, writer := range run.writers {
		if err := writer.Write(result); err != nil {
			return err
		}
	}

	return nil
}

func (run *Recon) Run(total int) {
	wg := sync.WaitGroup{}

    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        run.Running = false
    }()

	// will spawn Scan.Theads number of "workers" as goroutines
	for w := 0; w < run.options.Scan.Threads; w++ {
		wg.Add(1)

		// start a worker
		go func() {
			defer wg.Done()
			for run.Running {
				select {
				case <-run.ctx.Done():
					return
				case host, ok := <-run.Targets:
					if !ok || !run.Running {
						return
					}
					logger := run.log.With("FQDN", host)

					logger.Debug("Getting SOA")
					soa, err := run.GetSOAName(host)
					if err != nil {
						logger.Error("failed to get SOA", "err", err)
						continue
					}

					logger.Debug("SOA result", "soa", soa)
					if soa != host {
						logger = run.log.With("FQDN", soa)
					}

					//Check if has LDAP registers
					//nslookup -q=SRV _ldap._tcp.sec4us.com.br
					results := run.Probe("_ldap._tcp." + soa)
					if run.Running {
						for _, res := range results {
        					
        					if err := run.runWriters(res); err != nil {
        						logger.Error("failed to write result", "err", err)
        					}
                        }
                    }

					results = run.Probe(soa)
					if run.Running {
						for _, res := range results {
        					
        					if err := run.runWriters(res); err != nil {
        						logger.Error("failed to write result", "err", err)
        					}
                        }
                    }
				}
			}

		}()
	}

	wg.Wait()
	run.Running = false

    fmt.Fprintf(os.Stderr, "\n%s\n%s\r", 
        "                                                                                ",
        "                                                                                ",
    )

	return
}

func (run *Recon) GetSOAName(host string) (string, error) {
	s := ""
	host = strings.Trim(host, ". ")
	if host == "" {
		return "", errors.New("empty host string")
	}

	host = strings.ToLower(host) + "."

    m := new(dns.Msg)
    m.Id = dns.Id()
	m.RecursionDesired = true

	m.Question = make([]dns.Question, 1)
	m.Question[0] = dns.Question{host, dns.TypeSOA, dns.ClassINET}

	c := new(internal.SocksClient)
	in, err := c.Exchange(m, run.options.Proxy, run.dnsServer); 
	if err != nil {
		return "", err
	}else{
		
		for _, ans1 := range in.Answer {
			if soa, ok := ans1.(*dns.SOA); ok {
				s = strings.Trim(soa.Hdr.Name, ". ")
			}
		}
		
	}

	if s == "" {
		return "", errors.New("SOA not found for domain '"+ host + "'")
	}

	return s, nil

}

func (run *Recon) Probe(host string) []*models.Result {
	host = strings.Trim(host, ". ")
	host = strings.ToLower(host) + "."

	logger := run.log.With("FQDN", host)
	resList := []*models.Result{}

	resultBase := &models.Result{
		TestId: run.uid,
		FQDN: host,
		ProbedAt: time.Now(),
		Exists: true,
	}
    
	ips := []string{}

    for _, t := range run.searchOrder {
    	tName := dns.Type(t).String()
    	logger := run.log.With("FQDN", host, "Type", tName)
    	resultBase.RType = tName

    	good_to_go := false
		counter := 0
		for good_to_go != true && run.Running {
            m := new(dns.Msg)
            m.Id = dns.Id()
			m.RecursionDesired = true

			//m.Question = make([]dns.Question, 1)
			//m.Question[0] = dns.Question{host, t, dns.ClassINET}
			m.SetQuestion(host, t)

			//r, err := dns.Exchange(m, run.dnsServer); 
			c := new(internal.SocksClient)
			r, err := c.Exchange(m, run.options.Proxy, run.dnsServer); 
			counter += 1
			good_to_go = (err == nil)

			if err != nil {
				logger.Debug("Error running DNS request, trying again...", "type", t, "err", err)
				time.Sleep(time.Duration(rand.IntN(20)) * time.Second)
			}

			if !good_to_go && counter >= 5 {
				resultBase.Exists = false
				resultBase.Failed = true
				resultBase.FailedReason = err.Error()
				return []*models.Result{ resultBase }
			}
			
			if good_to_go {
				for _, ans := range r.Answer {
					run.log.Debug(ans.String())

					//SOA
					soa, ok := ans.(*dns.SOA)
					if ok {
						logger.Debug("SOA", "Name", soa.Hdr.Name)
						c1 := resultBase.Clone()
						c1.RType = "SOA"
						c1.Target = soa.Hdr.Name
						if !models.SliceHasResult(resList, c1) {
							cc, prodName, _ := ContainsCloudProduct(soa.Hdr.Name)
							if cc {
								c1.CloudProduct = prodName
							}
							resList = append(resList, c1)

							run.getHost(c, soa.Hdr.Name, &resList, &ips, resultBase)
						}
					}
					
					//TXT
					txt, ok := ans.(*dns.TXT)
					if ok {
						for _, t := range txt.Txt {
							logger.Debug("TXT", "Value", t)
							c1 := resultBase.Clone()
							c1.RType = "TXT"
							c1.Txt = t
							if !models.SliceHasResult(resList, c1) {
								resList = append(resList, c1)
							}
						}
					}

					//CNAME
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
							resList = append(resList, c1)

							run.getHost(c, cname.Target, &resList, &ips, resultBase)
						}
					}

					//SRV
					srv, ok := ans.(*dns.SRV)
					if ok {
						logger.Debug("SRV", "Target", srv.Target)
						c1 := resultBase.Clone()
						c1.RType = "SRV"
						c1.Target = srv.Target
						if !models.SliceHasResult(resList, c1) {
							cc, prodName, _ := ContainsCloudProduct(srv.Target)
							if cc {
								c1.CloudProduct = prodName
							}
							resList = append(resList, c1)

							run.getHost(c, srv.Target, &resList, &ips, resultBase)
						}
					}

					//MX
					mx, ok := ans.(*dns.MX)
					if ok {
						logger.Debug("MX", "Target", mx.Mx)
						c1 := resultBase.Clone()
						c1.RType = "MX"
						c1.Target = mx.Mx
						if !models.SliceHasResult(resList, c1) {
							cc, prodName, _ := ContainsCloudProduct(mx.Mx)
							if cc {
								c1.CloudProduct = prodName
							}
							resList = append(resList, c1)

							run.getHost(c, mx.Mx, &resList, &ips, resultBase)
						}
					}

					//NS
					ns, ok := ans.(*dns.NS)
					if ok {
						logger.Debug("NS", "Target", ns.Ns)
						c1 := resultBase.Clone()
						c1.RType = "NS"
						c1.Target = ns.Ns
						if !models.SliceHasResult(resList, c1) {
							cc, prodName, _ := ContainsCloudProduct(ns.Ns)
							if cc {
								c1.CloudProduct = prodName
							}
							resList = append(resList, c1)

							run.getHost(c, ns.Ns, &resList, &ips, resultBase)
						}
					}

					//IPv4
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

					//IPv6
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

	for _, ip := range ips {

		if arpa, err := dns.ReverseAddr(ip); err == nil {

			m := new(dns.Msg)
            m.Id = dns.Id()
			m.RecursionDesired = true

			m.SetQuestion(arpa, dns.TypePTR)

			//r, err := dns.Exchange(m, run.dnsServer); 
			c := new(internal.SocksClient)
			r, err := c.Exchange(m, run.options.Proxy, run.dnsServer); 
			if err != nil {
				logger.Error("Error", "err", err)
			}else{
				for _, ans := range r.Answer {
					ptr, ok := ans.(*dns.PTR)
					if ok {
						logger.Debug("PTR", "PTR", arpa, "CNAME", ptr.Ptr)
						a2 := resultBase.Clone()
						a2.FQDN = ptr.Ptr
						a2.RType = "PTR"
						if strings.Contains(arpa, "ip6.arpa") {
							a2.IPv6 = ip
						}else{
							a2.IPv4 = ip
						}
						a2.Ptr = arpa
						if !models.SliceHasResult(resList, a2) {
							resList = append(resList, a2)
						}

						cc, prodName, _ := ContainsCloudProduct(ptr.Ptr)
						for _, res := range resList {
							if res.RType != "PTR" && (res.IPv4 == ip || res.IPv6 == ip) {
								res.Ptr = ptr.Ptr
								if cc {
									res.CloudProduct = prodName
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
		return []*models.Result{ resultBase }
	}

	return resList
}

func (run *Recon) getHost(c *internal.SocksClient, host string, resList *[]*models.Result, ips *[]string, resultBase *models.Result) {
	logger := run.log.With("FQDN", host)

	m1 := new(dns.Msg)
    m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.SetQuestion(strings.Trim(host, ". ") + ".", dns.TypeANY)
	//r1, err := dns.Exchange(m1, run.dnsServer); 
	r1, err := c.Exchange(m1, run.options.Proxy, run.dnsServer); 
	if err != nil {
		return
	}else{
		for _, ans1 := range r1.Answer {
			a, ok := ans1.(*dns.A)
			if ok {
				logger.Debug("A", "IP", a.A.String())
				*ips = append(*ips, a.A.String())

				// With the same FQDN
				a1 := resultBase.Clone()
				a1.RType = "A"
				a1.IPv4 = a.A.String()
				if !models.SliceHasResult(*resList, a1) {
					*resList = append(*resList, a1)
				}

				// With CNAME fqdn
				a1 = resultBase.Clone()
				a1.FQDN = host
				a1.RType = "A"
				a1.IPv4 = a.A.String()
				if !models.SliceHasResult(*resList, a1) {
					*resList = append(*resList, a1)
				}
			}

			aaaa, ok := ans1.(*dns.AAAA)
			if ok {
				logger.Debug("AAAA", "IP", aaaa.AAAA.String())
				*ips = append(*ips, aaaa.AAAA.String())
				
				// With the same FQDN
				a2 := resultBase.Clone()
				a2.RType = "AAAA"
				a2.IPv6 = aaaa.AAAA.String()
				if !models.SliceHasResult(*resList, a2) {
					*resList = append(*resList, a2)
				}

				// With CNAME fqdn
				a2 = resultBase.Clone()
				a2.FQDN = host
				a2.RType = "AAAA"
				a2.IPv6 = aaaa.AAAA.String()
				if !models.SliceHasResult(*resList, a2) {
					*resList = append(*resList, a2)
				}
			}
		}
	}
}

func (run *Recon) Close() {
	for _, writer := range run.writers {
		writer.Finish()
	}
}

