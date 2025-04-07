package outputs

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPOptions_Validate(t *testing.T) {
	type args struct {
		opts HTTPOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "valid",
			args: args{
				opts: HTTPOptions{
					URL:    "https://example.com",
					Method: http.MethodPost,
					Headers: map[string]string{
						"Accept": "application/json",
					},
				},
			},
		},
		{
			name: "valid Cloud Events",
			args: args{
				opts: HTTPOptions{
					URL:         "https://example.com",
					Method:      http.MethodPost,
					ContentType: "application/cloudevents+json; charset=utf-8",
				},
			},
		},
		{
			name: "bad url",
			args: args{
				opts: HTTPOptions{
					Method: http.MethodPost,
				},
			},
			wantErr: errors.New("url is required"),
		},
		{
			name: "default method",
			args: args{
				opts: HTTPOptions{
					URL: "https://example.com",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.opts.Validate()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestHTTPWriter_SuccessfulRequest(t *testing.T) {
	ctx := context.Background()
	mockClient := newMockClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "application/json; charset=utf-8", req.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", req.Header.Get("Authorization"))
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "https://example.com/api", req.URL.String())

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`ok`)),
		}
	})

	var opts = HTTPOptions{
		URL:    "https://example.com/api",
		Method: "POST",
		Token:  "test-token",
		Headers: map[string]string{
			"Accept": "application/json",
		},
	}
	_ = opts.Validate() // applies defaults, usually called by `FromPlugin`
	writer := NewHTTPWriterWithClient(opts, mockClient)

	err := writer.WriteOutput(ctx, map[string]interface{}{"foo": "bar"})
	assert.NoError(t, err)
}

func TestHTTPWriter_InvalidStatusCode(t *testing.T) {
	ctx := context.Background()
	mockClient := newMockClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 418,
			Body:       io.NopCloser(strings.NewReader(`I'm a teapot`)),
		}
	})

	writer := NewHTTPWriterWithClient(HTTPOptions{
		URL:              "https://example.com/fail",
		Method:           "POST",
		ExpectStatusCode: 200,
	}, mockClient)

	err := writer.WriteOutput(ctx, map[string]interface{}{"fail": true})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected status code 200 but got 418")
}

func TestHTTPWriter_InvalidJSON(t *testing.T) {
	ctx := context.Background()
	writer := NewHTTPWriterWithClient(HTTPOptions{
		URL:    "https://example.com",
		Method: "POST",
	}, http.DefaultClient)

	// use a type that can't be marshalled
	invalid := map[string]interface{}{
		"bad": func() {},
	}

	err := writer.WriteOutput(ctx, invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json: unsupported type")
}

// roundTripFunc is a test helper
type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newMockClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}
