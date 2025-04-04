package writers

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStdoutWriter(t *testing.T) {
	stdout := &bytes.Buffer{}
	event := map[string]interface{}{"foo": "bar"}
	expected, _ := json.MarshalIndent(event, "", "  ")
	output := NewStdoutWriter(stdout, StdoutOptions{Pretty: true})

	// Act
	err := output.WriteOutput(event)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t,
		strings.TrimSpace(string(expected)),
		strings.TrimSpace(stdout.String()))
}
