# ----------------------
# Configuration
# ----------------------

COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html
GO=go
LINTER=golangci-lint
LINTER_REMOTE=github.com/golangci/golangci-lint/cmd/golangci-lint@latest
LINTER_OPTS=--timeout=2m

# ----------------------
# General Targets
# ----------------------

.PHONY: all check ci lint test-unit test-e2e test-full coverage example-simplest clean build-all ducto-orchestrator-macos ducto-orchestrator-windows gcp-pubsub-emulator

all: check

check: lint test-full coverage

build-all: ducto-orchestrator-macos ducto-orchestrator-windows

ci: check build-all

clean:
	@rm -f $(COVERAGE_OUT) $(COVERAGE_HTML) ducto-orchestrator*

# ----------------------
# Linting
# ----------------------

lint:
	@echo "==> Running linter"
	$(LINTER) run $(LINTER_OPTS)

lint-install:
	go install $(LINTER_REMOTE)

# ----------------------
# Testing
# ----------------------

gcp-pubsub-emulator:
	@echo "==> Starting Google Pub/Sub Emulator"
	@gcloud beta emulators pubsub start --project=test-project --host-port=0.0.0.0:8085

test-unit:
	@echo "==> Running short tests"
	$(GO) test -short -coverpkg=./... -coverprofile=$(COVERAGE_OUT) -covermode=atomic -v ./...
	$(GO) tool cover -func=$(COVERAGE_OUT)

test-full:
	@echo "==> Running all tests"
	@PUBSUB_EMULATOR_HOST=localhost:8085 GOOGLE_CLOUD_PROJECT=test-project $(GO) test -coverpkg=./... -coverprofile=$(COVERAGE_OUT) -covermode=atomic -v ./...
	$(GO) tool cover -func=$(COVERAGE_OUT)

test-e2e:
	@echo "==> Running full tests"
	$(GO) test -coverprofile=$(COVERAGE_OUT) -covermode=atomic -v -run E2E ./...
	$(GO) tool cover -func=$(COVERAGE_OUT)

coverage:
	@echo "==> Generating coverage HTML report"
	$(GO) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

# ----------------------
# CLI
# ----------------------

example-simplest:
	@echo "==> Building Example: Simplest"
	echo '{"foo":"bar"}' | $(GO) run ./cmd/ducto-orchestrator -debug -program examples/01-simplest.yaml

ducto-orchestrator-macos:
	@echo "==> Building macOS CLI"
	GOOS=darwin GOARCH=arm64 $(GO) build -o ducto-orchestrator ./cmd/ducto-orchestrator

ducto-orchestrator-windows:
	@echo "==> Building Windows CLI"
	GOOS=windows GOARCH=amd64 $(GO) build -o ducto-orchestrator.exe ./cmd/ducto-orchestrator