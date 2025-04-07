package outputs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

type HTTPOptions struct {
	URL              string            `mapstructure:"url"`
	Method           string            `mapstructure:"method"` // Assumes POST if unset
	TimeoutSeconds   int               `mapstructure:"timeout_seconds"`
	ExpectStatusCode int               `mapstructure:"expect_status_code"` // Optional
	ContentType      string            `mapstructure:"content_type"`       // Optional (assumes application/json)
	Headers          map[string]string `mapstructure:"headers"`            // Optional additional headers

	// For tokens, you can embed or access via an environment variable
	TokenType string `mapstructure:"token_type"` // Assumes Bearer if not set but Token/EnvToken is
	Token     string `mapstructure:"token"`
	EnvToken  string `mapstructure:"env_token"` // Name of env token to get token from
}

func (opts *HTTPOptions) Validate() error {
	if opts.URL == "" {
		return errors.New("url is required")
	}
	if opts.Method == "" {
		opts.Method = http.MethodPost
	}
	if opts.ContentType == "" {
		opts.ContentType = "application/json; charset=utf-8"
	}
	return nil
}

type httpOutput struct {
	opts   HTTPOptions
	client *http.Client // Set in unit tests manually
}

func NewHTTPWriter(opts HTTPOptions) OutputWriter {
	return NewHTTPWriterWithClient(opts, &http.Client{
		Timeout: time.Duration(opts.TimeoutSeconds) * time.Second,
	})
}

// NewHTTPWriterWithClient allows injection of a custom HTTP client (for testing).
func NewHTTPWriterWithClient(opts HTTPOptions, client *http.Client) OutputWriter {
	return &httpOutput{
		opts:   opts,
		client: client,
	}
}

// WriteOutput sends the JSON payload to the configured endpoint.
func (h *httpOutput) WriteOutput(ctx context.Context, input map[string]interface{}) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}

	// Build request with context
	req, err := http.NewRequestWithContext(ctx, h.opts.Method, h.opts.URL, bytes.NewReader(data))
	if err != nil {
		return err
	}

	// Additional headers
	for k, v := range h.opts.Headers {
		req.Header.Set(k, v)
	}

	// Security handling
	token := h.opts.Token
	if h.opts.EnvToken != "" {
		token = os.Getenv(h.opts.EnvToken)
	}
	if token != "" {
		tokenType := h.opts.TokenType
		if tokenType == "" {
			tokenType = "Bearer"
		}
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", tokenType, token))
	}

	// Content type (which we could eventually inject if needed)
	req.Header.Set("Content-Type", h.opts.ContentType)

	// Do the request
	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// Validate the status code
	if h.opts.ExpectStatusCode > 0 && resp.StatusCode != h.opts.ExpectStatusCode {
		return fmt.Errorf("expected status code %d but got %d", h.opts.ExpectStatusCode, resp.StatusCode)
	}

	return nil
}
