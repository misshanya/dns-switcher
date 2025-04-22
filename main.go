package main

import (
	"fmt"
	"log"
	"time"

	"github.com/miekg/dns"
)

func main() {
	cfg := NewConfig()
	StartDNSServer(cfg)
}

func StartDNSServer(cfg *Config) {
	handler := NewDNSHandler(cfg.Upstreams)
	server := &dns.Server{
		Addr:      cfg.Address,
		Net:       "udp",
		Handler:   handler,
		UDPSize:   65535,
		ReusePort: true,
	}

	fmt.Printf("Starting DNS server on %s\n", cfg.Address)

	go watchUpstreams(handler, cfg.Upstreams)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %s\n", err)
	}
}

type dnsHandler struct {
	upstreams []string
	Upstream  string
}

func NewDNSHandler(upstreams []string) *dnsHandler {
	return &dnsHandler{
		upstreams: upstreams,
		Upstream:  upstreams[0],
	}
}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true

	for _, question := range r.Question {
		answers, err := resolve(h.Upstream, question.Name, question.Qtype)
		if err != nil {
			fmt.Println("Got err while resolving. Trying to get new upstream...")
			newUpstream := getWorkingUpstream(h.upstreams)
			if newUpstream == "" {
				log.Println("Not found working upstream")
				break
			}
			h.Upstream = newUpstream
			answers, _ = resolve(h.Upstream, question.Name, question.Qtype)
		}

		msg.Answer = append(msg.Answer, answers...)
	}

	if len(msg.Answer) > 0 {
		fmt.Printf("[RESPONSE] %s  [SERVER] %s\n", msg.Answer[0].Header().Name, h.Upstream)
	}

	w.WriteMsg(msg)
}

func resolve(server string, domain string, qtype uint16) ([]dns.RR, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), qtype)
	m.RecursionDesired = true

	c := &dns.Client{Timeout: 5 * time.Second}

	response, _, err := c.Exchange(m, server)
	if err != nil {
		log.Printf("[ERROR]: %v\n", err)
		return nil, err
	}

	if response == nil {
		log.Printf("[NO RESPONSE] from server\n")
		return nil, err
	}

	return response.Answer, nil
}

func getWorkingUpstream(upstreams []string) string {
	c := &dns.Client{Timeout: 5 * time.Second}

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("google.com"), dns.TypeA)
	m.RecursionDesired = true

	for _, upstream := range upstreams {
		log.Printf("Testing %s...\n", upstream)

		// Send test-request and check server accessibility
		_, _, err := c.Exchange(m, upstream)
		if err == nil {
			log.Printf("Found new working upstream: %s...\n", upstream)
			return upstream
		}
	}

	return ""
}

// watchUpstreams runs as background goroutine to monitor and switch upstreams.
// It makes switch to first available server.
func watchUpstreams(handler *dnsHandler, upstreams []string) {
	ticker := time.NewTicker(1 * time.Minute)

	c := &dns.Client{Timeout: 5 * time.Second}

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("google.com"), dns.TypeA)
	m.RecursionDesired = true

	for range ticker.C {
		for _, upstream := range upstreams {
			log.Printf("Testing %s...\n", upstream)
			_, _, err := c.Exchange(m, upstream)
			if err == nil && upstream != handler.Upstream {
				log.Printf("Switching upstream from %s to %s...\n", handler.Upstream, upstream)
				handler.Upstream = upstream
				break
			}
			if err == nil && upstream == handler.Upstream {
				log.Printf("Current upstream works right")
				break
			}
		}
	}
}
