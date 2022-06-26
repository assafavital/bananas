# Build tags, this skips the btrfs dependency in https://github.com/containers/storage imported from https://github.com/containers/image
BUILD_TAGS = containers_image_storage_stub
ifdef GIT_TAG
 	# Non empty versions are considered production, set empty if building a dev version
	ifeq ($(GIT_TAG), dev)
		BUILD_TAGS += dev
	else
		RAFTT_VERSION := $(GIT_TAG)
	endif
else
	RAFTT_VERSION := $(shell git tag --points-at HEAD)
endif

# If RAFTT_VERSION exists and starts with "v"
ifdef RAFTT_VERSION
	ifeq ($(filter v%,$(RAFTT_VERSION)), $(RAFTT_VERSION))
		BUILD_TAGS += prod
	endif
endif

ifdef RC_TARGET
	ifneq ($(RC_TARGET), "")
		BUILD_TAGS += release_candidate
	endif
endif

WHEN := $(shell date "+DATE:%Y-%m-%d,TIME:%H:%M:%S%z")

SYNCTHING_BUILD_TAGS += noassets
BUILD_TAGS += $(SYNCTHING_BUILD_TAGS)

BUILD_DIR = build

# a macro that joins a list of values using the given arg, e.g $(call joinlist,a b c,-) -> a-b-c
joinlist = $(subst $(eval) ,$2,$1)

GO = go
ifdef GIT_COMMIT
	GIT_COMMIT := $(GIT_COMMIT)
else
	GIT_COMMIT := $(shell git rev-list -1 HEAD)
endif
PKG_PATH := raftt.io/bananas/pkg

ifeq ($(OS),Windows_NT)
	BUILD_MACHINE := windows
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Darwin)
		BUILD_MACHINE := mac
	else ifeq ($(UNAME_S),Linux)
		BUILD_MACHINE := linux
	else
		$(error Unsupported OS $(UNAME_S))
	endif
endif

BUILD_PROPS = '$(PKG_PATH)/version/properties.Commit=$(GIT_COMMIT)'
BUILD_PROPS += '$(PKG_PATH)/version/properties.Version=$(RAFTT_VERSION)'
BUILD_PROPS += '$(PKG_PATH)/version/properties.When=$(WHEN)'

REMOVE_SYMBOLS_FLAGS =
ifndef WITH_SYMBOLS
	REMOVE_SYMBOLS_FLAGS = -s -w
endif

LDFLAGS = $(REMOVE_SYMBOLS_FLAGS) $(addprefix -X ,$(BUILD_PROPS))

comma := ,
JOINED_TAGS = $(call joinlist,$(BUILD_TAGS),$(comma))

GOFLAGS = -ldflags "$(LDFLAGS)" -tags "$(JOINED_TAGS)"

GOBUILD = $(GO) build $(GOFLAGS)
GORUN = $(GO) run $(GOFLAGS)
GOCLEAN = $(GO) clean
GOLINT = golangci-lint run --max-same-issues 50 --allow-parallel-runners
GOPATH = $(shell $(GO) env GOPATH)

ifeq ($(BUILD_MACHINE), windows)
	GOTESTSUM = $(GOPATH)\\bin\\gotestsum
else
	GOTESTSUM = $(GOPATH)/bin/gotestsum
endif
DOCKER = docker

# Protobuf gen cmd
PBGEN_GO := protoc --go_out=paths=source_relative:./

# Ignore protobuf changes; for container builders
ifndef NOPROTOGEN
PB_SOURCES := $(shell find . -type f -name '*.proto')
PB_TARGETS := $(patsubst %.proto,%.pb.go,$(PB_SOURCES))
endif

GO_SOURCES := $(shell find . -type f -name '*.go' -or -name '*.go.tmpl' -or -name 'go.mod') $(PB_SOURCES)

# Run all build steps
.PHONY: all
all: tidy lint build test

# Fix golang dependencies
.PHONY: tidy
tidy:
	$(GO) mod tidy

PYTHON_SOURCE_DIRS := $(shell git ls-files '*.py' | cut -d '/' -f 1 | uniq)

linters = lint-py lint-go lint-yaml lint-protobuf lint-shell lint-scripts

.PHONY: $(linters)
lint-py:
	pylint $(PYTHON_SOURCE_DIRS)
	mypy $(PYTHON_SOURCE_DIRS)

lint-go: lint-protobuf
	$(GOLINT) -v

lint-protobuf:
	buf lint

lint-yaml:
	yamllint .

lint-shell:
	find . -not \( -path "./deploy/backend_kube/.terraform" -prune \) -type f -name "*.sh" | xargs shellcheck

lint-scripts: lint-shell lint-py lint-yaml

# Run all linters
.PHONY: lint
lint: $(linters)

# Run tests
.PHONY: gotestsum
gotestsum:
	@if ! [ -x $(GOTESTSUM) ]; then \
		echo Installing gotestsum; \
		$(GO) install gotest.tools/gotestsum@latest; \
		$(GO) mod tidy; \
	fi

.PHONY: test
test: PKGS_TO_TEST = ./...
test: --run-tests

.PHONY: testcli
testcli: PKGS_TO_TEST = ./cmd/raftt/... ./pkg/...
testcli: --run-tests

.PHONY: --run-tests
--run-tests: gotestsum
--run-tests: BUILD_TAGS += test_version
--run-tests:
	gotestsum --packages="$(PKGS_TO_TEST)" -- $(GOFLAGS)

.PHONY: docs
docs:
	$(GORUN) cmd/docs/main.go > "docs/docs/cli_reference.md"
# All build targets / entry points
TARGETS := admiral lifeguard cli backoffice sandcastle in-env-raftt container-proxy
.PHONY: $(TARGETS)

# default rules for local builds
$(TARGETS): %: $(BUILD_DIR)/$(BUILD_MACHINE)/%

PKG_DIR := .

$(BUILD_DIR)/%/admiral: PKG_MAIN := cmd/admiral/admiral.go
$(BUILD_DIR)/%/lifeguard: PKG_MAIN := cmd/lifeguard/*.go
$(BUILD_DIR)/%/cli: PKG_MAIN := cmd/raftt/raftt.go
$(BUILD_DIR)/%/sandcastle: PKG_MAIN := cmd/sandcastle/main.go
$(BUILD_DIR)/%/backoffice: PKG_MAIN := cmd/backoffice/*.go
$(BUILD_DIR)/%/in-env-raftt: PKG_MAIN := cmd/in-env-raftt/*.go
$(BUILD_DIR)/%/container-proxy: PKG_MAIN := cmd/container-proxy/*.go

# macro for per-OS build target
TARGETS_OS = $(addprefix $(BUILD_DIR)/$1/,$(TARGETS))

TARGETS_MAC := $(call TARGETS_OS,mac)
TARGETS_LINUX := $(call TARGETS_OS,linux)
TARGETS_WINDOWS := $(call TARGETS_OS,windows)

# per-OS build targets
OS_LIST := mac linux windows
BUILD_OS = $(addprefix $(BUILD_DIR)/,$(OS_LIST))

$(BUILD_DIR)/mac: $(TARGETS_MAC)
$(BUILD_DIR)/linux: $(TARGETS_LINUX)
$(BUILD_DIR)/windows: $(TARGETS_WINDOWS)

.PHONY: $(BUILD_DIR)/linux $(BUILD_DIR)/mac $(BUILD_DIR)/windows

$(TARGETS_MAC): $(BUILD_DIR)/%: $(GO_SOURCES)
	GOOS=darwin $(GOBUILD) -o "$@" $(PKG_MAIN)

$(TARGETS_LINUX): $(BUILD_DIR)/%: $(GO_SOURCES)
	GOOS=linux CGO_ENABLED=0 $(GOBUILD) -o "$@" $(PKG_MAIN)

$(TARGETS_WINDOWS): $(BUILD_DIR)/%: $(GO_SOURCES)
	GOOS=windows CGO_ENABLED=0 $(GOBUILD) -o "$@".exe $(PKG_MAIN)

# Build all entry points
.PHONY: build
build: $(TARGETS)

# Build protobuf targets
.PHONY: protobuf
protobuf: $(PB_TARGETS)

$(PB_TARGETS): %.pb.go: %.proto
	$(PBGEN_GO) $<

# Build docker images
DOCKER_IMAGES := docker/sandcastle docker/admiral docker/backoffice
.PHONY: docker
docker: $(DOCKER_IMAGES)

docker/%: DOCKER_LABEL=$*

.PHONY: $(DOCKER_IMAGES)
$(DOCKER_IMAGES): docker/%: docker/%/Dockerfile
	$(DOCKER) build -t $(DOCKER_LABEL) -f $< .

# Build raftt compose
.PHONY: compose
compose: admiral docker
	docker-compose -f .raftt/raftt-compose.yml build
	docker-compose -f .raftt/raftt-compose.yml up

# Clean build files
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)/*

# Format all go files
.PHONY: fmt
fmt:
	./scripts/gofmtall.py

.PHONY: debug-starlark
debug-starlark:
	$(GORUN) cmd/debug-starlark/main.go "$(RAFTTLARK_ARGS)"


.PHONY: help
help:
	@echo "Raftt makefile"
	@echo "=============="
	@echo "Building:"
	@echo "  build            - builds all executable targets for current machine:"
	@echo "    "$(call joinlist,$(TARGETS),"\n    ")
	@echo
	@echo "  <build-dir>/<os> - cross-compilation targets:"
	@echo "    "$(call joinlist,$(BUILD_OS),"\n    ")
	@echo "                   - Cross-compile targets"
	@echo "  protobuf         - build precompiled protobuf targets"
	@echo
	@echo "Docker images:"
	@echo "  docker - build all docker images"
	@echo "      $(DOCKER_IMAGES)"
	@echo "  compose - build docker compose targets"
	@echo
	@echo "Linters:"
	@echo "  $(linters)"
	@echo
	@echo "Misc targets:"
	@echo "  docs           - Generate markdown documentation for the CLI"
	@echo "  clean          - Delete build artifacts and go cache"
	@echo "  fmt            - Run code formatters"
	@echo "  test  	        - Run unit tests"
	@echo "  testcli        - Run cli related unit tests"
	@echo "  lint           - Run all linters"
	@echo "  tidy           - Run \`go mod tidy\`"
	@echo "  debug-starlark - Simulate a starlark config debug using local raftt.yml"
	@echo "  help           - Show this message"
	@echo
	@echo "Environment:"
	@echo "  BUILD_DIR=$(BUILD_DIR)/"
	@echo "  BUILD_MACHINE=$(BUILD_MACHINE) ($(OS))"
	@echo "  GO=$(GO) ($$($(GO) version))"
	@echo "  GOPATH=$(GOPATH)"
	@echo "  RAFTT_VERSION=$(RAFTT_VERSION)"
	@echo "  PYTHON_SOURCE_DIRS=$(PYTHON_SOURCE_DIRS)"
	@echo "  GIT_TAG - Set to override RAFTT_VERSION, set to 'dev' to make it empty"
	@echo "  WITH_SYMBOLS - set to enable debug symbols"
