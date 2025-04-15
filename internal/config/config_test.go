package config

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tommed/ducto-dsl/transform"
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

func TestParseConfig_InvalidPath(t *testing.T) {
	_, err := ParseConfig("testdata/does-not-exist.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestParseConfig_InvalidYAML(t *testing.T) {
	tmp := writeTempFile(t, []byte(`{ invalid_yaml`))
	defer os.Remove(tmp)

	_, err := ParseConfig(tmp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "decode config: yaml:")
}

//goland:noinspection GoUnhandledErrorResult
func TestLoad_ProgramPathResolution(t *testing.T) {
	programPath := "prog.json"
	tmpCfg := []byte(fmt.Sprintf("program_file: '%s'", programPath))
	tmpFile := writeTempFile(t, tmpCfg)
	tmpDir := filepath.Dir(tmpFile)
	defer os.Remove(tmpFile)

	t.Run("no program file", func(t *testing.T) {
		cfg, err := ParseConfig(tmpFile)
		assert.NoError(t, err)
		err = FinalizeConfig(cfg, tmpDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no such file or directory")
	})

	fullProgramPath := filepath.Join(tmpDir, programPath)
	programJSON, _ := json.Marshal(transform.Program{
		Version: 1,
		Instructions: []transform.Instruction{
			{
				Op:    "set",
				Key:   "field1",
				Value: "test1",
			},
		},
	})
	err := os.WriteFile(fullProgramPath, programJSON, os.ModePerm)
	assert.NoError(t, err)
	defer os.Remove(fullProgramPath)

	t.Run("with program file", func(t *testing.T) {
		cfg, err := ParseConfig(tmpFile)
		assert.NoError(t, err)
		err = FinalizeConfig(cfg, tmpDir)
		assert.Equal(t, "", cfg.ProgramFile)
		assert.NotNil(t, "", cfg.Program)
	})
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
