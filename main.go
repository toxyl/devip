// Package main provides a command line tool to manage loopback aliases on Linux.
package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

// as applies ANSI escape codes to format a string with bold, italic, and/or underline styles.
func as(str string, bold, italic, underline bool) string {
	formatStr := ""
	if bold {
		formatStr += "\033[1m"
	}
	if italic {
		formatStr += "\033[3m"
	}
	if underline {
		formatStr += "\033[4m"
	}
	return formatStr + str + "\033[0m"
}

// asColor applies ANSI escape codes to color a string and optionally apply bold, italic, and/or underline styles.
func asColor(str string, color int, bold, italic, underline bool) string {
	formatStr := "\033[" + fmt.Sprint(30+color) + "m"
	if bold {
		formatStr += "\033[1m"
	}
	if italic {
		formatStr += "\033[3m"
	}
	if underline {
		formatStr += "\033[4m"
	}
	return formatStr + str + "\033[0m"
}

// asBold applies ANSI escape codes to make a string bold and optionally apply italic and/or underline styles.
func asBold(str string, italic bool, underline bool) string {
	return as(str, true, italic, underline)
}

// asItalic applies ANSI escape codes to make a string italic and optionally apply bold and/or underline styles.
func asItalic(str string, bold bool, underline bool) string {
	return as(str, bold, true, underline)
}

// asUnderline applies ANSI escape codes to make a string underlined and optionally apply bold and/or italic styles.
func asUnderline(str string, bold bool, italic bool) string {
	return as(str, bold, italic, true)
}

// asError applies ANSI escape codes to color a string red and make it bold.
func asError(str string) string {
	return asColor(str, 1, true, false, false)
}

// asWarning applies ANSI escape codes to color a string yellow and make it italic.
func asWarning(str string) string {
	return asColor(str, 3, false, true, false)
}

// asOK applies ANSI escape codes to color a string green and make it bold.
func asOK(str string) string {
	return asColor(str, 2, true, false, false)
}

// asNeutral applies ANSI escape codes to color a string blue.
func asNeutral(str string) string {
	return asColor(str, 4, false, false, false)
}

// isLocalhost checks whether the given alias is a loopback address.
func isLocalhost(alias string) bool {
	ip := net.ParseIP(alias)
	if ip == nil {
		return false
	}
	return ip.IsLoopback()
}

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
		fmt.Printf("%s: ", ip.String())
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

// sudoExec runs the given command with sudo privileges and returns its output.
func sudoExec(name string, arg ...string) ([]byte, error) {
	args := []string{"-S", "-p", "", "sh", "-c", fmt.Sprintf("%s %s", name, quoteArgs(arg))}
	cmd := exec.Command("sudo", args...)
	cmd.Stdin = strings.NewReader("")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, err
	}
	return output, nil
}

// quoteArgs quotes the given arguments so they can be passed as a single argument to a shell command.
func quoteArgs(args []string) string {
	for i, arg := range args {
		if strings.Contains(arg, " ") {
			args[i] = fmt.Sprintf("'%s'", arg)
		}
	}
	return strings.Join(args, " ")
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
