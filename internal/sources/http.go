package sources

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type HTTPOptions struct {
	Addr      string `mapstructure:"addr"`
	MetaField string `mapstructure:"meta_field"`
}

func (opts *HTTPOptions) Validate() error {
	if opts.Addr == "" {
		return errors.New("addr is required")
	}
	return nil
}

type httpEventSource struct {
	Addr          string // e.g., ":8080"
	MetadataField string // e.g., "_http_meta" (empty string disables)
	events        chan map[string]interface{}
	server        *http.Server
}

func NewHTTPEventSource(opts HTTPOptions) EventSource {
	return &httpEventSource{
		Addr:          opts.Addr,
		MetadataField: opts.MetaField,
		events:        make(chan map[string]interface{}),
	}
}

func (h *httpEventSource) Start(ctx context.Context) (<-chan map[string]interface{}, error) {
	mux := http.NewServeMux()

	//goland:noinspection GoUnhandledErrorResult
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// get input from body
		var input map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if h.MetadataField != "" {
			input[h.MetadataField] = map[string]interface{}{
				"method":       r.Method,
				"path":         r.URL.Path,
				"http_version": r.Proto,
				"headers":      r.Header,
				"remote_addr":  r.RemoteAddr,
			}
		}

		select {
		case h.events <- input:
			w.WriteHeader(http.StatusAccepted)
			break
		default:
			http.Error(w, "event queue full", http.StatusServiceUnavailable)
		}

	})

	h.server = &http.Server{
		Addr:    h.Addr,
		Handler: mux,
	}

	go func() {
		_ = h.server.ListenAndServe()
	}()

	go func() {
		<-ctx.Done()
		_ = h.server.Shutdown(context.Background())
		close(h.events)
	}()

	return h.events, nil
}

func (h *httpEventSource) Close() error {
	return h.server.Shutdown(context.Background())
}
