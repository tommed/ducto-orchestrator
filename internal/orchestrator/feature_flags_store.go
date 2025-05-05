package orchestrator

import (
	"context"
	"errors"
	"fmt"
	flagsdk "github.com/tommed/ducto-featureflags/sdk"
	"github.com/tommed/ducto-orchestrator/internal/config"
	"os"
	"time"
)

// NewStoreFromConfig only needs to cover custom (usually remote) feature flag sources.
// Local file and embedded flags are already supported at a higher level
func NewStoreFromConfig(ctx context.Context, raw map[string]interface{}) (*flagsdk.DynamicStore, error) {
	sourceType, ok := raw["type"].(string)
	if !ok {
		return nil, fmt.Errorf("missing orchestrator.feature_flags.source.type")
	}

	var provider flagsdk.StoreProvider
	switch sourceType {

	// Watch a local file
	case "file":
		path, ok := raw["file"].(string)
		if !ok {
			return nil, fmt.Errorf("missing orchestrator.feature_flags.source.file.path")
		}
		path, _ = config.ResolvePath(path)
		provider = flagsdk.NewFileProvider(path)

	// HTTP Watcher with 304 smart re-loading
	case "http":
		opts, err := config.Decode[flagsHTTPSourceOptions](raw)
		if err != nil {
			return nil, fmt.Errorf("invalid http source config: %w", err)
		}
		if err := opts.Validate(); err != nil {
			return nil, fmt.Errorf("invalid http source config: %w", err)
		}
		provider = flagsdk.NewHTTPProvider(opts.URL, opts.Token(), opts.PollInterval())

	// Fail otherwise
	default:
		return nil, fmt.Errorf("unsupported orchestrator.feature_flags.source.type: %s", sourceType)
	}

	// Run this from inside a DynamicStore
	return flagsdk.NewDynamicStore(ctx, provider), nil
}

type flagsHTTPSourceOptions struct {
	URL                 string `json:"url" mapstructure:"url"`
	TokenLiteral        string `json:"token" mapstructure:"token"`
	TokenEnv            string `json:"token_env" mapstructure:"token_env"`
	PollIntervalSeconds int    `json:"poll_interval_seconds" mapstructure:"poll_interval_seconds"`
}

func (o *flagsHTTPSourceOptions) Token() string {
	if o.TokenEnv != "" {
		return os.Getenv(o.TokenEnv)
	}
	return o.TokenLiteral
}

func (o *flagsHTTPSourceOptions) PollInterval() time.Duration {
	return time.Duration(o.PollIntervalSeconds) * time.Second
}

func (o *flagsHTTPSourceOptions) Validate() error {
	if o.URL == "" {
		return errors.New("missing url")
	}
	if o.PollInterval() < time.Second {
		return errors.New("poll_interval too small")
	}
	return nil
}
