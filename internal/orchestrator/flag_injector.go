package orchestrator

import (
	"context"
	flagsdk "github.com/tommed/ducto-featureflags/sdk"
)

type FlagInjectorOptions struct {
	File  string                  `mapstructure:"file,omitempty"`
	Flags map[string]flagsdk.Flag `mapstructure:"flags,omitempty"`
	Keys  []string                `mapstructure:"keys"`
}

type FlagInjector struct {
	store *flagsdk.Store
	keys  []string
}

func NewFlagInjector(store *flagsdk.Store, keys []string) *FlagInjector {
	return &FlagInjector{store: store, keys: keys}
}

func (f *FlagInjector) Process(_ context.Context, input map[string]interface{}) error {
	ctx := flagsdk.EvalContext{}

	for k, v := range input {
		if str, ok := v.(string); ok {
			ctx[k] = str
		}
	}

	res := map[string]bool{}
	for _, key := range f.keys {
		res[key] = f.store.IsEnabled(key, ctx)
	}

	input["_flags"] = res
	return nil
}
