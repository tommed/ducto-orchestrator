// File: ducto-orchestrator/internal/config/fileloader_test.go
package config_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tommed/ducto-orchestrator/internal/config"
)

func TestFileLoader_JSONConfig(t *testing.T) {
	temp := t.TempDir()
	file := filepath.Join(temp, "config.json")
	yaml := `{
		"program": {
			"version": 1,
			"instructions": [
				{"op": "set", "key": "x", "value": 1}
			]
		}
	}`
	assert.NoError(t, os.WriteFile(file, []byte(yaml), 0644))

	loader := &config.FileLoader{}
	cfg, err := loader.Load(context.Background(), file)
	assert.NoError(t, err)
	assert.NotNil(t, cfg.Program)
}

func TestFileLoader_UnsupportedExtension(t *testing.T) {
	temp := t.TempDir()
	file := filepath.Join(temp, "config.unsupported")
	assert.NoError(t, os.WriteFile(file, []byte("{}"), 0644))

	_, err := config.LoadFromPath(file)
	assert.ErrorContains(t, err, "unsupported config file extension")
}

func TestFileLoader_MissingConfigPath(t *testing.T) {
	loader := &config.FileLoader{}
	_, err := loader.Load(context.Background(), "")
	assert.ErrorContains(t, err, "missing config file path")
}
