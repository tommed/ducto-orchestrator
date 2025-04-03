package orchestrator

import (
	"encoding/json"
	"io"
)

type stdoutWriter struct {
	writer io.Writer
}

func NewStdoutWriter(stdout io.Writer) OutputWriter {
	return &stdoutWriter{writer: stdout}
}

func (w *stdoutWriter) WriteOutput(data map[string]interface{}) error {
	encoder := json.NewEncoder(w.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
