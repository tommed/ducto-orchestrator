package orchestrator

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func WithSignalContext(parent context.Context) context.Context {
	return withSignalContextRaw(parent, syscall.SIGINT, syscall.SIGTERM)
}

func withSignalContextRaw(parent context.Context, signals ...os.Signal) context.Context {
	ctx, stop := signal.NotifyContext(parent, signals...)

	go func() {
		<-ctx.Done()
		stop() // clean up signal handlers
	}()

	return ctx
}
