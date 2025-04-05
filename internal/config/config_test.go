package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDecode_Success(t *testing.T) {
	raw := map[string]interface{}{
		"addr":       ":8080",
		"meta_field": "_http",
	}

	type HTTPOpts struct {
		Addr      string `mapstructure:"addr"`
		MetaField string `mapstructure:"meta_field"`
	}

	opts, err := Decode[HTTPOpts](raw)
	assert.NoError(t, err)
	assert.Equal(t, ":8080", opts.Addr)
	assert.Equal(t, "_http", opts.MetaField)
}

func TestLoad_InvalidPath(t *testing.T) {
	_, err := Load("testdata/does-not-exist.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config")
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmp := writeTempFile(t, []byte(`{ invalid_yaml`))
	defer os.Remove(tmp)

	_, err := Load(tmp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config")
}

func TestLoad_ProgramPathResolution(t *testing.T) {
	programPath := "relative/prog.json"
	tmpCfg := []byte(fmt.Sprintf("program_file: %s", programPath))
	tmpFile := writeTempFile(t, tmpCfg)
	defer os.Remove(tmpFile)

	cfg, err := Load(tmpFile)
	assert.NoError(t, err)
	assert.True(t, filepath.IsAbs(cfg.ProgramFile))
	assert.Contains(t, cfg.ProgramFile, "relative/prog.json")
}

func TestDecode_Errors(t *testing.T) {
	// force decode error
	raw := map[string]interface{}{
		"callback": func() {},
	}
	type Bad struct {
		Callback int `mapstructure:"callback"`
	}

	_, err := Decode[Bad](raw)
	assert.Error(t, err)
}

func writeTempFile(t *testing.T, content []byte) string {
	t.Helper()
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("test-%d.yaml", time.Now().UnixNano()))
	require.NoError(t, os.WriteFile(tmp, content, 0644))
	return tmp
}
