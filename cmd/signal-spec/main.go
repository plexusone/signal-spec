// Command signal-spec is the CLI for signal-spec operations.
package main

import (
	"os"

	"github.com/plexusone/signal-spec/cmd/signal-spec/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
