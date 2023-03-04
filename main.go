// Package main provides a command line tool to manage loopback aliases on Linux.
package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"regexp"
	"sort"
	"strings"
)

// list prints the status of all loopback aliases.
func list() {
	output, err := sudoExec("ip", "address", "show", "dev", "lo")
	if err != nil {
		fmt.Println(asError("Failed to list loopback aliases"))
		return
	}
	re := regexp.MustCompile(`inet\s+(\d+\.\d+\.\d+\.\d+)\/32`)
	matches := re.FindAllStringSubmatch(string(output), -1)
	ips := make([]net.IP, 0, len(matches))
	for _, match := range matches {
		ip := net.ParseIP(match[1])
		if ip != nil {
			ips = append(ips, ip)
		}
	}
	sort.Slice(ips, func(i, j int) bool {
		return bytes.Compare(ips[i], ips[j]) < 0
	})
	for _, ip := range ips {
		fmt.Printf("%s: ", asIP(ip.String()))
		_, err := sudoExec("ping", "-c", "1", "-w", "1", ip.String())
		if err != nil {
			if strings.HasSuffix(err.Error(), "exit status 1") {
				fmt.Printf("%s\n", asError("down"))
				continue
			}
			fmt.Printf("Ping error: %s\n", asError(err.Error()))
			continue
		}
		fmt.Printf("%s\n", asOK("up"))
	}
}

// remove removes the loopback alias of the given IP.
func remove(ip string) {
	fmt.Printf("Removing %s ... ", asIP(ip))

	if isLocalhost(ip) {
		fmt.Printf("%s %s\n", as("seriously?", true, true, false), asError("Not gonna do that!"))
		return
	}

	output, err := sudoExec("ip", "address", "del", ip+"/32", "dev", "lo")
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

// add adds a loopback alias for the given IP.
func add(ip string) {
	fmt.Printf("Adding %s ... ", asIP(ip))
	output, err := sudoExec("ip", "address", "add", ip+"/32", "dev", "lo")
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

// main is the entry point of the program.
//
// If no arguments are provided, it will display a list of available aliases.
// Otherwise, it will execute the appropriate function based on the first argument.
func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		list()
		return
	}
	action := args[0]
	args = args[1:]
	switch action {
	case "add":
		for _, alias := range args {
			add(alias)
		}
	case "remove":
		for _, alias := range args {
			remove(alias)
		}
	case "test":
		for _, alias := range args {
			test(alias)
		}
	default:
		list()
	}
}
