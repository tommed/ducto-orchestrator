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

.PHONY: all check ci lint test-unit test-e2e test-full coverage example-simplest clean # build-all ducto-orchestrator-macos ducto-orchestrator-windows

all: check

check: lint test-full coverage

#build-all: ducto-orchestrator-macos ducto-orchestrator-windows

ci: check # build-all

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

test-unit:
	@echo "==> Running short tests"
	$(GO) test -short -coverprofile=$(COVERAGE_OUT) -covermode=atomic -v ./...
	$(GO) tool cover -func=$(COVERAGE_OUT)

test-full:
	@echo "==> Running all tests"
	$(GO) test -coverprofile=$(COVERAGE_OUT) -covermode=atomic -v ./...
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
	echo '{"foo":"bar"}' | $(GO) run ./cmd/ducto-orchestrator -debug -program examples/simplest.json

#ducto-orchestrator-macos:
#	@echo "==> Building macOS CLI"
#	$(GO) build -o ducto-dsl ./cmd/ducto-dsl
#
#ducto-orchestrator-windows:
#	@echo "==> Building Windows CLI"
#	GOOS=windows GOARCH=amd64 $(GO) build -o ducto-dsl.exe ./cmd/ducto-dsl