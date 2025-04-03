package cli

import (
	"context"
	"flag"
	"fmt"
	"github.com/tommed/ducto-dsl/transform"
	"github.com/tommed/ducto-orchestrator/internal/orchestrator"
	"io"
)

//goland:noinspection GoUnhandledErrorResult
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	ctx := orchestrator.WithSignalContext(context.Background())

	// Parse flags from args provided not os.Args (for testing)
	fs := flag.NewFlagSet("ducto-orchestrator", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// User provided flags
	sourceType := fs.String("source", "stdin", "event source (stdin, http, pubsub, ...)")
	programPath := fs.String("program", "", "Path to the Ducto DSL program JSON file")
	debug := fs.Bool("debug", false, "Enable debug output")
	addr := fs.String("addr", "127.0.0.1:8080", "the address to bind to when --source http")
	metaField := fs.String("meta", "", "when supplied, this is the field http meta-data will be set to with --source http")

	// Parse the flags
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "failed to parse args: %v\n", err)
		return 1
	}

	if *sourceType == "" {
		fmt.Fprintf(stderr, "missing --source stdin|http|...\n")
		fs.Usage()
		return 1
	}

	if *programPath == "" {
		fmt.Fprintf(stderr, "missing --program <file>\n")
		fs.Usage()
		return 1
	}

	prog, err := transform.LoadProgram(*programPath)
	if err != nil {
		fmt.Fprintf(stderr, "failed to load program: %v\n", err)
		return 1
	}

	// Create components
	var source orchestrator.EventSource

	switch *sourceType {
	case "stdin":
		source = orchestrator.NewStdinEventSource(stdin)
	case "http":
		source = orchestrator.NewHTTPEventSource(*addr, *metaField)
	default:
		fmt.Fprintf(stderr, "unknown source type: %s\n", sourceType)
		return 1
	}

	//source := orchestrator.NewStdinEventSource(stdin)
	writer := orchestrator.NewStdoutWriter(stdout)

	// Run the Loop
	err = orchestrator.New(prog, *debug).RunLoop(ctx, source, writer)
	if err != nil {
		fmt.Fprintf(stderr, "execution error: %v\n", err)
		return 1
	}

	return 0
}
