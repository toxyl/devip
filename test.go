package main

import (
	"fmt"
	"strings"
)

// test pings the given IP and prints its status.
func test(ip string) {
	fmt.Printf("Testing %s ... ", asIP(ip))
	_, err := sudoExec("ping", "-c", "1", "-w", "1", ip)
	if err != nil {
		if strings.HasSuffix(err.Error(), "exit status 1") {
			fmt.Printf("%s\n", asError("down"))
			return
		}
		fmt.Printf("Ping error: %s\n", asError(err.Error()))
	} else {
		fmt.Printf("%s\n", asOK("up"))
	}
}
