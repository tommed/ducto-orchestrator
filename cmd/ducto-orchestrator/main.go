package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/tommed/ducto-dsl/transform"
)

func main() {
	programPath := flag.String("program", "", "Path to the Ducto DSL program JSON file")
	debug := flag.Bool("debug", false, "Enable debug output")
	flag.Parse()

	// our context with debug switch set
	ctx := context.WithValue(context.Background(),
		transform.ContextKeyDebug, *debug)

	if *programPath == "" {
		fmt.Println("missing --program <file>")
		os.Exit(1)
	}

	prog, err := transform.LoadProgram(*programPath)
	if err != nil {
		fmt.Printf("failed to load program: %v\n", err)
		os.Exit(1)
	}

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("failed to read stdin: %v\n", err)
		os.Exit(1)
	}

	var input map[string]interface{}
	if err := json.Unmarshal(inputData, &input); err != nil {
		fmt.Printf("invalid input JSON: %v\n", err)
		os.Exit(1)
	}

	output, err := transform.New().Apply(ctx, input, prog)
	if err != nil {
		fmt.Printf("failed to execute program: %v\n", err)
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(output)
}
