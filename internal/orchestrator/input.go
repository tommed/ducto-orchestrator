package orchestrator

import "context"

type EventSource interface {
	Start(ctx context.Context) (<-chan map[string]interface{}, <-chan error)
	Close() error
}

type FakeEventSource struct {
	events chan map[string]interface{}
}

func NewFakeEventSource(events ...map[string]interface{}) EventSource {
	ch := make(chan map[string]interface{}, len(events))
	for _, e := range events {
		ch <- e
	}
	close(ch)
	return &FakeEventSource{
		events: ch,
	}
}

func (f *FakeEventSource) Start(_ context.Context) (<-chan map[string]interface{}, <-chan error) {
	return f.events, nil
}

func (f *FakeEventSource) Close() error {
	return nil
}
