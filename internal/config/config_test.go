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
	//goland:noinspection GoUnhandledErrorResult
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
		assert.NoError(t, FinalizeConfig(cfg, tmpDir))
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

func TestResolvePath(t *testing.T) {
	type args struct {
		configDir string
		relOrAbs  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "no config directory set",
			args: args{
				configDir: "",
				relOrAbs:  "./docs/doc1.md",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				if err == nil {
					return false
				}
				return err.Error() == "no config directory set"
			},
		},
		{
			name: "rel path",
			args: args{
				configDir: "/tmp",
				relOrAbs:  "docs/doc1.md",
			},
			want: "/tmp/docs/doc1.md",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "abs path",
			args: args{
				configDir: "/var/run",
				relOrAbs:  "/tmp/docs/doc1.md",
			},
			want: "/tmp/docs/doc1.md",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
	}
	for _, tt := range tests {
		// NOT run in parallel as setting configDir in each iteration
		configDir = tt.args.configDir
		got, err := ResolvePath(tt.args.relOrAbs)
		if !tt.wantErr(t, err, fmt.Sprintf("ResolvePath(%v)", tt.args.relOrAbs)) {
			continue
		}
		assert.Equal(t, tt.want, got, "ResolvePath[%s](%v) is not '%s'", tt.name, tt.args.relOrAbs, tt.want)
	}
}
