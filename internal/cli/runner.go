package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tommed/ducto-dsl/transform"
	"io"
)

//goland:noinspection GoUnhandledErrorResult
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("ducto-orchestrator", flag.ContinueOnError)
	fs.SetOutput(stderr)

	programPath := fs.String("program", "", "Path to the Ducto DSL program JSON file")
	debug := fs.Bool("debug", false, "Enable debug output")

	// Parse the injected args
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "failed to parse args: %v\n", err)
		return 1
	}

	if *programPath == "" {
		fmt.Fprintf(stderr, "missing --program <file>\n")
		fs.Usage()
		return 1
	}

	// Our context
	ctx := context.WithValue(context.Background(),
		transform.ContextKeyDebug, *debug)

	prog, err := transform.LoadProgram(*programPath)
	if err != nil {
		fmt.Fprintf(stderr, "failed to load program: %v\n", err)
		return 1
	}

	inputData, err := io.ReadAll(stdin)
	if err != nil {
		fmt.Fprintf(stderr, "failed to read stdin: %v\n", err)
		return 1
	}

	var input map[string]interface{}
	if err := json.Unmarshal(inputData, &input); err != nil {
		fmt.Fprintf(stderr, "invalid input JSON: %v\n", err)
		return 1
	}

	output, err := transform.New().Apply(ctx, input, prog)
	if err != nil {
		fmt.Fprintf(stderr, "failed to execute program: %v\n", err)
		return 1
	}

	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(output)
	return 0
}
