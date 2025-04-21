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
	handler := new(dnsHandler)
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
	Upstream string
}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true

	for _, question := range r.Question {
		answers := resolve(h.Upstream, question.Name, question.Qtype)
		msg.Answer = append(msg.Answer, answers...)
	}

	w.WriteMsg(msg)
}

func resolve(server string, domain string, qtype uint16) []dns.RR {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), qtype)
	m.RecursionDesired = true

	c := &dns.Client{Timeout: 5 * time.Second}

	response, _, err := c.Exchange(m, server)
	if err != nil {
		log.Printf("[ERROR]: %v\n", err)
		return nil
	}

	if response == nil {
		log.Printf("[NO REPONSE] from server\n")
		return nil
	}

	return response.Answer
}

func watchUpstreams(handler *dnsHandler, upstreams []string) {
	ticker := time.NewTicker(10 * time.Second)

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
			}
		}
	}
}
