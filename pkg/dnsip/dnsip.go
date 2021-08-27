package dnsip

import (
	"fmt"
	"net"
	"strings"
)

// lookup looks domain.
// It returns a slice of that host's IPv4 and IPv6 (or "none" if no DNS record found), tab-separated addresses and domain.
func lookup(domain string) (res []string) {
	iplist, err := net.LookupIP(domain)
	if err != nil {
		res = append(res, fmt.Sprintf("none\t\t%s", domain))
	}
	if len(iplist) > 0 {
		for _, ip := range iplist {
			res = append(res, fmt.Sprintf("%s\t%s", ip.String(), domain))
		}
	}
	return
}

func worker(domainCh chan string, resCh chan []string) {
	for {
		domain := <-domainCh
		res := lookup(domain)
		resCh <- res
	}
}

func generator(domain string, domainCh chan string) {
	domainCh <- domain
}

// ADnsLookup looks up domains.
// It returns a results slice of that host's IPv4 and IPv6 addresses and a none slice containing domains with no DNS record.
func Lookup(threads int, domains []string, out bool) (results []string, none []string) {
	var res []string
	var resNone []string

	// Empty domain list
	if len(domains) == 0 {
		return res, resNone
	}

	// Single target
	if len(domains) == 1 {
		for _, d := range lookup(domains[0]) {
			if strings.HasPrefix(d, "none") {
				resNone = append(resNone, d[6:]) // "Remove" none string.
				if out {
					fmt.Println(d)
				}
				continue
			}
			res = append(res, d)
			if out {
				fmt.Println(d)
			}
		}
	} else {
		domainCh := make(chan string)
		resCh := make(chan []string)

		for i := 0; i < threads; i++ {
			go worker(domainCh, resCh)
		}

		for _, domain := range domains {
			go generator(domain, domainCh)
		}

		if out {
			fmt.Print("\n:Domains found:\n\n")
		}
		for i := 0; i < len(domains); i++ {
			r := <-resCh
			for _, d := range r {
				if strings.HasPrefix(d, "none") {
					resNone = append(resNone, d[6:]) // "Remove" none string.
					continue
				}
				res = append(res, d)
				if out {
					fmt.Println(d)
				}
			}
		}

		if out && len(resNone) > 0 {
			fmt.Print("\n:No DNS record:\n\n")
			for _, none := range resNone {
				fmt.Println(none)
			}
		}
	}

	return res, resNone

}
