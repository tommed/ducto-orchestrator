package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tommed/ducto-dsl/transform"
	"gopkg.in/yaml.v3"
)

type FileLoader struct{}

// Load implements Loader. It accepts a path to a local YAML or JSON file.
func (f *FileLoader) Load(_ context.Context, uri string) (*Config, error) {
	// Accept either bare path or file:// URI
	path := strings.TrimPrefix(uri, "file://")
	if path == "" {
		return nil, fmt.Errorf("fileloader: missing config file path")
	}

	return LoadFromPath(path)
}

// LoadFromPath loads the config from a local file, resolves program paths, and stores the source path.
func LoadFromPath(path string) (*Config, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("get abs path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	switch {
	case strings.HasSuffix(absPath, ".json"):
		err = json.Unmarshal(data, &cfg)
	case strings.HasSuffix(absPath, ".yaml"), strings.HasSuffix(absPath, ".yml"):
		err = yaml.Unmarshal(data, &cfg)
	default:
		return nil, fmt.Errorf("unsupported config file extension: %s", absPath)
	}
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	cfgDir := filepath.Dir(absPath)
	if cfg.ProgramFile != "" && !filepath.IsAbs(cfg.ProgramFile) {
		cfg.ProgramFile = filepath.Join(cfgDir, cfg.ProgramFile)
	}

	// Resolve program if necessary
	switch {
	case cfg.Program != nil:
		// already embedded
	case cfg.ProgramFile != "":
		prog, err := transform.LoadProgram(cfg.ProgramFile)
		if err != nil {
			return nil, fmt.Errorf("load program file: %w", err)
		}
		cfg.Program = prog
		cfg.ProgramFile = ""
	default:
		return nil, fmt.Errorf("no program defined in config")
	}

	SetConfigFilePath(absPath)
	return &cfg, nil
}
