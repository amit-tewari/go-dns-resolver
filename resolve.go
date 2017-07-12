package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"sort"
	"strings"
	"time"
)

func do_read_domains(domains chan<- string,
	domainSlotAvailable <-chan bool) {
	in := bufio.NewReader(os.Stdin)

	for _ = range domainSlotAvailable {

		input, err := in.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "read(stdin): %s\n", err)
			os.Exit(1)
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		domain := input + "."
		domains <- domain
	}
	close(domains)
}

var sendingDelay time.Duration
var retryDelay time.Duration
var packetsPerSecond int
var concurrency int
var dnsServerPool string
var retryTime string
var verbose bool
var ipv6 bool
var soa bool
var txt bool
var mx bool
var all bool
var a bool

func init() {
	flag.StringVar(&dnsServerPool, "serverPool", "8.8.8.8,8.8.4.4",
		"comma seperated DNS server address")
	flag.IntVar(&concurrency, "concurrency", 5000,
		"Internal buffer")
	flag.IntVar(&packetsPerSecond, "pps", 120,
		"Send up to PPS DNS queries per second")
	flag.StringVar(&retryTime, "retry", "1s",
		"Resend unanswered query after RETRY")
	flag.BoolVar(&verbose, "v", false,
		"Verbose logging")
	flag.BoolVar(&soa, "soa", false,
		"Query SOA records")
	flag.BoolVar(&mx, "mx", false,
		"Query MX records")
	flag.BoolVar(&txt, "txt", false,
		"Query TXT records")
	flag.BoolVar(&a, "a", false,
		"Query A records")
	flag.BoolVar(&ipv6, "6", false,
		"Ipv6 - ask for AAAA, not A")
	flag.BoolVar(&all, "all", false,
		"Perform lookups for all query types")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, strings.Join([]string{
			"\"resolve\" mass resolve DNS A records for domains names read from stdin.",
			"",
			"Usage: resolve [option ...]",
			"",
		}, "\n"))
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(1)
	}

	allQTypes := []uint16{dnsTypeA, dnsTypeSOA, dnsTypeMX, dnsTypeTXT, dnsTypeAAAA}
	queryFor := make([]uint16, 0, len(allQTypes))
	if all {
		queryFor = allQTypes
	} else {
		if a {
			queryFor = append(queryFor, dnsTypeA)
		}
		if soa {
			queryFor = append(queryFor, dnsTypeSOA)
		}
		if mx {
			queryFor = append(queryFor, dnsTypeMX)
		}
		if txt {
			queryFor = append(queryFor, dnsTypeTXT)
		}
		if ipv6 {
			queryFor = append(queryFor, dnsTypeAAAA)
		}
		if len(queryFor) == 0 {
			queryFor = append(queryFor, dnsTypeA)
		}
	}

	dnsServers := strings.Split(dnsServerPool, ",")
	healthyDnsServers := make([]string, 0, len(dnsServers))
	dnsConnectionPool := make([]net.Conn, 0, len(dnsServers))

	var err error
	for _, server := range dnsServers {
		c, err := net.Dial("udp", server+":53")
		if err != nil {
			fmt.Fprintf(os.Stderr, "bind(udp, %s): %s\n", server, err)
		} else {
			id_sent := uint16(rand.Int())
			msg := packDns("google.com", id_sent, dnsTypeA)
			n, _ := c.Write(msg)
			c.SetReadDeadline(time.Now().Add(2 * time.Second))
			fmt.Fprintf(os.Stderr, "Checking if server %s responding correcly under 2 second delay! ", server)
			buf := make([]byte, 4096)
			n, err = c.Read(buf)
			if err == nil {
				c.SetReadDeadline(time.Time{})
				domain, id_returned, _, _, _ := unpackDns(buf[:n])
				if domain == "google.com." && id_sent == id_returned {
					fmt.Fprintf(os.Stderr, "server verified and added to pool: %s\n", server)
					dnsConnectionPool = append(dnsConnectionPool, c)
					healthyDnsServers = append(healthyDnsServers, server)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Encountered error with %s, not adding to DNS pool\n", server)
			}
		}
	}
	if len(dnsConnectionPool) == 0 {
		fmt.Println("No connection could be established")
		os.Exit(1)
	}
	sendingDelay = time.Duration(1000000000/packetsPerSecond) * time.Nanosecond
	retryDelay, err = time.ParseDuration(retryTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't parse duration %s\n", retryTime)
		os.Exit(1)
	}
	qtypes := fmt.Sprintf("types queried : %v", queryFor)
	fmt.Fprintf(os.Stderr, "\n%d Server in pool : %v, %s sending delay: %s (%d pps), retry delay: %s\n\n",
		len(healthyDnsServers), healthyDnsServers, qtypes, sendingDelay, packetsPerSecond, retryDelay)

	domainSlotAvailable := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		domainSlotAvailable <- true
	}
	domains := make(chan string, concurrency)

	go do_read_domains(domains,
		domainSlotAvailable)

	// Used as a queue. Make sure it has plenty of storage available.
	timeoutRegister := make(chan *domainRecord, concurrency*1000)
	timeoutExpired := make(chan *domainRecord)

	resolved := make(chan *domainAnswer, concurrency*5)
	tryResolving := make(chan *domainRecord, concurrency)

	go do_timeouter(timeoutRegister, timeoutExpired)

	go do_send(dnsConnectionPool, tryResolving)
	for poolIndex, rConn := range dnsConnectionPool {
		go do_receive(rConn, resolved, poolIndex)
	}

	t0 := time.Now()
	domainsCount, avgTries := do_map_guard(domains,
		domainSlotAvailable,
		timeoutRegister,
		timeoutExpired,
		tryResolving,
		resolved,
		allQTypes,
		queryFor)
	td := time.Now().Sub(t0)
	fmt.Fprintf(os.Stderr, "\nResolved %d domains in %.3fs. Average retries %.3f. Domains per second: %.3f\n",
		domainsCount,
		td.Seconds(),
		avgTries,
		float64(domainsCount)/td.Seconds())
}

type domainRecord struct {
	id          uint16
	domain      string
	timeout     time.Time
	resend      int
	time_queued time.Time
	time_sent   time.Time
	dnsQtype    uint16
}

type domainAnswer struct {
	id           uint16
	domain       string
	ips          []net.IP
	soa_t        string
	dnsQtype     uint16
	time_resolve time.Time
}

func do_map_guard(domains <-chan string,
	domainSlotAvailable chan<- bool,
	timeoutRegister chan<- *domainRecord,
	timeoutExpired <-chan *domainRecord,
	tryResolving chan<- *domainRecord,
	resolved <-chan *domainAnswer,
	allQTypes []uint16,
	queryFor []uint16) (int, float64) {

	m := make(map[uint16]*domainRecord)

	done := false

	sumTries := 0
	domainCount := 0
	for done == false || len(m) > 0 {
		select {
		case domain := <-domains:
			if domain == "" {
				domains = make(chan string)
				done = true
				break
			}
			for _, qtype := range queryFor {
				var id uint16
				for {
					id = uint16(rand.Int())
					if id != 0 && m[id] == nil {
						break
					}
				}
				time_now := time.Now()
				dr := &domainRecord{id,
					domain,
					time_now,
					1,
					time_now,
					time_now,
					qtype}
				m[id] = dr
				if verbose {
					fmt.Fprintf(os.Stderr, "0x%04x resolving %s\n", id, domain)
				}
				timeoutRegister <- dr
				tryResolving <- dr
			}

		case dr := <-timeoutExpired:
			if m[dr.id] == dr {
				dr.resend += 1
				dr.timeout = time.Now()
				if verbose {
					fmt.Fprintf(os.Stderr, "0x%04x resend (try:%d) %s\n", dr.id,
						dr.resend, dr.domain)
				}
				if dr.resend < 3 {
					timeoutRegister <- dr
					tryResolving <- dr
				}
			}

		case da := <-resolved:
			if m[da.id] != nil {
				dr := m[da.id]
				if dr.domain != da.domain {
					if verbose {
						fmt.Fprintf(os.Stderr, "0x%04x error, unrecognized domain: %s != %s\n",
							da.id, dr.domain, da.domain)
					}
					break
				}

				if verbose {
					fmt.Fprintf(os.Stderr, "0x%04x resolved %s\n",
						dr.id, dr.domain)
				}

				s := make([]string, 0, 16)
				for _, ip := range da.ips {
					s = append(s, ip.String())
				}
				sort.Sort(sort.StringSlice(s))

				// without trailing dot
				domain := dr.domain[:len(dr.domain)-1]
				//
				fmt.Printf("%s %d %d %d %s%s\n",
					domain,
					da.dnsQtype,
					int64(dr.time_sent.Sub(dr.time_queued)/time.Millisecond),
					int64(da.time_resolve.Sub(dr.time_sent)/time.Millisecond),
					strings.Join(s, " "),
					da.soa_t)

				sumTries += dr.resend
				domainCount += 1

				delete(m, dr.id)
				select {
				case domainSlotAvailable <- true:
				default:
				}
			}
		}
	}
	return domainCount, float64(sumTries) / float64(domainCount)
}

func do_timeouter(timeoutRegister <-chan *domainRecord,
	timeoutExpired chan<- *domainRecord) {
	for {
		dr := <-timeoutRegister
		t := dr.timeout.Add(retryDelay)
		now := time.Now()
		if t.Sub(now) > 0 {
			delta := t.Sub(now)
			time.Sleep(delta)
		}
		timeoutExpired <- dr
	}
}

func do_send(c []net.Conn, tryResolving <-chan *domainRecord) {
	poolLength := len(c)
	target := 0
	for {
		dr := <-tryResolving

		msg := packDns(dr.domain, dr.id, dr.dnsQtype)
		if verbose {
			fmt.Printf("sending %s type %d ", dr.domain, dr.dnsQtype)
		}
		dr.time_sent = time.Now()
		target = target % poolLength
		_, err := c[target].Write(msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "write(udp): %s\n", err)
			fmt.Fprintf(os.Stderr, "error writing to Servers [%d]! %s will be retried with another server\nError: %s\n", target, dr.domain, err)
			//os.Exit(1)
		}
		target++
		time.Sleep(sendingDelay)
	}
}

func do_receive(c net.Conn, resolved chan<- *domainAnswer, poolIndex int) {
	buf := make([]byte, 4096)
	for {
		n, err := c.Read(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Server %d in pool seem to have issues reading replies! TAKING OUT OF POOL [Error %s]\n",
				poolIndex+1, err)
			return
		}

		domain, id, ips, soa_t, dnsQType := unpackDns(buf[:n])
		if verbose {
			fmt.Printf(" received %s type %d\n", domain, dnsQType)
		}
		resolved <- &domainAnswer{id, domain, ips, soa_t, dnsQType, time.Now()}
		if verbose {
			fmt.Printf(" pushed on received chan %s type %d\n", domain, dnsQType)
		}
	}
}
