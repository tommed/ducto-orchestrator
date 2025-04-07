package outputs

import (
	"context"
	"fmt"
	"io"

	"github.com/tommed/ducto-orchestrator/internal/config"
)

func FromPlugin(ctx context.Context, block config.PluginBlock, stdout io.Writer) (OutputWriter, error) {
	switch block.Type {
	case "stdout":
		opts, err := config.Decode[StdoutOptions](block.Config)
		if err != nil {
			return nil, err
		}
		return NewStdoutWriter(stdout, *opts), nil

	case "http":
		opts, err := config.Decode[HTTPOptions](block.Config)
		if err != nil {
			return nil, err
		}
		return NewHTTPWriter(*opts), nil

	case "pubsub":
		opts, err := config.Decode[PubSubOptions](block.Config)
		if err != nil {
			return nil, err
		}
		return NewPubSubWriter(ctx, *opts)

	default:
		return nil, fmt.Errorf("unsupported output type: %q", block.Type)
	}
}
