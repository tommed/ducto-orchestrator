package config

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

type FileLoader struct{}

// Load implements Loader. It accepts a path to a local YAML or JSON file.
func (f *FileLoader) Load(_ context.Context, uri string) (*Config, error) {
	if uri == "" {
		return nil, fmt.Errorf("fileloader: missing config file path")
	}

	// Accept either bare path or file:// URI
	path := strings.TrimPrefix(uri, "file://")
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("get abs path: %w", err)
	}

	cfg, err := ParseConfig(absPath)
	if err != nil {
		return nil, err
	}

	if err := FinalizeConfig(cfg, filepath.Dir(absPath)); err != nil {
		return nil, err
	}
	
	return cfg, nil
}
