package orchestrator

import (
	"context"
	"errors"
	"github.com/tommed/ducto-dsl/transform"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrchestrator_RunLoop_Failure(t *testing.T) {
	ctx := context.Background()
	o := &Orchestrator{}
	writer := &FakeWriter{}

	source := NewErrorEventSource(errors.New("expected setup failure"))
	defer source.Close()

	err := o.RunLoop(ctx, source, writer)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected setup failure")
}

func TestOrchestrator_RunLoop_CancelledCtx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel context immediately

	source := NewSubjectEventSource()
	go func() {
		// Shouldn't need to wait as the context should immediately cancel underneath
		time.Sleep(5 * time.Second)
		source.Push(map[string]interface{}{})
	}()

	o := &Orchestrator{}
	err := o.RunLoop(ctx,
		source,
		&FakeWriter{})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "context canceled")
}

func TestOrchestrator_Execute(t *testing.T) {
	prog := &transform.Program{
		Version: 1,
		Instructions: []transform.Instruction{
			{Op: "set", Key: "greeting", Value: "hello"},
		},
	}

	input := map[string]interface{}{"foo": "bar"}
	source := NewValuesEventSource(input)
	writer := &FakeWriter{}

	err := New(prog, false).RunLoop(context.Background(), source, writer)

	assert.NoError(t, err)
	assert.Equal(t, "bar", writer.Written["foo"])
	assert.Equal(t, "hello", writer.Written["greeting"])
}
