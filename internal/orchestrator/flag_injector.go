package orchestrator

import (
	"context"
	flagsdk "github.com/tommed/ducto-featureflags/sdk"
	"os"
	"strings"
)

type FlagInjectorOptions struct {
	RawTags  flagsdk.EvalContext     `json:"tags,omitempty" mapstructure:"tags,omitempty"`
	TagsEnv  string                  `json:"tags_env,omitempty" mapstructure:"tags_env,omitempty"`
	Flags    map[string]flagsdk.Flag `json:"flags,omitempty" mapstructure:"flags,omitempty"`
	Provider map[string]interface{}  `json:"provider,omitempty" mapstructure:"provider,omitempty"`
}

func (o FlagInjectorOptions) Tags() flagsdk.EvalContext {
	if o.TagsEnv != "" {
		return parseTagsEnv(os.Getenv(o.TagsEnv))
	}
	return o.RawTags
}

type FlagInjector struct {
	store flagsdk.AnyStore
	tags  flagsdk.EvalContext
}

func NewFlagInjector(store flagsdk.AnyStore, tags flagsdk.EvalContext) *FlagInjector {
	return &FlagInjector{store: store, tags: tags}
}

func (f *FlagInjector) Process(_ context.Context, input map[string]interface{}) error {
	res := make(map[string]interface{})
	for key := range f.store.AllFlags() {
		flag, ok := f.store.Get(key)
		if ok {
			result := flag.Evaluate(f.tags)
			if result.OK {
				res[key] = result.Value
			}
		}
	}

	input["_flags"] = res
	return nil
}

func parseTagsEnv(env string) flagsdk.EvalContext {
	out := flagsdk.EvalContext{}
	for _, pair := range strings.Split(env, ";") {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			out[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return out
}
