package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

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

// isLocalhost checks whether the given alias is a loopback address.
func isLocalhost(alias string) bool {
	ip := net.ParseIP(alias)
	if ip == nil {
		return false
	}
	return ip.IsLoopback()
}
