package config

import "context"

// Loader defines an interface for loading Ducto configuration files.
type Loader interface {
	Load(ctx context.Context, uri string) (*Config, error)
}

//// WatchingLoader interface for loaders that support watching for changes
//type WatchingLoader interface {
//	Loader
//	Watch(ctx context.Context, uri string) <-chan struct{}
//}
