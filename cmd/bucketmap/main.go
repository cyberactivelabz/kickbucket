package main

import (
	"fmt"
	"os"
)

// Set at build time by GoReleaser (see .goreleaser.yaml).
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("bucketmap %s (%s, %s)\n", version, commit, date)
		return
	}
	fmt.Fprintln(os.Stderr, "usage: bucketmap scan <keyword> (not implemented yet)")
	os.Exit(2)
}
