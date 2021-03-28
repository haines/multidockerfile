.DEFAULT_GOAL := all

export CACHE_DIR := $(shell scripts/cache-dir)
export MOD_DIR := $(CACHE_DIR)/mod
export TARGET_DIR := $(CURDIR)/target

export GO_VERSION := $(shell scripts/go-version)

export VERSION := $(file < VERSION)
export GIT_COMMIT := $(shell git rev-parse HEAD)
export BUILT := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

export LDFLAGS := \
	-s \
	-w \
	-X github.com/haines/multidockerfile/internal/version.version=$(VERSION) \
	-X github.com/haines/multidockerfile/internal/version.gitCommit=$(GIT_COMMIT) \
	-X github.com/haines/multidockerfile/internal/version.built=$(BUILT)

export DOCKER_IMAGE := ghcr.io/haines/multidockerfile
DOCKER_PLATFORMS := linux/amd64,linux/arm/v7,linux/arm64/v8

SOURCES := go.mod go.sum $(shell find . -type f -name '*.go' -print)
export BINARY_PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm linux/arm64 windows/amd64
BINARIES := $(foreach platform,$(BINARY_PLATFORMS),$(TARGET_DIR)/multidockerfile-$(subst /,-,$(platform:windows/%=windows/%.exe)))
CHECKSUMS := $(BINARIES:=.sha256)
export GPG_KEY := 6E225DD62262D98AAC77F9CDB16A6F178227A23E
SIGNATURES := $(BINARIES:=.asc)

$(CACHE_DIR):
	@mkdir -p $(CACHE_DIR)

$(MOD_DIR):
	@mkdir -p $(MOD_DIR)

$(TARGET_DIR):
	@mkdir -p $(TARGET_DIR)

$(BINARIES)&: $(SOURCES) | $(CACHE_DIR) $(MOD_DIR) $(TARGET_DIR)
	@scripts/docker-compose run --rm build

$(CHECKSUMS): %.sha256: %
	@scripts/checksum $< $@

$(SIGNATURES): %.asc: %
	@scripts/sign $< $@

.PHONY: help
help:
	@scripts/makefile-help

.PHONY: all
all: lint test ## Lint code and run tests (default)

.PHONY: build
build: $(BINARIES) ## Build binaries for all platforms

.PHONY: clean
clean: ## Clean all build artifacts
	@scripts/clean

.PHONY: docker-build
docker-build: ## Build Docker image for current platform
	@scripts/docker-build --load

.PHONY: docker-push
docker-push: ## Build and push Docker image for all platforms
	@scripts/docker-build --platform $(DOCKER_PLATFORMS) --push

.PHONY: lint
lint: | $(CACHE_DIR) $(MOD_DIR) ## Lint code
	@scripts/docker-compose run --rm lint

.PHONY: release
release: ## Publish a release
	@scripts/release

.PHONY: release-artifacts
release-artifacts: $(BINARIES) $(CHECKSUMS) $(SIGNATURES) docker-push

.PHONY: test
test: | $(CACHE_DIR) $(MOD_DIR) ## Run tests
	@scripts/docker-compose run --rm test
