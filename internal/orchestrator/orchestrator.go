package orchestrator

import (
	"context"
	"fmt"
	"github.com/tommed/ducto-dsl/transform"
	"github.com/tommed/ducto-orchestrator/internal/sources"
	"github.com/tommed/ducto-orchestrator/internal/writers"
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

func (o *Orchestrator) RunLoop(ctx context.Context, source sources.EventSource, writer writers.OutputWriter) error {

	// Setup Source
	events, err := source.Start(ctx)
	if err != nil {
		return err
	}

	// Teardown Source
	defer func(source sources.EventSource) {
		_ = source.Close()
	}(source)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case evt, ok := <-events:
			if !ok {
				return nil // stream closed naturally
			}

			if err := o.RunOnce(ctx, evt, writer); err != nil {
				return err
			}
		}
	}
}

func (o *Orchestrator) RunOnce(ctx context.Context, input map[string]interface{}, writer writers.OutputWriter) error {

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
