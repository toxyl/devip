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

// List prints the status of all loopback aliases.
func List() {
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
		fmt.Printf("%s: ", asNeutral(ip.String()))
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

// Remove removes the loopback alias of the given IP.
func Remove(ip string) {
	fmt.Printf("Removing %s ... ", asNeutral(ip))

	if isLocalhost(ip) {
		fmt.Printf("%s %s\n", asItalic("seriously?", true, false), asError("Not gonna do that!"))
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

// Add adds a loopback alias for the given IP.
func Add(ip string) {
	fmt.Printf("Adding %s ... ", asNeutral(ip))
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

// Test pings the given IP and prints its status.
func Test(ip string) {
	fmt.Printf("Testing %s ... ", asNeutral(ip))
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

// Run executes the appropriate function based on the action and arguments provided.
//
// Parameters:
//   - action: The action to perform (add, remove, test, or no action which default to list).
//   - args: The arguments to pass to the function.
func Run(action string, args []string) {
	switch action {
	case "add":
		for _, alias := range args {
			Add(alias)
		}
	case "remove":
		for _, alias := range args {

			Remove(alias)
		}
	case "test":
		for _, alias := range args {
			Test(alias)
		}
	default:
		List()
	}
}

// main is the entry point of the program.
//
// If no arguments are provided, it will display a list of available aliases.
// Otherwise, it will execute the appropriate function based on the first argument.
func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		List()
		return
	}
	action := args[0]
	args = args[1:]
	Run(action, args)
}
