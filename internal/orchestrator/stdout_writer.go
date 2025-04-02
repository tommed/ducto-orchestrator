package orchestrator

import (
	"encoding/json"
	"io"
)

type StdoutWriter struct {
	writer io.Writer
}

func NewStdoutWriter(stdout io.Writer) *StdoutWriter {
	return &StdoutWriter{writer: stdout}
}

func (w *StdoutWriter) WriteOutput(data map[string]interface{}) error {
	encoder := json.NewEncoder(w.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
