package orchestrator

import (
	"context"
	"os/signal"
	"syscall"
)

func WithSignalContext(parent context.Context) context.Context {
	ctx, stop := signal.NotifyContext(parent, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-ctx.Done()
		stop() // clean up signal handlers
	}()

	return ctx
}
