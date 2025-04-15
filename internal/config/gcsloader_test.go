// File: ducto-orchestrator/internal/config/gcsloader_test.go
package config

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGCSLoader_InvalidURIScheme(t *testing.T) {
	loader := &GCSLoaderMock{}
	_, err := loader.Load(context.Background(), "/not/a/gsuri.yaml")
	assert.ErrorContains(t, err, "invalid URI scheme")
}

func TestGCSLoader_InvalidGCSURIFormat(t *testing.T) {
	loader := &GCSLoaderMock{}
	_, err := loader.Load(context.Background(), "gs://bucketonly")
	assert.ErrorContains(t, err, "invalid GCS URI")
}

// GCSLoaderMock is a stub to isolate URI parsing logic
type GCSLoaderMock struct{}

func (m *GCSLoaderMock) Load(ctx context.Context, uri string) (*Config, error) {
	l := &gcsLoader{}
	return l.Load(ctx, uri)
}

func TestGCSLoader_Load_PublicExample(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	loader, err := NewGCSLoader()
	require.NoError(t, err)

	ctx := context.Background()
	cfg, err := loader.Load(ctx, "gs://ducto-public/02-gcs_config.yaml")
	require.NoError(t, err)
	require.NotNil(t, cfg.Program, "expected a loaded DSL program")
	require.Equal(t, 1, len(cfg.Program.Instructions))
}
