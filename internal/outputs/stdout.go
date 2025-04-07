package outputs

import (
	"context"
	"encoding/json"
	"io"
)

type StdoutOptions struct {
	Pretty bool `mapstructure:"pretty"`
}

type stdoutWriter struct {
	writer io.Writer
	opts   StdoutOptions
}

func NewStdoutWriter(stdout io.Writer, opts StdoutOptions) OutputWriter {
	return &stdoutWriter{writer: stdout, opts: opts}
}

func (w *stdoutWriter) WriteOutput(_ context.Context, data map[string]interface{}) error {
	encoder := json.NewEncoder(w.writer)
	if w.opts.Pretty {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(data)
}
