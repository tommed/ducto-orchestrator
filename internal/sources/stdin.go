package sources

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type StdinOptions struct {
	JSONL bool `json:"jsonl" mapstructure:"jsonl"`
}

type stdinEventSource struct {
	reader   io.Reader
	useJSONL bool
}

func NewStdinEventSource(stdin io.Reader, opts StdinOptions) EventSource {
	return &stdinEventSource{
		reader:   stdin,
		useJSONL: opts.JSONL,
	}
}

// Start reads from stdin either as a single JSON object or as a stream of JSONL events.
func (s *stdinEventSource) Start(ctx context.Context) (<-chan map[string]interface{}, error) {
	if s.useJSONL {
		return s.readJSONLStream(ctx)
	} else {
		return s.readSingleJSON(ctx)
	}
}

func (s *stdinEventSource) readSingleJSON(ctx context.Context) (<-chan map[string]interface{}, error) {
	var input map[string]interface{}
	decoder := json.NewDecoder(s.reader)
	if err := decoder.Decode(&input); err != nil {
		return nil, fmt.Errorf("failed to decode stdin input: %w", err)
	}
	return NewValuesEventSourceRaw(input).Start(ctx)
}

func (s *stdinEventSource) readJSONLStream(ctx context.Context) (<-chan map[string]interface{}, error) {
	ch := make(chan map[string]interface{})
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(s.reader)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
			}

			var obj map[string]interface{}
			line := scanner.Bytes()
			if err := json.Unmarshal(line, &obj); err != nil {
				// skip invalid lines but consider logging
				continue
			}

			ch <- obj
		}
	}()
	return ch, nil
}

func (s *stdinEventSource) Close() error {
	if c, ok := s.reader.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
