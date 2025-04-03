package orchestrator

import "context"

type EventSource interface {
	Start(ctx context.Context) (<-chan map[string]interface{}, error)
	Close() error
}

type ValuesEventSource struct {
	stream chan map[string]interface{}
	values []map[string]interface{}
}

func NewValuesEventSource(values ...map[string]interface{}) EventSource {
	ch := make(chan map[string]interface{}, len(values))
	return &ValuesEventSource{
		stream: ch,
		values: values,
	}
}

func (f *ValuesEventSource) Start(_ context.Context) (<-chan map[string]interface{}, error) {
	for _, v := range f.values {
		f.stream <- v
	}
	close(f.stream)
	return f.stream, nil
}

func (f *ValuesEventSource) Close() error {
	return nil // nothing to close
}
