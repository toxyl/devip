package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func listAliases() {
	fmt.Println("Current loopback aliases:")
	output, err := sudoExec("ip", "address", "show", "dev", "lo")
	if err != nil {
		fmt.Println("Failed to list loopback aliases")
		return
	}
	re := regexp.MustCompile(`inet\s+(\d+\.\d+\.\d+\.\d+)\/32`)
	matches := re.FindAllStringSubmatch(string(output), -1)
	for _, match := range matches {
		fmt.Println(match[1])
	}
	fmt.Println("")
}

func removeAlias(alias string) {
	output, err := sudoExec("ip", "address", "del", alias+"/32", "dev", "lo")
	if err != nil {
		if string(output) != "Cannot assign requested address\n" {
			fmt.Printf("Failed to remove alias %s: %s", alias, string(output))
			return
		} else {
			fmt.Printf("Could not find alias %s, therefore not removed\n", alias)
			return
		}
	}
	fmt.Printf("Removed %s\n", alias)
}

func addAlias(alias string) {
	output, err := sudoExec("ip", "address", "add", alias+"/32", "dev", "lo")
	if err != nil {
		if string(output) != "File exists\n" {
			fmt.Printf("Failed to add alias %s: %s", alias, string(output))
			return
		} else {
			fmt.Printf("Alias %s already exists\n", alias)
			return
		}
	}
	fmt.Printf("Added %s\n", alias)
}

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

func quoteArgs(args []string) string {
	for i, arg := range args {
		if strings.Contains(arg, " ") {
			args[i] = fmt.Sprintf("'%s'", arg)
		}
	}
	return strings.Join(args, " ")
}

func main() {
	action := ""
	if len(os.Args) > 1 {
		action = os.Args[1]
	}

	switch action {
	case "add":
		for _, alias := range os.Args[2:] {
			addAlias(alias)
		}
	case "remove":
		for _, alias := range os.Args[2:] {
			removeAlias(alias)
		}
	default:
		listAliases()
	}
}
