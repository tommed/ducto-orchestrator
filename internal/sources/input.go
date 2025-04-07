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

type ValuesOptions struct {
	Values []map[string]interface{}
}

func NewValuesEventSource(opts ValuesOptions) EventSource {
	ch := make(chan map[string]interface{}, len(opts.Values))
	return &valuesEventSource{
		stream: ch,
		values: opts.Values,
	}
}

func NewValuesEventSourceRaw(values ...map[string]interface{}) EventSource {
	return NewValuesEventSource(ValuesOptions{Values: values})
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

func NewErrorEventSource(err error) EventSource {
	return &errorEventSource{err: err}
}

func (e *errorEventSource) Start(ctx context.Context) (<-chan map[string]interface{}, error) {
	return nil, e.err
}

func (e *errorEventSource) Close() error {
	return nil
}
