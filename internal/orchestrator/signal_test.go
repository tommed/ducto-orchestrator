package orchestrator

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWithSignalContext(t *testing.T) {
	// Use SIGUSR1 as it's user-defined and safe for testing
	ctx := context.Background()
	ctx = withSignalContextRaw(ctx, syscall.SIGUSR1)

	// Send the signal to ourselves
	p, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)
	err = p.Signal(syscall.SIGUSR1)
	require.NoError(t, err)

	select {
	case <-ctx.Done():
		// success
	case <-time.After(1 * time.Second):
		t.Fatal("context was not cancelled by signal")
	}
}
