package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/u-root/dhcp4/dhcp4server"
)

var (
	subnet = flag.String("subnet", "192.168.1.0/24", "IP subnet to use to allocate over DHCP (must be CIDR notation)")
	self   = flag.String("self", "192.168.0.1", "My own IP")
)

func main() {
	_, sn, err := net.ParseCIDR(*subnet)
	if err != nil {
		log.Fatalf("Could not parse CIDR for subnet %q: %v", *subnet, err)
	}

	l, err := net.ListenPacket("udp4", ":67")
	if err != nil {
		log.Fatalf("Could not listen on udp port 67: %v", err)
	}
	defer l.Close()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	s := dhcp4server.New(net.ParseIP(*self), sn, "", "")

	// This should be an "infinite loop".
	if err := s.Serve(logger, l); err != nil {
		log.Fatalf("Serve DHCP failed: %v", err)
	}
}
