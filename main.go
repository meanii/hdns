package main

import (
	"fmt"
	"hdns/dns"
	"net"
)

func main() {
	conn, err := net.ListenPacket("udp", ":53")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("DNS server listening on port 53...")

	// init dns resolver
	dns.Resolver(conn)
}
