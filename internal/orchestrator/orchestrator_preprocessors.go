package orchestrator

import (
	"context"
	"fmt"
	"github.com/tommed/ducto-orchestrator/internal/config"
)

func (o *Orchestrator) InstallPreprocessors(ctx context.Context, preprocessors []config.PluginBlock) error {
	for _, block := range preprocessors {
		switch block.Type {
		case "feature_flags":
			p, err := NewFlagInjectorFromConfig(ctx, block.Config)
			if err != nil {
				return fmt.Errorf("load feature flag preprocessor: %w", err)
			}
			o.AddPreprocessor(p)
		default:
			return fmt.Errorf("unknown preprocessor type: %q", block.Type)
		}
	}
	return nil
}
