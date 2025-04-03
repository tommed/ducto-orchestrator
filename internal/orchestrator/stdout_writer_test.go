package orchestrator

import (
	"context"
	"encoding/json"
	"github.com/tommed/ducto-dsl/transform"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStdoutWriter(t *testing.T) {
	stdout := buf()
	event := map[string]interface{}{"foo": "bar"}
	input := NewValuesEventSource(event)
	output := NewStdoutWriter(stdout)
	prog := &transform.Program{
		Version:      1,
		Instructions: []transform.Instruction{{Op: "noop"}},
	}
	o := New(prog, false)
	expected, _ := json.MarshalIndent(event, "", "  ")

	// Act
	err := o.RunLoop(context.Background(), input, output)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t,
		strings.TrimSpace(string(expected)),
		strings.TrimSpace(stdout.String()))
}
