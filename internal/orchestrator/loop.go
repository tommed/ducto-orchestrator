package orchestrator

import (
	"context"
)

func (o *Orchestrator) RunLoop(ctx context.Context, source EventSource, writer OutputWriter) error {

	// Setup Source
	events, errc := source.Start(ctx)

	// Teardown Source
	defer func(source EventSource) {
		_ = source.Close()
	}(source)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case err, ok := <-errc:
			if ok && err != nil {
				return err // fail fast
			}

		case evt, ok := <-events:
			if !ok {
				return nil // channel closed naturally
			}

			if err := o.Execute(ctx, evt, writer); err != nil {
				return err
			}
		}
	}
}
