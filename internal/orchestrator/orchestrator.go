package orchestrator

import (
	"context"
	"fmt"
	"github.com/tommed/ducto-dsl/transform"
	"github.com/tommed/ducto-orchestrator/internal/outputs"
	"github.com/tommed/ducto-orchestrator/internal/sources"
)

type Orchestrator struct {
	Program       *transform.Program
	Debug         bool
	preprocessors []Preprocessor
}

func New(program *transform.Program, debug bool) *Orchestrator {
	return &Orchestrator{
		Program: program,
		Debug:   debug,
	}
}

func (o *Orchestrator) AddPreprocessor(p Preprocessor) {
	o.preprocessors = append(o.preprocessors, p)
}

func (o *Orchestrator) RunLoop(ctx context.Context, source sources.EventSource, writer outputs.OutputWriter) error {

	// Setup Provider
	events, err := source.Start(ctx)
	if err != nil {
		return err
	}

	// Teardown Provider
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

func (o *Orchestrator) RunOnce(ctx context.Context, input map[string]interface{}, writer outputs.OutputWriter) error {

	// Context (with flags)
	ctx = context.WithValue(ctx,
		transform.ContextKeyDebug, o.Debug)

	// 🔁 Apply preprocessors here
	for _, p := range o.preprocessors {
		if err := p.Process(ctx, input); err != nil {
			return fmt.Errorf("preprocessor failed: %w", err)
		}
	}

	// ➕ Transform
	output, err := transform.New().Apply(ctx, input, o.Program)
	if err != nil {
		return fmt.Errorf("failed to apply program: %w", err)
	}
	if output == nil {
		// The program author determined this input should be disregarded/dropped
		return nil
	}

	// ➡️ Output
	if err := writer.WriteOutput(ctx, output); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}
