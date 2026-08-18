package main

import (
	"fmt"
	"os"
)

const defaultVersion = "0.1.0"

// version is overridden at link time via -ldflags "-X main.version=..."
var version = defaultVersion

func main() {
	if len(os.Args) < 2 {
		printUsage(os.Stderr)
		os.Exit(2)
	}

	cmd, args := os.Args[1], os.Args[2:]
	switch cmd {
	case "serve":
		if err := runServe(args); err != nil {
			fmt.Fprintf(os.Stderr, "git-lore: %v\n", err)
			os.Exit(1)
		}
	case "version", "--version", "-V":
		fmt.Printf("git-lore %s\n", version)
	case "help", "--help", "-h":
		printUsage(os.Stdout)
	default:
		fmt.Fprintf(os.Stderr, "git-lore: unknown command %q\n\n", cmd)
		printUsage(os.Stderr)
		os.Exit(2)
	}
}

func printUsage(w *os.File) {
	fmt.Fprintf(w, `git-lore — optional CLI for Lore Works (refs/lore/*)

Usage:
  git-lore <command> [flags]

Commands:
  serve      Serve the local Lore browser UI
  version    Print version
  help       Show this help

The primary git-lore surface remains agent skills under skills/.
This CLI is for humans who prefer commands over skills; more
subcommands (list, show, create, sync, …) can land here later.

`)
}
