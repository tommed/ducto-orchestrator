package outputs

import (
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

func (w *stdoutWriter) WriteOutput(data map[string]interface{}) error {
	encoder := json.NewEncoder(w.writer)
	if w.opts.Pretty {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(data)
}
