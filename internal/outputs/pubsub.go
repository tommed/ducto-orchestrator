package outputs

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// PubSubOptions configure the Pub/Sub writer.
type PubSubOptions struct {
	Topic                string        `mapstructure:"topic"`
	AttributesFromFields []string      `mapstructure:"attributes_from_fields"`
	OrderingKeyField     string        `mapstructure:"ordering_key_field"`
	EnableOrdering       bool          `mapstructure:"enable_message_ordering"`
	Timeout              time.Duration `mapstructure:"timeout"` // defaults to 5s if unset
}

func (o *PubSubOptions) Validate() error {
	if o.Topic == "" {
		return fmt.Errorf("topic required")
	}
	if o.Timeout == 0 {
		o.Timeout = 30 * time.Second
	}
	if o.Timeout < 0 {
		return fmt.Errorf("invalid timeout")
	}
	if o.Timeout > 5*time.Minute {
		return fmt.Errorf("timeout too high")
	}
	if o.EnableOrdering && o.OrderingKeyField == "" {
		return fmt.Errorf("ordering_key_field required when enable_message_ordering set")
	}
	parts := strings.Split(o.Topic, "/")
	if len(parts) < 4 || parts[0] != "projects" || parts[2] != "topics" {
		return fmt.Errorf("invalid topic format")
	}
	return nil
}

func NewPubSubWriter(ctx context.Context, opts PubSubOptions) (OutputWriter, error) {
	parts := strings.Split(opts.Topic, "/")
	projectID := parts[1]
	topicID := parts[3]
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %w", err)
	}

	topic := client.Topic(topicID)
	topic.EnableMessageOrdering = opts.EnableOrdering

	if opts.Timeout == 0 {
		opts.Timeout = 5 * time.Second
	}

	return &pubSubWriter{
		client: client,
		topic:  topic,
		opts:   opts,
	}, nil
}

// pubSubWriter is our concrete pub/sub writer class
type pubSubWriter struct {
	client *pubsub.Client
	topic  *pubsub.Topic
	opts   PubSubOptions
}

func (w pubSubWriter) WriteOutput(ctx context.Context, input map[string]interface{}) error {
	ctxCncl, cancel := context.WithTimeout(ctx, w.opts.Timeout)
	defer cancel()

	data, err := json.Marshal(input)
	if err != nil {
		return err
	}

	msg := &pubsub.Message{
		Data: data,
	}

	// Extract attributes
	if w.opts.AttributesFromFields != nil {
		msg.Attributes = map[string]string{}
		for _, field := range w.opts.AttributesFromFields {
			if val, ok := input[field]; ok {
				msg.Attributes[field] = fmt.Sprintf("%v", val)
			}
		}
	}

	// Ordering Key
	if w.opts.OrderingKeyField != "" {
		if val, ok := input[w.opts.OrderingKeyField]; ok {
			msg.OrderingKey = fmt.Sprintf("%v", val)
		}
	}

	// Publish
	res := w.topic.Publish(ctxCncl, msg)
	_, err = res.Get(ctxCncl)
	return err
}
