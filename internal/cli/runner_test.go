package cli_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/tommed/ducto-orchestrator/internal/cli"
	"io"
	"testing"
)

func TestRun_CoreCases(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		configFile     string // optional
		stdin          string // optional
		expectCode     int
		expectInStdErr string
		expectInStdOut string
	}{
		{
			name:           "invalid flags",
			args:           []string{"-doesnt -exist"},
			expectCode:     1,
			expectInStdErr: "failed to parse args",
		},
		{
			name:           "missing --config",
			args:           []string{},
			expectCode:     1,
			expectInStdErr: "missing required --config",
		},
		{
			name:           "nonexistent config",
			args:           []string{"--config", "../../testdata/does-not-exist.yaml"},
			expectCode:     1,
			expectInStdErr: "failed to load config",
		},
		{
			name:           "no program in config",
			args:           []string{"--config", "../../testdata/no_program.yaml"},
			expectCode:     1,
			expectInStdErr: "no DSL program or program_file defined",
		},
		{
			name:           "invalid program path",
			args:           []string{"--config", "../../testdata/invalid_program_file.yaml"},
			expectCode:     1,
			expectInStdErr: "failed to read program",
		},
		{
			name:           "invalid source",
			args:           []string{"--config", "../../testdata/invalid_source.yaml"},
			expectCode:     1,
			expectInStdErr: "failed to load source",
		},
		{
			name:           "invalid output",
			args:           []string{"--config", "../../testdata/invalid_output.yaml"},
			expectCode:     1,
			expectInStdErr: "failed to load output",
		},
		{
			name:           "program fails",
			args:           []string{"--config", "../../testdata/valid_but_program_fails.yaml"},
			expectCode:     1,
			expectInStdErr: "orchestrator failed",
		},
		{
			name:           "valid embedded program",
			args:           []string{"--debug", "--config", "../../examples/embedded-program-stdin.yaml"},
			stdin:          `{"foo":"bar"}`,
			expectCode:     0,
			expectInStdOut: `"greeting": "hello world"`,
		},
		{
			name:           "valid config + stdin + transform",
			args:           []string{"--config", "../../examples/simplest.yaml"},
			stdin:          `{"foo":"bar"}`,
			expectCode:     0,
			expectInStdOut: `"greeting": "hello world"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdin := io.NopCloser(bytes.NewReader([]byte(tt.stdin)))
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			code := cli.Run(tt.args, stdin, stdout, stderr)

			assert.Equal(t, tt.expectCode, code)

			if tt.expectInStdErr != "" {
				assert.Contains(t, stderr.String(), tt.expectInStdErr)
			}
			if tt.expectInStdOut != "" {
				assert.Contains(t, stdout.String(), tt.expectInStdOut)
			}
		})
	}
}
