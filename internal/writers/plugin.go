package writers

import (
	"fmt"
	"io"

	"github.com/tommed/ducto-orchestrator/internal/config"
)

func FromPlugin(block config.PluginBlock, stdout io.Writer) (OutputWriter, error) {
	switch block.Type {
	case "stdout":
		opts, err := config.Decode[StdoutOptions](block.Config)
		if err != nil {
			return nil, err
		}
		return NewStdoutWriter(stdout, *opts), nil

	default:
		return nil, fmt.Errorf("unsupported output type: %q", block.Type)
	}
}
