package cli

import (
	"context"
	"flag"
	"fmt"
	"github.com/tommed/ducto-orchestrator/internal/outputs"
	"github.com/tommed/ducto-orchestrator/internal/sources"
	"io"

	"github.com/tommed/ducto-orchestrator/internal/config"
	"github.com/tommed/ducto-orchestrator/internal/orchestrator"
)

//goland:noinspection GoUnhandledErrorResult
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer, cfgLoader config.Loader) int {
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

	// Listen out for Ctl+C / Signal Interrupt
	ctx := WithSignalContext(context.Background())

	cfg, err := cfgLoader.Load(ctx, configPath)
	if err != nil {
		fmt.Fprintf(stderr, "failed to load config: %v\n", err)
		return 1
	}

	// Load program (the cfg loader will have already set this inline for us)
	if cfg.Program == nil {
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
	output, err := outputs.FromPlugin(ctx, cfg.Output, stdout)
	if err != nil {
		fmt.Fprintf(stderr, "failed to load output: %v\n", err)
		return 1
	}

	// Run orchestrator
	o := orchestrator.New(cfg.Program, cfg.Debug)

	// Install our preprocessors
	if err := o.InstallPreprocessorsFromConfig(ctx, cfg.Preprocessors); err != nil {
		fmt.Fprintf(stderr, "failed to install preprocessors: %v\n", err)
		return 1
	}

	// Run the Loop
	if err := o.RunLoop(ctx, source, output); err != nil {
		fmt.Fprintf(stderr, "orchestrator failed: %v\n", err)
		return 1
	}

	return 0
}
