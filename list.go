package main

import (
	"bytes"
	"fmt"
	"net"
	"regexp"
	"sort"
	"strings"
)

func getList() map[string]string {
	res := map[string]string{}
	output, err := sudoExec("ip", "address", "show", "dev", "lo")
	if err != nil {
		fmt.Println(asError("Failed to list loopback aliases"))
		return res
	}
	re := regexp.MustCompile(`inet\s+(\d+\.\d+\.\d+\.\d+)/(\d+)`)
	matches := re.FindAllStringSubmatch(string(output), -1)
	ips := make([]net.IP, 0, len(matches))
	cidrs := make([]string, 0, len(matches))
	for _, match := range matches {
		ip := net.ParseIP(match[1])
		if ip != nil && !isLocalhost(ip.String()) {
			ips = append(ips, ip)
			cidrs = append(cidrs, match[2])
		}
	}
	sort.Slice(ips, func(i, j int) bool {
		return bytes.Compare(ips[i], ips[j]) < 0
	})
	for i, ip := range ips {
		res[ip.String()] = fmt.Sprintf("%s/%s", ip.String(), cidrs[i])
	}
	return res
}

// list prints the status of all loopback aliases.
func list() {
	for ip, cidr := range getList() {
		fmt.Printf("%s: ", asIP(cidr))
		_, err := sudoExec("ping", "-c", "1", "-w", "1", ip)
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
