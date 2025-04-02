package orchestrator

import (
	"context"
	"fmt"
	"github.com/tommed/ducto-dsl/transform"
)

type Orchestrator struct {
	Program *transform.Program
	Debug   bool
}

func New(program *transform.Program, debug bool) *Orchestrator {
	return &Orchestrator{
		Program: program,
		Debug:   debug,
	}
}

func (o *Orchestrator) Execute(ctx context.Context, input map[string]interface{}, writer OutputWriter) error {

	// Context (with flags)
	ctx = context.WithValue(ctx,
		transform.ContextKeyDebug, o.Debug)

	// Apply transformation
	output, err := transform.New().Apply(ctx, input, o.Program)
	if err != nil {
		return fmt.Errorf("failed to apply program: %w", err)
	}

	// Write Output
	if err := writer.WriteOutput(output); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}
