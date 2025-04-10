package orchestrator

import (
	"context"
	flagsdk "github.com/tommed/ducto-featureflags/sdk"
)

type FlagInjectorOptions struct {
	Tags     flagsdk.EvalContext     `json:"tags" mapstructure:"tags"`
	Flags    map[string]flagsdk.Flag `json:"flags,omitempty" mapstructure:"flags,omitempty"`
	Provider map[string]interface{}  `json:"provider,omitempty" mapstructure:"provider,omitempty"`
}

type FlagInjector struct {
	store flagsdk.AnyStore
	tags  flagsdk.EvalContext
}

func NewFlagInjector(store flagsdk.AnyStore, tags flagsdk.EvalContext) *FlagInjector {
	return &FlagInjector{store: store, tags: tags}
}

func (f *FlagInjector) Process(_ context.Context, input map[string]interface{}) error {
	res := map[string]bool{}
	for key := range f.store.AllFlags() {
		res[key] = f.store.IsEnabled(key, f.tags)
	}

	input["_flags"] = res
	return nil
}
