<!--suppress HtmlDeprecatedAttribute -->
<p align="right">
    <a href="https://github.com/tommed" title="See Project Ducto">
        <img src="../assets/ducto-logo-small.png" alt="A part of Project Ducto"/>
    </a>
</p>

# Ducto Config-Loader Spec.

This document defines the architecture and responsibilities of `ConfigLoader` implementations across supported environments.

---
## 🎯 Purpose

The `ConfigLoader` abstraction enables ducto-orchestrator to be configured from multiple sources (files, cloud buckets, secrets managers, etc.) while maintaining a clean, portable, testable CLI.

Each platform-specific `cmd/ducto-*` entrypoint passes a different `ConfigLoader` implementation, responsible for:

- Fetching and resolving the main config file (`config.Config`)
- Resolving the embedded or external program (`transform.Program`)
- Performing optional polling or watch-based reload

---
## ✅ Interface

```go
type ConfigLoader interface {
    Load(ctx context.Context, configURI string) (*config.Config, error)
}

// Wishlist: Make hot-reloadable

type WatchingConfigLoader interface {
    ConfigLoader
    Watch(ctx context.Context, configURI string) <-chan struct{}
}
```

---
## 🧱 Entry Point Structure

Each `main.go` entrypoint stays minimal:

```go
package main

import (
    "github.com/tommed/ducto-orchestrator/internal/cli"
    "os"
)

func main() {
    os.Exit(cli.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, &cli.FileLoader{}))
}
```

---
## 🌍 Supported Config Loaders

### 📁 FileLoader (default)
- Accepts: `./config.yaml`, `file://path/to/config.yaml`
- Supports `program_file:` pointing to local file
- Usage:
  ```bash
  ducto-orchestrator --config ./config.yaml
  ```

### ☁️ GCSLoader (GCP only)
- Accepts: `gs://bucket/config.yaml`
- Supports `program_uri: gs://bucket/prog.json`
- Uses GCP IAM credentials for access
- Ideal for Cloud Run / Cloud Functions

### ☁️ S3Loader (AWS)
- Accepts: `s3://bucket/path.yaml`
- Uses AWS credentials
- Supports `program_uri: s3://...`
- Future implementation

### ☁️ AzureBlobLoader (Azure)
- Accepts: `azblob://container/path.yaml`
- Future implementation

---
## 🧠 Program Resolution

The following fields can be used to resolve the DSL program:

Priority order:
1. `program` (inline in YAML)
2. `program_inline` (CLI or ENV override)
3. `program_uri` (resolved by loader)
4. `program_file` (local path)

---
## 🧪 Example Usage

### Local Dev
```bash
ducto-orchestrator --config ./examples/http-source.yaml
```

### GCP Production
```bash
ducto-gcp-cloudrun --config gs://my-org/prod/config.yaml
```

---
## 🔁 Future: Hot Reload

`WatchingConfigLoader` can emit events when the config or DSL program changes.

Strategies:
- Restart the orchestrator (loop + reload)
- Exit and let Cloud Run restart
- Trigger a graceful in-process restart (advanced)

---
## 🧩 Next Steps
- [ ] Implement `FileLoader`
- [ ] Implement `GCSLoader`
- [ ] Wire CLI to accept `--config <uri>`
- [ ] Support program resolution hierarchy
- [ ] Optional: `NewLoaderForScheme()` helper for automatic dispatch
