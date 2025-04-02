package orchestrator

import (
	"context"
	"github.com/tommed/ducto-dsl/transform"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrchestrator_Execute(t *testing.T) {
	prog := &transform.Program{
		Version: 1,
		Instructions: []transform.Instruction{
			{Op: "set", Key: "greeting", Value: "hello"},
		},
	}

	o := New(prog, false)

	input := map[string]interface{}{"foo": "bar"}

	reader := &FakeReader{Data: input}
	writer := &FakeWriter{}

	err := o.Execute(context.Background(), reader, writer)

	assert.NoError(t, err)
	assert.Equal(t, "bar", writer.Written["foo"])
	assert.Equal(t, "hello", writer.Written["greeting"])
}
