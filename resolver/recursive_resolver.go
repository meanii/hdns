package resolver

import (
	"fmt"
	"github.com/miekg/dns"
	"net"
)

func RResolveDNS(name string) ([]net.IP, error) {
	ips, err := net.LookupIP(name)
	if err == nil {
		return ips, nil
	}

	// If the error is not a "no such host" error, return the error
	if _, ok := err.(*net.DNSError); !ok {
		return nil, err
	}

	// If the error is a "no such host" error, try resolving the name with the
	// DNS servers listed in /etc/resolv.conf
	config, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return nil, err
	}

	for _, server := range config.Servers {
		ips, err = resolveDNSWithServer(name, server)
		if err == nil {
			return ips, nil
		}
	}

	return nil, fmt.Errorf("could not resolve %s", name)
}

func resolveDNSWithServer(name string, server string) ([]net.IP, error) {
	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(dns.Fqdn(name), dns.TypeA)
	r, _, err := c.Exchange(&m, net.JoinHostPort(server, "53"))
	if err != nil {
		return nil, err
	}

	if r.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("DNS query failed with error code %d", r.Rcode)
	}

	ips := make([]net.IP, 0)
	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.A); ok {
			ips = append(ips, a.A)
		}
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("no A records found for %s", name)
	}

	return ips, nil
}
