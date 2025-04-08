package orchestrator

import (
	"fmt"
	"github.com/tommed/ducto-orchestrator/internal/config"

	flagsdk "github.com/tommed/ducto-featureflags/sdk"
)

func (o *Orchestrator) InstallPreprocessors(preprocessors []config.PluginBlock) error {
	for _, block := range preprocessors {
		switch block.Type {
		case "feature_flags":
			p, err := NewFlagInjectorFromConfig(block.Config)
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

func NewFlagInjectorFromConfig(raw map[string]interface{}) (*FlagInjector, error) {
	opts, err := config.Decode[FlagInjectorOptions](raw)
	if err != nil {
		return nil, err
	}

	var store *flagsdk.Store

	switch {
	case opts.File != "":
		path, err := config.ResolvePath(opts.File)
		if err != nil {
			return nil, fmt.Errorf("resolve path: %w", err)
		}
		store, err = flagsdk.NewStoreFromFile(path)
		if err != nil {
			return nil, fmt.Errorf("load flag file: %w", err)
		}
	case opts.Flags != nil:
		store = flagsdk.NewStore(opts.Flags)
	default:
		return nil, fmt.Errorf("featureflags preprocessor requires either 'file' or inline 'flags'")
	}

	return NewFlagInjector(store, opts.Keys), nil
}
