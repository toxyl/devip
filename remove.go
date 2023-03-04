package main

import (
	"fmt"
	"net"
	"strings"
)

// remove removes the loopback alias of the given IP.
func remove(ip string) {
	fmt.Printf("Removing %s ... ", asIP(ip))

	if isLocalhost(ip) {
		fmt.Printf("%s %s\n", as("seriously?", true, true, false), asError("Not gonna do that!"))
		return
	}

	_, _, err := net.ParseCIDR(ip)
	if err != nil {
		// we assume it didn't contain a CIDR, so it's a single machine
		ip += "/32"
	}

	output, err := sudoExec("ip", "address", "del", ip, "dev", "lo")
	if err != nil {
		if strings.Contains(string(output), "Cannot assign requested address") {
			fmt.Printf("not found, skipping!\n")
			return
		} else {
			fmt.Printf("failed with error: %s\n", asError(string(output)))
			return
		}
	}
	fmt.Printf("%s\n", asOK("done!"))
}
