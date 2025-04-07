package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/tommed/ducto-dsl/transform"
	"path/filepath"
)

type PluginBlock struct {
	Type   string                 `mapstructure:"type"`
	Config map[string]interface{} `mapstructure:"config"`
}

type Config struct {
	Debug       bool               `mapstructure:"debug"`
	Program     *transform.Program `mapstructure:"program"`
	ProgramFile string             `mapstructure:"program_file"`

	Source PluginBlock `mapstructure:"source"`
	Output PluginBlock `mapstructure:"output"`
}

func Load(path string) (*Config, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("get abs path: %w", err)
	}

	v := viper.New()
	v.SetConfigFile(absPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Change relative paths so they are relative to this config file NOT the current working directory.
	// This is the behaviour of the least surprise.
	cfgDir := filepath.Dir(absPath)
	if cfg.ProgramFile != "" && !filepath.IsAbs(cfg.ProgramFile) {
		cfg.ProgramFile = filepath.Join(cfgDir, cfg.ProgramFile)
	}

	return &cfg, nil
}

type Options interface {
	Validate() error
}

func Decode[T any](raw map[string]interface{}) (*T, error) {
	var target T
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &target,
		TagName: "mapstructure",
	})
	if err != nil {
		return nil, err
	}
	if err := dec.Decode(raw); err != nil {
		return nil, err
	}
	if val, ok := any(&target).(Options); ok {
		if err := val.Validate(); err != nil {
			return nil, err
		}
	}
	return &target, nil
}
