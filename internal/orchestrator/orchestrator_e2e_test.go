package orchestrator

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/tommed/ducto-dsl/transform"
	"github.com/tommed/ducto-orchestrator/internal/config"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type preprocessorCase struct {
	Preprocessors []config.PluginBlock `json:"preprocessors" mapstructure:"preprocessors"`
}

func TestOrchestrator_E2E_Matrix(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E tests in short mode")
	}

	base := "../../testdata/cases"

	entries, err := os.ReadDir(base)
	assert.NoError(t, err)

	cases := map[string]struct{}{}

	// Auto-detect case prefixes
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".input.json") {
			prefix := strings.TrimSuffix(name, ".input.json")
			cases[prefix] = struct{}{}
		}
	}
	if len(cases) == 0 {
		t.Fatal("no cases found in ./testdata/cases")
	}

	// Run each case
	for prefix := range cases {
		t.Run(prefix, func(t *testing.T) {
			input := loadJSON(t, filepath.Join(base, prefix+".input.json"))
			prog := loadProgram(t, filepath.Join(base, prefix+".program.json"))
			expected := loadJSON(t, filepath.Join(base, prefix+".expected.json"))

			// Load the optional preprocessor config (e.g. feature flags)
			var preprocessorCfg []config.PluginBlock
			preprocessorCfgRaw, err := config.Decode[preprocessorCase](loadJSON(t, filepath.Join(base, prefix+".preprocessors.json")))
			assert.NoError(t, err)
			if preprocessorCfgRaw != nil && len(preprocessorCfgRaw.Preprocessors) > 0 {
				preprocessorCfg = preprocessorCfgRaw.Preprocessors
			}

			ctx := context.Background()
			o := New(prog, false)
			err = o.InstallPreprocessorsFromConfig(ctx, preprocessorCfg)
			assert.NoError(t, err, "unable to install preprocessors")
			writer := &fakeWriter{}

			err = o.RunOnce(context.Background(), input, writer)
			assert.NoError(t, err)

			// We need to compare as JSON objects to avoid map type differences
			expectedJSON, _ := json.Marshal(expected)
			gotJSON, _ := json.Marshal(writer.Written)
			assert.Equal(t, string(expectedJSON), string(gotJSON))
		})
	}
}

func loadJSON(t *testing.T, path string) map[string]interface{} {
	t.Helper()
	data, err := os.ReadFile(path)
	assert.NoError(t, err)
	var out map[string]interface{}
	assert.NoError(t, json.Unmarshal(data, &out))
	return out
}

func loadProgram(t *testing.T, path string) *transform.Program {
	t.Helper()
	prog, err := transform.LoadProgram(path)
	assert.NoError(t, err)
	return prog
}
