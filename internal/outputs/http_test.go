package outputs

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPWriter_SuccessfulRequest(t *testing.T) {
	mockClient := newMockClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", req.Header.Get("Authorization"))
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "https://example.com/api", req.URL.String())

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`ok`)),
		}
	})

	writer := NewHTTPWriterWithClient(HTTPOptions{
		URL:    "https://example.com/api",
		Method: "POST",
		Token:  "test-token",
	}, mockClient)

	err := writer.WriteOutput(map[string]interface{}{"foo": "bar"})
	assert.NoError(t, err)
}

func TestHTTPWriter_InvalidStatusCode(t *testing.T) {
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

	err := writer.WriteOutput(map[string]interface{}{"fail": true})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected status code 200 but got 418")
}

func TestHTTPWriter_InvalidJSON(t *testing.T) {
	writer := NewHTTPWriterWithClient(HTTPOptions{
		URL:    "https://example.com",
		Method: "POST",
	}, http.DefaultClient)

	// use a type that can't be marshalled
	invalid := map[string]interface{}{
		"bad": func() {},
	}

	err := writer.WriteOutput(invalid)
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
