package orchestrator

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	flagsdk "github.com/tommed/ducto-featureflags/sdk"
)

func TestParseTagsEnv(t *testing.T) {
	t.Setenv("DUCTO_TAGS", "env=prod;group=beta;customer=mps")

	got := parseTagsEnv(os.Getenv("DUCTO_TAGS"))
	assert.Equal(t, flagsdk.EvalContext{
		"env":      "prod",
		"group":    "beta",
		"customer": "mps",
	}, got)
}

func TestParseTagsEnv_HandlesWhitespaceAndMalformedPairs(t *testing.T) {
	input := " key1 = value1 ; key2=value2 ; malformed ; empty= "
	expected := flagsdk.EvalContext{
		"key1":  "value1",
		"key2":  "value2",
		"empty": "",
	}

	got := parseTagsEnv(input)
	assert.Equal(t, expected, got)
}

func TestFlagInjectorOptions_Tags_PrefersEnvOverRaw(t *testing.T) {
	t.Setenv("DUCTO_TAGS", "env=prod;foo=bar")

	opts := FlagInjectorOptions{
		RawTags: flagsdk.EvalContext{
			"env": "dev",
			"foo": "nope",
		},
		TagsEnv: "DUCTO_TAGS",
	}

	got := opts.Tags()
	assert.Equal(t, flagsdk.EvalContext{
		"env": "prod",
		"foo": "bar",
	}, got)
}

func TestFlagInjectorOptions_Tags_UsesRawWhenEnvMissing(t *testing.T) {
	opts := FlagInjectorOptions{
		RawTags: flagsdk.EvalContext{
			"env":   "staging",
			"group": "alpha",
		},
		TagsEnv: "", // no override
	}

	got := opts.Tags()
	assert.Equal(t, opts.RawTags, got)
}
