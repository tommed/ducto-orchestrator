package sources

import (
	"context"
	"fmt"
	"io"

	"github.com/tommed/ducto-orchestrator/internal/config"
)

func FromPlugin(ctx context.Context, block config.PluginBlock, stdin io.Reader) (EventSource, error) {
	switch block.Type {
	case "stdin":
		return NewStdinEventSource(stdin), nil

	case "http":
		opts, err := config.Decode[HTTPOptions](block.Config)
		if err != nil {
			return nil, err
		}
		return NewHTTPEventSource(*opts), nil

	default:
		return nil, fmt.Errorf("unsupported source type: %q", block.Type)
	}
}
