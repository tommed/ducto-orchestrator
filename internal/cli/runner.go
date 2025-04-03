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
	programPath := fs.String("program", "", "Path to the Ducto DSL program JSON file")
	debug := fs.Bool("debug", false, "Enable debug output")

	// Parse the flags
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "failed to parse args: %v\n", err)
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
	//source := orchestrator.NewStdinEventSource(stdin)
	source := orchestrator.NewHTTPEventSource("127.0.0.1:8080", "_http")
	writer := orchestrator.NewStdoutWriter(stdout)

	// Run the Loop
	err = orchestrator.New(prog, *debug).RunLoop(ctx, source, writer)
	if err != nil {
		fmt.Fprintf(stderr, "execution error: %v\n", err)
		return 1
	}

	return 0
}
