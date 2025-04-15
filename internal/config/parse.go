package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tommed/ducto-dsl/transform"
	"gopkg.in/yaml.v3"
)

// ParseConfig reads and parses the config file without resolving program_file.
func ParseConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	switch {
	case strings.HasSuffix(path, ".json"):
		err = json.Unmarshal(data, &cfg)
	case strings.HasSuffix(path, ".yaml"), strings.HasSuffix(path, ".yml"):
		err = yaml.Unmarshal(data, &cfg)
	default:
		return nil, fmt.Errorf("unsupported config extension: %s", path)
	}
	if err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	return &cfg, nil
}

// FinalizeConfig resolves relative paths and loads the DSL program if needed.
func FinalizeConfig(cfg *Config, basedir string) error {
	SetConfigFilePath(basedir)

	if cfg.Program != nil {
		return nil
	}

	if cfg.ProgramFile != "" {
		progPath := cfg.ProgramFile
		if !filepath.IsAbs(progPath) {
			progPath = filepath.Join(basedir, progPath)
		}
		prog, err := transform.LoadProgram(progPath)
		if err != nil {
			return fmt.Errorf("load program: %w", err)
		}
		cfg.Program = prog
		cfg.ProgramFile = ""
		return nil
	}

	return fmt.Errorf("no program or program_file specified")
}
