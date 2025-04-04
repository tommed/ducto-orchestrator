package sources

import "context"

type EventSource interface {
	Start(ctx context.Context) (<-chan map[string]interface{}, error)
	Close() error
}

type valuesEventSource struct {
	stream chan map[string]interface{}
	values []map[string]interface{}
}

func NewValuesEventSource(values ...map[string]interface{}) EventSource {
	ch := make(chan map[string]interface{}, len(values))
	return &valuesEventSource{
		stream: ch,
		values: values,
	}
}

func (f *valuesEventSource) Start(_ context.Context) (<-chan map[string]interface{}, error) {
	for _, v := range f.values {
		f.stream <- v
	}
	close(f.stream)
	return f.stream, nil
}

func (f *valuesEventSource) Close() error {
	return nil // nothing to close
}

type errorEventSource struct {
	err error
}

func NewErrorEventSource(err error) *errorEventSource {
	return &errorEventSource{err: err}
}

func (e *errorEventSource) Start(ctx context.Context) (<-chan map[string]interface{}, error) {
	return nil, e.err
}

func (e *errorEventSource) Close() error {
	return nil
}
