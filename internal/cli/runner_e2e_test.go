package cli

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/tommed/ducto-dsl/transform"
	"github.com/tommed/ducto-orchestrator/internal/config"
	"github.com/tommed/ducto-orchestrator/internal/orchestrator"
	"github.com/tommed/ducto-orchestrator/internal/outputs"
	"github.com/tommed/ducto-orchestrator/internal/sources"
	"os"
	"testing"
	"time"
)

func TestOrchestrator_PubSub_E2E(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E tests in short mode")
	}
	if os.Getenv("PUBSUB_EMULATOR_HOST") == "" {
		t.Skip("PUBSUB_EMULATOR_HOST not set")
	}
	ctx := WithSignalContext(context.Background())
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	require.NotEmpty(t, projectID)

	client, err := pubsub.NewClient(ctx, projectID)
	require.NoError(t, err)

	topicID := "test-topic"
	subID := "test-sub"

	topic, err := client.CreateTopic(ctx, topicID)
	if err != nil {
		topic = client.Topic(topicID)
	}
	defer topic.Stop()

	sub, err := client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
		Topic: topic,
	})
	require.NoError(t, err)

	defer func(sub *pubsub.Subscription, ctx context.Context) {
		err := sub.Delete(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}(sub, ctx)

	// Configure the orchestrator
	cfg := &config.Config{
		Debug: true,
		Program: &transform.Program{
			Version:      1,
			OnError:      "ignore",
			Instructions: []transform.Instruction{}, // No-op, just pass-through
		},
		Source: config.PluginBlock{
			Type: "values",
			Config: map[string]interface{}{
				"values": []map[string]interface{}{
					{"message": "hello world"},
				},
			},
		},
		Output: config.PluginBlock{
			Type: "pubsub",
			Config: map[string]interface{}{
				"topic": fmt.Sprintf("projects/%s/topics/%s", projectID, topicID),
			},
		},
	}

	// Init orchestrator
	source, err := sources.FromPlugin(ctx, cfg.Source, nil)
	require.NoError(t, err)

	output, err := outputs.FromPlugin(ctx, cfg.Output, nil)
	require.NoError(t, err)

	// Run the orchestrator
	o := orchestrator.New(cfg.Program, cfg.Debug)
	require.NoError(t, o.RunLoop(ctx, source, output))

	// Pull the message back from subscription
	msgs := make(chan *pubsub.Message, 1)
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		msgs <- m
		m.Ack()
		cancel()
	})
	require.NoError(t, err)

	select {
	case msg := <-msgs:
		var result map[string]interface{}
		require.NoError(t, json.Unmarshal(msg.Data, &result))
		require.Equal(t, "hello world", result["message"])
	default:
		t.Fatal("No message received from pubsub")
	}
}
