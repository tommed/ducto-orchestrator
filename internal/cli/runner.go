package cli

import (
	"context"
	"flag"
	"fmt"
	"github.com/tommed/ducto-orchestrator/internal/outputs"
	"github.com/tommed/ducto-orchestrator/internal/sources"
	"io"

	"github.com/tommed/ducto-dsl/transform"

	"github.com/tommed/ducto-orchestrator/internal/config"
	"github.com/tommed/ducto-orchestrator/internal/orchestrator"
)

//goland:noinspection GoUnhandledErrorResult
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("ducto-orchestrator", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var configPath string
	var debug bool
	fs.StringVar(&configPath, "config", "", "Path to the config file")
	fs.BoolVar(&debug, "debug", false, "Enable debug mode")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "failed to parse args: %v\n", err)
		return 1
	}
	if configPath == "" {
		fmt.Fprintln(stderr, "missing required --config path")
		return 1
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(stderr, "failed to load config: %v\n", err)
		return 1
	}

	// Listen out for Ctl+C / Signal Interrupt
	ctx := WithSignalContext(context.Background())

	// Load program
	var prog *transform.Program
	if cfg.Program != nil {
		prog = cfg.Program
	} else if cfg.ProgramFile != "" {
		prog, err = transform.LoadProgram(cfg.ProgramFile)
		if err != nil {
			fmt.Fprintf(stderr, "%v\n", err)
			return 1
		}
	} else {
		fmt.Fprintln(stderr, "no DSL program or program_file defined")
		return 1
	}
	// Debug can come from cli flag or config, if any are true
	if debug {
		cfg.Debug = true
	}

	// Load source
	source, err := sources.FromPlugin(ctx, cfg.Source, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "failed to load source: %v\n", err)
		return 1
	}
	defer source.Close()

	// Load output
	output, err := outputs.FromPlugin(cfg.Output, stdout)
	if err != nil {
		fmt.Fprintf(stderr, "failed to load output: %v\n", err)
		return 1
	}

	// Run orchestrator
	o := orchestrator.New(prog, cfg.Debug)
	if err := o.RunLoop(ctx, source, output); err != nil {
		fmt.Fprintf(stderr, "orchestrator failed: %v\n", err)
		return 1
	}

	return 0
}
