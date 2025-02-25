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

	//"github.com/helviojunior/enumdns/internal/islazy"
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
}

type Status struct {
	Total int
	Complete int
	Error int
	Label string
	Running bool
}

func (st *Status) Print() { 
	switch st.Label {
		case "[=====]":
            st.Label = "[ ====]"
        case  "[ ====]":
            st.Label = "[  ===]"
        case  "[  ===]":
            st.Label = "[=  ==]"
        case "[=  ==]":
            st.Label = "[==  =]"
        case  "[==  =]":
            st.Label = "[===  ]"
        case "[===  ]":
            st.Label = "[==== ]"
        default:
            st.Label = "[=====]"
	}

	fmt.Fprintf(os.Stderr, "%s\n    %s (%d/%d) failed: %d               \r\033[A", 
    	"                                                                        ",
    	st.Label, st.Complete, st.Total, st.Error)
	
} 

func (run *Runner) AddSkiped() { 
    run.status.Complete += 1
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
		Targets:      make(chan string),
		uid: fmt.Sprintf("%d", time.Now().UnixMilli()),
		ctx:        ctx,
		cancel:     cancel,
		log:        logger,
		writers:    writers,
		options:    opts,
		searchOrder: []uint16{ dns.TypeCNAME, dns.TypeA, dns.TypeAAAA, dns.TypeANY },
		dnsServer: opts.DnsServer + ":" + fmt.Sprintf("%d", opts.DnsPort),
		status:     &Status{
			Total: 0,
			Complete: 0,
			Error: 0,
			Label: "[=====]",
			Running: true,
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

// runWriters takes a result and passes it to writers
func (run *Runner) runWriters(result *models.Result) error {
	for _, writer := range run.writers {
		if err := writer.Write(result); err != nil {
			return err
		}
	}

	return nil
}

func (run *Runner) Run(total int) Status {
	wg := sync.WaitGroup{}
	swg := sync.WaitGroup{}

	run.status.Total = total

    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        fmt.Fprintf(os.Stderr, "\n%s\n", 
            "                                                                                ",
        )
        run.log.Warn("interrupted, shutting down...                            \n")

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
			        	time.Sleep(time.Duration(time.Second/4))
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
			for run.status.Running {
				select {
				case <-run.ctx.Done():
					return
				case host, ok := <-run.Targets:
					if !ok || !run.status.Running {
						return
					}
					logger := run.log.With("FQDN", host)

					results := run.Probe(host)
					if run.status.Running {
						
						if len(results) == 0 {
							//Host not found
							run.status.Complete += 1
						}else{
							run.status.AddResult(results[0])
						}
						
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
	run.status.Running = false
	swg.Wait()

    fmt.Fprintf(os.Stderr, "\n%s\n%s\r", 
        "                                                                                ",
        "                                                                                ",
    )

	return *run.status
}

func (run *Runner) Probe(host string) []*models.Result {
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
		for good_to_go != true && run.status.Running {
            m := new(dns.Msg)
            m.Id = dns.Id()
			m.RecursionDesired = true

			//m.Question = make([]dns.Question, 1)
			//m.Question[0] = dns.Question{host, t, dns.ClassINET}
			m.SetQuestion(host, t)

			r, err := dns.Exchange(m, run.dnsServer); 
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
						}
					}

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

			r, err := dns.Exchange(m, run.dnsServer); 
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

func (run *Runner) Close() {
	for _, writer := range run.writers {
		writer.Finish()
	}
}

