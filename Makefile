BINDIR     := $(CURDIR)/bin
BINNAME    ?= gograph

GOPATH        = $(shell go env GOPATH)
DEP           = $(GOPATH)/bin/dep
GOX           = $(GOPATH)/bin/gox
GOIMPORTS     = $(GOPATH)/bin/goimports
GOLANGCI_LINT = $(GOPATH)/bin/golangci-lint

TARGETS := ${notdir ${shell find ./cmd/* -type d -print}}

# go option
PKG        := ./...
TAGS       :=
TESTS      := .
TESTFLAGS  :=
LDFLAGS    := -w -s
GOFLAGS    :=
SRC        := $(shell find . -type f -name '*.go' -print)

# Required for globs to work correctly
SHELL      = /bin/bash

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")


ifdef VERSION
	BINARY_VERSION = $(VERSION)
endif
BINARY_VERSION ?= ${GIT_TAG}

# Only set Version if building a tag or VERSION is set
ifneq ($(BINARY_VERSION),)
	LDFLAGS += -X cybermats/gograph/internal/version.version=${BINARY_VERSION}
endif
# Clear the "unreleased" string in BuildMetadata
ifneq ($(GIT_TAG),)
	LDFLAGS += -X cybermats/gograph/internal/version.metadata=
endif
LDFLAGS += -X cybermats/gograph/internal/version.gitCommit=${GIT_COMMIT}
LDFLAGS += -X cybermats/gograph/internal/version.gitTreeState=${GIT_DIRTY}

.PHONY: all
all: build vet test

# ------------------------------------------------------------------------------
#  build

.PHONY: build
build: $(TARGETS)
	@true

$(TARGETS): $(SRC)
	GO111MODULE=on go build $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $(BINDIR)/$@ ./cmd/$@


# ------------------------------------------------------------------------------
#  vet

.PHONY: vet
vet: build
vet: vet-code

.PHONY: vet-code
vet-code:
	@echo
	@echo "==> Running vet <=="
	GO111MODULE=on go vet $(GOFLAGS) $(PKG)

# ------------------------------------------------------------------------------
#  test

.PHONY: test
test: build
test: vet-code
test: TESTFLAGS += -v
test: test-unit

.PHONY: test-unit
test-unit:
	@echo
	@echo "==> Running unit tests <=="
	GO111MODULE=on go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)


.PHONY: format
format: $(GOIMPORTS)
	GO111MODULE=on go list -f '{{.Dir}}' ./... | xargs $(GOIMPORTS) -w -local cybermats/gograph

# ------------------------------------------------------------------------------

.PHONY: clean
clean:
	@rm -rf $(BINDIR) ./_dist

.PHONY: info
info:
	 @echo "Version:           ${VERSION}"
	 @echo "Git Tag:           ${GIT_TAG}"
	 @echo "Git Commit:        ${GIT_COMMIT}"
	 @echo "Git Tree State:    ${GIT_DIRTY}"
