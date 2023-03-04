// Package main provides a command line tool to manage loopback aliases on Linux.
package main

import (
	"os"
)

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
		if len(args) == 1 && args[0] == "all" {
			for _, cidr := range getList() {
				remove(cidr)
			}
			return
		}
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
