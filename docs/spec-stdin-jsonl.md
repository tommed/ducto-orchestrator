# Ducto EventSource: stdin jsonl mode

This spec defines the optional `jsonl` mode for the `stdin` EventSource plugin in Ducto. 
It allows streaming newline-delimited JSON events (JSONL) into the orchestrator from 
standard input, enabling integration with data generation tools like `ducto-faker` or 
external pipelines.

---
## 🎯 Goals

- Retain support for single JSON input (default mode)
- Add support for streaming JSONL input (one JSON object per line)
- Make the feature opt-in via a `jsonl: true` flag
- Ensure full compatibility with orchestrator's existing event loop
- Maintain graceful shutdown behavior

---
## 📦 Configuration

### JSON Example

```json5
{
  "source": {
    "type": "stdin",
    "config": {
      "jsonl": true
    }
  }
}
```

### YAML Example

```yaml
source:
  type: stdin
  config:
    jsonl: true
```

### Go Struct

```go
type StdinOptions struct {
	JSONL bool `mapstructure:"jsonl,omitempty" json:"jsonl,omitempty"`
}
```

If `jsonl` is omitted or false, stdin behaves as before, reading a single JSON object.

---
## 🧠 Runtime Behavior

- When `jsonl` is **false or omitted**:
    - The orchestrator reads one JSON object from stdin and terminates the stream.

- When `jsonl` is **true**:
    - Each line is parsed as an individual JSON object and pushed into the event channel
    - EOF or invalid JSON terminates the stream

---
## ✅ Use Case: Ducto Faker

```bash
ducto-faker generate -format jsonl | ducto-orchestrator --config my-pipeline.yaml
```

This allows high-volume pseudo-events to be generated and streamed through the 
orchestrator via stdin.

---
## 🔐 Graceful Shutdown

- The implementation respects the provided context
- On SIGINT/SIGTERM, reading stops cleanly and closes the event channel

---
## 🔮 Future Extensions

To support additional streaming input formats in the future:

```yaml
source:
  type: stdin
  config:
    format: jsonl   # or "json", "csv", etc.
```

But for now, the simpler `jsonl: true` design keeps it lean and focused.

---
## 📁 Implementation Path

- Update `stdinEventSource` in `internal/sources/stdin.go`
- Add `StdinOptions` struct
- Parse `jsonl` field from plugin config
- Stream objects line-by-line into event channel if enabled

---
## 🧪 Testing

- Validate single-object mode remains unchanged
- Pipe JSONL content via `echo` or `cat` to verify streaming mode
- Add unit tests for both modes

---
## ✅ Summary

`stdin` now supports a simple, extensible `jsonl` mode that turns standard input 
into a JSON event stream. Ideal for dev tools and batch workflows, this feature 
lays the foundation for full synthetic pipelines with Ducto.
