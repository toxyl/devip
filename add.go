package main

import (
	"fmt"
	"net"
	"strings"
)

// add adds a loopback alias for the given IP.
func add(ip string) {
	fmt.Printf("Adding %s ... ", asIP(ip))
	_, _, err := net.ParseCIDR(ip)
	if err != nil {
		// we assume it didn't contain a CIDR, so it's a single machine
		ip += "/32"
	}
	output, err := sudoExec("ip", "address", "add", ip, "dev", "lo")
	if err != nil {
		if !strings.Contains(string(output), "File exists") {
			fmt.Printf("failed with error: %s\n", asError(err.Error()))
			return
		} else {
			fmt.Printf("already exists, skipping!\n")
			return
		}
	}
	fmt.Printf("%s\n", asOK("done!"))
}
