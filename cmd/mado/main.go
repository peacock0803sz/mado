// Package main is the entry point for the mado CLI.
package main

import (
	"fmt"
	"os"

	"github.com/peacock0803sz/mado/internal/ax"
	"github.com/peacock0803sz/mado/internal/cli"
)

// version is injected at build time via -ldflags.
var version = "dev"

func main() {
	svc := ax.NewWindowService()
	cmd := cli.NewRootCmd(svc)

	// inject the version string into the root command
	cmd.Version = version

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
