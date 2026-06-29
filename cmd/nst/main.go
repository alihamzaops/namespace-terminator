package main

import (
	"os"

	"github.com/alihamzaops/namespace-terminator/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
