package orchestrator

import (
	"bytes"
	"context"
)

func buf() *bytes.Buffer {
	return &bytes.Buffer{}
}

// SubjectEventSource has a Push event for manually sending events through the
// channel. This should only really be used for unit testing - otherwise consider
// is a compromise on code - a code smell.
type SubjectEventSource struct {
	channel chan map[string]interface{}
}

func NewSubjectEventSource() *SubjectEventSource {
	return &SubjectEventSource{channel: make(chan map[string]interface{})}
}

func (s *SubjectEventSource) Push(inputs ...map[string]interface{}) {
	for _, input := range inputs {
		s.channel <- input
	}
}

func (s *SubjectEventSource) Start(ctx context.Context) (<-chan map[string]interface{}, error) {
	return s.channel, nil
}

func (s *SubjectEventSource) Close() error {
	if s.channel != nil {
		close(s.channel)
		s.channel = nil
	}
	return nil
}
