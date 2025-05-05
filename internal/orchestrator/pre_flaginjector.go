package orchestrator

import (
	"context"
	"fmt"
	flagsdk "github.com/tommed/ducto-featureflags/sdk"
	"github.com/tommed/ducto-orchestrator/internal/config"
)

func NewFlagInjectorFromConfig(ctx context.Context, raw map[string]interface{}) (*FlagInjector, error) {
	opts, err := config.Decode[FlagInjectorOptions](raw)
	if err != nil {
		return nil, err
	}

	var store flagsdk.AnyStore
	switch {

	// 'Provider' is a dynamic loader of feature flags
	// For example:  NewFileProvider, or NewHTTPProvider
	case opts.Provider != nil:
		reloadingStore, err := NewStoreFromConfig(ctx, opts.Provider)
		if err != nil {
			return nil, err
		}
		err = reloadingStore.Start()
		if err != nil {
			return nil, fmt.Errorf("load flag preprocessor: %w", err)
		}
		store = reloadingStore

	// Hard-coded flags (inline)
	case opts.Flags != nil:
		store = flagsdk.NewStore(opts.Flags)

	// Not supported
	default:
		return nil, fmt.Errorf("featureflags preprocessor requires either 'file' or inline 'flags'")
	}

	return NewFlagInjector(store, opts.Tags()), nil
}
