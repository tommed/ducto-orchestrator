package orchestrator

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"github.com/tommed/ducto-dsl/transform"
	"github.com/tommed/ducto-orchestrator/internal/outputs"
	"github.com/tommed/ducto-orchestrator/internal/sources"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//goland:noinspection GoUnhandledErrorResult
func TestOrchestrator_RunLoop_Failure(t *testing.T) {
	type args struct {
		source sources.EventSource
		output outputs.OutputWriter
		op     transform.Instruction
	}
	tests := []struct {
		name      string
		args      args
		wantInErr string
	}{
		{
			name: "failing source",
			args: args{
				source: sources.NewErrorEventSource(errors.New("expected setup failure")),
				output: &outputs.FakeWriter{},
				op:     transform.Instruction{Op: "noop"},
			},
			wantInErr: "expected setup failure",
		},
		{
			name: "failing source",
			args: args{
				source: sources.NewValuesEventSource(map[string]interface{}{}),
				output: outputs.NewFailingWriter(errors.New("expected output failure")),
				op:     transform.Instruction{Op: "noop"},
			},
			wantInErr: "expected output failure",
		},
		{
			name: "failing program",
			args: args{
				source: sources.NewValuesEventSource(map[string]interface{}{}),
				output: &outputs.FakeWriter{},
				op:     transform.Instruction{Op: "fail", Value: "expected operation failure"},
			},
			wantInErr: "execution halted due to an error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			o := &Orchestrator{
				Program: &transform.Program{
					Version: 1,
					OnError: "fail",
					Instructions: []transform.Instruction{
						tt.args.op,
					},
				},
			}
			defer tt.args.source.Close()

			err := o.RunLoop(ctx, tt.args.source, tt.args.output)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantInErr)
		})
	}
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
		&outputs.FakeWriter{})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "context canceled")
}

func TestOrchestrator_RunLoop_Success(t *testing.T) {
	prog := &transform.Program{
		Version: 1,
		Instructions: []transform.Instruction{
			{Op: "set", Key: "greeting", Value: "hello"},
		},
	}

	input := map[string]interface{}{"foo": "bar"}
	source := sources.NewValuesEventSource(input)
	writer := &outputs.FakeWriter{}

	err := New(prog, false).RunLoop(context.Background(), source, writer)

	assert.NoError(t, err)
	assert.Equal(t, "bar", writer.Written["foo"])
	assert.Equal(t, "hello", writer.Written["greeting"])
}
