package main

import (
	"github.com/tommed/ducto-orchestrator/internal/cli"
	"github.com/tommed/ducto-orchestrator/internal/config"
	"os"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, &config.FileLoader{}))
}
