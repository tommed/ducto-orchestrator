package sources

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//goland:noinspection GoUnhandledErrorResult
//goland:noinspection HttpUrlsUsage,GoUnhandledErrorResult
func TestHTTPEventSource_SuccessfulEvent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Build server and hijack it to httptest
	addr := findFreePort(t)
	source := NewHTTPEventSource(HTTPOptions{
		Addr:      addr,
		MetaField: "_meta",
	}).(*httpEventSource)
	server := hijackHTTPEventSource(ctx, source)
	defer server.Close()
	time.Sleep(100 * time.Millisecond)

	events, err := source.Start(ctx)
	require.NoError(t, err)
	require.NotNil(t, events)

	// Prepare input
	input := map[string]interface{}{"foo": "bar"}
	body, _ := json.Marshal(input)

	// Act
	go func() {
		resp, err := http.Post("http://"+source.server.Addr, "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()
		_, _ = io.ReadAll(resp.Body) // drain body
	}()

	// Assert
	select {
	case event := <-events:
		assert.Equal(t, "bar", event["foo"])
		meta, ok := event["_meta"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "POST", meta["method"])
		assert.Equal(t, "/", meta["path"])
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event")
	}
}

//goland:noinspection GoUnhandledErrorResult
func findFreePort(t *testing.T) string {
	ln, err := net.Listen("tcp", "127.0.0.1:0") // use :0 for random port
	require.NoError(t, err)
	defer ln.Close()
	return ln.Addr().String()
}

//goland:noinspection GoUnhandledErrorResult
func hijackHTTPEventSource(ctx context.Context, source *httpEventSource) *httptest.Server {
	// Hijack server to a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = source.Start(ctx)
		defer source.Close()
	}))
	// Inject the real server mux
	source.server = &http.Server{
		Addr:    source.Addr,
		Handler: http.NewServeMux(),
	}
	return server
}
