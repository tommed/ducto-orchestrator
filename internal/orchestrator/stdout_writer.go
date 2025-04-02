package orchestrator

import (
	"encoding/json"
	"io"
)

type StdoutWriter struct {
	Writer io.Writer
}

func (w *StdoutWriter) WriteOutput(data map[string]interface{}) error {
	encoder := json.NewEncoder(w.Writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
