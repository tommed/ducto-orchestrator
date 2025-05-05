package config

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
)

type gcsLoader struct {
	client *storage.Client
}

func MustCreateGCSLoader() Loader {
	client, err := NewGCSLoader()
	if err != nil {
		panic(err)
	}
	return client
}

// NewGCSLoader builds a new Loader which supports `gs://` GCS URIs.
func NewGCSLoader() (Loader, error) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}
	return &gcsLoader{client: client}, nil
}

// Load downloads the config from a GCS URI (gs://bucket/key.yaml) and resolves DSL if needed.
//
//goland:noinspection GoUnhandledErrorResult
func (g *gcsLoader) Load(ctx context.Context, uri string) (*Config, error) {
	if !strings.HasPrefix(uri, "gs://") {
		return nil, fmt.Errorf("gcsloader: invalid URI scheme: %s", uri)
	}

	bucket, object, err := parseGCSURI(uri)
	if err != nil {
		return nil, err
	}

	tmpFile, err := os.CreateTemp("", "ducto-config-*.yaml")
	if err != nil {
		return nil, fmt.Errorf("gcsloader: temp file: %w", err)
	}
	defer tmpFile.Close()

	rc, err := g.client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("gcsloader: read object: %w", err)
	}
	defer rc.Close()

	if _, err := io.Copy(tmpFile, rc); err != nil {
		return nil, fmt.Errorf("gcsloader: download config: %w", err)
	}

	cfg, err := ParseConfig(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	// If the program_file is itself a gs:// URI, resolve and load it manually
	if strings.HasPrefix(cfg.ProgramFile, "gs://") {
		progBucket, progObject, err := parseGCSURI(cfg.ProgramFile)
		if err != nil {
			return nil, fmt.Errorf("gcsloader: program_file invalid URI: %w", err)
		}

		progTmp, err := os.CreateTemp("", "ducto-program-*.json")
		if err != nil {
			return nil, fmt.Errorf("gcsloader: temp file for program: %w", err)
		}
		defer progTmp.Close()

		prc, err := g.client.Bucket(progBucket).Object(progObject).NewReader(ctx)
		if err != nil {
			return nil, fmt.Errorf("gcsloader: read program_file object: %w", err)
		}
		defer prc.Close()

		if _, err := io.Copy(progTmp, prc); err != nil {
			return nil, fmt.Errorf("gcsloader: download program_file: %w", err)
		}

		cfg.ProgramFile = progTmp.Name()
	}

	if err := FinalizeConfig(cfg, filepath.Dir(tmpFile.Name())); err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseGCSURI(uri string) (bucket, object string, err error) {
	parts := strings.SplitN(strings.TrimPrefix(uri, "gs://"), "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("gcsloader: invalid GCS URI: %s", uri)
	}
	return parts[0], parts[1], nil
}
