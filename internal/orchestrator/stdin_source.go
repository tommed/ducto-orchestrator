package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type StdinEventSource struct {
	reader io.Reader
	events chan map[string]interface{}
}

func NewStdinEventSource(stdin io.Reader) EventSource {
	return &StdinEventSource{
		reader: stdin,
		events: make(chan map[string]interface{}),
	}
}

// Start reads exactly one JSON object from stdin and then closes the channel.
func (s *StdinEventSource) Start(_ context.Context) (<-chan map[string]interface{}, <-chan error) {
	errc := make(chan error, 1) // always buffered

	go func() {
		defer close(s.events)
		defer close(errc)

		var input map[string]interface{}
		decoder := json.NewDecoder(s.reader)
		if err := decoder.Decode(&input); err != nil {
			errc <- fmt.Errorf("failed to decode stdin input: %w", err)
			return
		}

		s.events <- input
	}()

	return s.events, errc
}

func (s *StdinEventSource) Close() error {
	return nil
}
