package main

import (
	"fmt"
	"net"
	"os"
)

func unpackDns(msg []byte) (domain string, id uint16, ips []net.IP, soa_t string, dnsQType uint16) {
	d := new(dnsMsg)
	if !d.Unpack(msg) {
		// fmt.Fprintf(os.Stderr, "dns error (unpacking)\n")
		return
	}

	domain = d.question[0].Name
	if len(domain) < 1 {
		// fmt.Fprintf(os.Stderr, "dns error (wrong domain in question)\n")
		return
	}
	id = d.id
	dnsQType = d.question[0].Qtype

	switch d.rcode {
	case dnsRcodeSuccess:
		soa_t = ""
	case dnsRcodeFormatError:
		soa_t = "-ERR-FORMAT"
	case dnsRcodeServerFailure:
		soa_t = "-ERR-SERVFAIL"
	case dnsRcodeNameError:
		soa_t = "-ERR-NAME-ERROR"
	case dnsRcodeNotImplemented:
		soa_t = "-ERR-NOT-IMPLEMENTED"
	case dnsRcodeRefused:
		soa_t = "-ERR-REFUSED"
	}
	if d.rcode != dnsRcodeSuccess {
		return
	}

	if len(d.question) < 1 {
		// fmt.Fprintf(os.Stderr, "dns error (wrong question section)\n")
		return
	}

	_, addrs, err := answer(domain, "server", d, dnsQType)
	//fmt.Printf (printStruct(d.answer[0]) + "\n")
	//fmt.Println(d.String())
	//fmt.Println(dnsQType)
	if err == nil {
		switch dnsQType {
		case dnsTypeA:
			ips = convertRR_A(addrs)
		case dnsTypeAAAA:
			ips = convertRR_AAAA(addrs)
		case dnsTypeSOA:
			soa_t = convertRR_SOA(addrs)
		case dnsTypeMX:
			soa_t = convertRR_MX(addrs)
		case dnsTypeTXT:
			soa_t = convertRR_TXT(addrs)
		}
	}
	return
}

func packDns(domain string, id uint16, dnsType uint16) []byte {

	out := new(dnsMsg)
	out.id = id
	out.recursion_desired = true
	out.question = []dnsQuestion{
		{domain, dnsType, dnsClassINET},
	}

	msg, ok := out.Pack()
	if !ok {
		fmt.Fprintf(os.Stderr, "can't pack domain %s\n", domain)
		os.Exit(1)
	}
	return msg
}
