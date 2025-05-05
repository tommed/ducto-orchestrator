package sources

import (
	"context"
	"fmt"
	"io"

	"github.com/tommed/ducto-orchestrator/internal/config"
)

func FromPlugin(_ context.Context, block config.PluginBlock, stdin io.Reader) (EventSource, error) {
	switch block.Type {
	case "stdin":
		opts, err := config.Decode[StdinOptions](block.Config)
		if err != nil {
			return nil, err
		}
		return NewStdinEventSource(stdin, *opts), nil

	case "values":
		opts, err := config.Decode[ValuesOptions](block.Config)
		if err != nil {
			return nil, err
		}
		return NewValuesEventSource(*opts), nil

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
