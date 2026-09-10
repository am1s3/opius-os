package main

import (
	"os"

	"github.com/opius-os/opius/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}