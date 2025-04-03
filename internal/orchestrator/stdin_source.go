package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type StdinEventSource struct {
	reader io.Reader
}

func NewStdinEventSource(stdin io.Reader) EventSource {
	return &StdinEventSource{
		reader: stdin,
	}
}

// Start reads exactly one JSON object from stdin and then closes the stream.
func (s *StdinEventSource) Start(ctx context.Context) (<-chan map[string]interface{}, error) {
	// Read input
	var input map[string]interface{}
	decoder := json.NewDecoder(s.reader)
	if err := decoder.Decode(&input); err != nil {
		return nil, fmt.Errorf("failed to decode stdin input: %w", err)
	}

	// Delegate to ValuesEventSource without exporting it
	return NewValuesEventSource(input).Start(ctx)
}

func (s *StdinEventSource) Close() error {
	if c, ok := s.reader.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
