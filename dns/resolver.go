package dns

import (
	"fmt"
	"github.com/miekg/dns"
	"hdns/resolver"
	"net"
)

func Resolver(conn net.PacketConn) {
	// Loop indefinitely to handle incoming requests
	for {
		// Read the incoming DNS request
		buffer := make([]byte, 512)
		_, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Error reading request:", err)
			continue
		}

		// Parse the DNS request
		msg := &dns.Msg{}
		err = msg.Unpack(buffer)
		if err != nil {
			fmt.Println("Error parsing request:", err)
			continue
		}

		// Handle the DNS request
		reply := &dns.Msg{}
		reply.SetReply(msg)

		for _, q := range msg.Question {
			if q.Qtype == dns.TypeA {
				// Respond with an A record for example.com
				ips, err := resolver.RResolveDNS(q.Name)
				if err != nil {
					fmt.Println("Error while resolving:", err)
				}
				for _, ip := range ips {
					resolvedIP := ip
					rr := &dns.A{
						Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
						A:   resolvedIP,
					}
					reply.Answer = append(reply.Answer, rr)
					fmt.Printf("resolving (%s --> %s) for %s\n", q.Name, resolvedIP, addr)
					// #TODO: Add the record to the cache
				}
			}
		}

		// Send the DNS response
		buffer, err = reply.Pack()
		if err != nil {
			fmt.Println("Error packing response:", err)
			continue
		}
		_, err = conn.WriteTo(buffer, addr)
		if err != nil {
			fmt.Println("Error sending response:", err)
			continue
		}
		fmt.Println("Sent DNS response to ", addr)
	}
}
