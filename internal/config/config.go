package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/tommed/ducto-dsl/transform"
	"path/filepath"
)

type PluginBlock struct {
	Type   string                 `json:"type" yaml:"type" mapstructure:"type"`
	Config map[string]interface{} `json:"config" yaml:"config" mapstructure:"config"`
}

type Config struct {
	Debug       bool               `json:"debug" yaml:"debug" mapstructure:"debug"`
	Program     *transform.Program `json:"program" yaml:"program" mapstructure:"program"`
	ProgramFile string             `json:"program_file" yaml:"program_file" mapstructure:"program_file"`

	Preprocessors []PluginBlock `json:"preprocessors" yaml:"preprocessors" mapstructure:"preprocessors"`
	Source        PluginBlock   `json:"source" yaml:"source" mapstructure:"source"`
	Output        PluginBlock   `json:"output" yaml:"output" mapstructure:"output"`
}

var configDir string

func SetConfigFilePath(absPath string) {
	configDir = filepath.Dir(absPath)
}

func ResolvePath(relOrAbs string) (string, error) {
	if filepath.IsAbs(relOrAbs) {
		return relOrAbs, nil
	}
	if configDir == "" {
		return "", fmt.Errorf("no config directory set")
	}
	return filepath.Join(configDir, relOrAbs), nil
}

type Options interface {
	Validate() error
}

func Decode[T any](raw map[string]interface{}) (*T, error) {
	var target T
	dec, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &target,
		TagName: "mapstructure",
	})
	// Not sure if it's possible with mapstructure to return an error here?!
	//if err != nil {
	//	return nil, err
	//}
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
