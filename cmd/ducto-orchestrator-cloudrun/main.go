package main

import (
	"github.com/tommed/ducto-orchestrator/internal/cli"
	"github.com/tommed/ducto-orchestrator/internal/config"
	"os"
)

func main() {
	loader, err := config.NewGCSLoader() // supports `gs://` paths
	if err != nil {
		panic(err)
	}
	os.Exit(cli.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, loader))
}
