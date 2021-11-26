package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/miekg/dns"
)

var records = map[string]dns.RR{}

func CNameFlatenning(r dns.RR) ([]dns.RR, error) {
	config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
	c := new(dns.Client)

	m := new(dns.Msg)
	m.SetQuestion(r.Header().Name, dns.TypeA)
	m.RecursionDesired = true

	response, _, err := c.Exchange(m, net.JoinHostPort(config.Servers[0], config.Port))
	if r == nil {
		return nil, err
	}

	if response.Rcode != dns.RcodeSuccess {
		log.Fatalf(" *** invalid answer name %s after A query for %s\n", os.Args[1], os.Args[1])
	}
	// Stuff must be in the answer section
	for _, a := range response.Answer {
		log.Println("Found an A Record!!!!", a)
	}

	flatRecords := []dns.RR{}

	for _, answer := range response.Answer {
		fmt.Println("Flattening", m.Answer)
		if answer.Header().Rrtype == dns.TypeA {
			aRecord := answer.(*dns.A)

			fmt.Println("Let's go for answer:", aRecord)
			flatRecord := fmt.Sprintf("%s A %s", r.Header().Name, aRecord.A)
			rr, err := dns.NewRR(flatRecord)
			if err != nil {
				log.Fatal("houlalalala")
			}
			flatRecords = append(flatRecords, rr)
		}
	}

	// Query deeper
	return flatRecords, nil
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		log.Printf("Query for %s\n", q.Name)
		record := records[q.Name]
		if record.Header().Rrtype == dns.TypeCNAME {
			aRecord, err := CNameFlatenning(record)
			if err == nil {
				m.Answer = append(m.Answer, aRecord...)
			}
		} else {
			m.Answer = append(m.Answer, record)

		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

func main() {
	server := dns.Server{
		Addr:     "127.0.0.1:5553",
		Net:      "udp",
		Listener: nil,
	}

	dns.HandleFunc("mydomain.com.", handleDnsRequest)

	f, err := os.Open("./zone.txt")
	if err != nil {
		log.Fatal("Could not open file", err)
	}

	r := bufio.NewReader(f)

	zp := dns.NewZoneParser(r, "mydomain.com.", "./zone.txt")

	for record, ok := zp.Next(); ok; record, ok = zp.Next() {
		log.Println(record, ok)
		if record != nil {
			records[record.Header().Name] = record
		}
	}

	for name, record := range records {
		log.Println(name, ": ", record)
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Could not create server: ", err)
	}
}
