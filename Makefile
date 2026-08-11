# Build, test and install waybar-nordvpn

BINARY      := waybar-nordvpn
PKG         := ./...
MAIN        := .

# Install locations. Override on the command line, e.g.
#   make install PREFIX=/usr/local
PREFIX      ?= $(HOME)/.local
BINDIR      ?= $(PREFIX)/bin
ICONDIR     ?= $(PREFIX)/share/icons/hicolor/scalable/apps

# Version metadata, derived from git when available. Stamped into the binary via
# main.version and used to name dist artifacts.
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS     := -s -w -X main.version=$(VERSION)

GO          ?= go
GOFLAGS     ?=
DIST        := dist

.DEFAULT_GOAL := build

.PHONY: all build run test test-race cover bench fmt vet lint tidy check clean install uninstall dist help

## all: format, vet, test and build
all: check build

## build: compile the binary into the project root
build:
	$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BINARY) $(MAIN)

## run: build and run once, printing the waybar JSON line
run: build
	./$(BINARY)

## test: run the unit tests
test:
	$(GO) test $(GOFLAGS) $(PKG)

## test-race: run the unit tests under the race detector
test-race:
	$(GO) test $(GOFLAGS) -race $(PKG)

## cover: run tests with coverage and print a per-function summary
cover:
	$(GO) test $(GOFLAGS) -coverprofile=coverage.out $(PKG)
	$(GO) tool cover -func=coverage.out

## cover-html: render the coverage profile as HTML
cover-html: cover
	$(GO) tool cover -html=coverage.out -o coverage.html

## bench: run benchmarks
bench:
	$(GO) test $(GOFLAGS) -run '^$$' -bench=. -benchmem $(PKG)

## fmt: format all Go sources
fmt:
	$(GO) fmt $(PKG)

## vet: run go vet
vet:
	$(GO) vet $(PKG)

## lint: run golangci-lint when it is installed
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed; skipping (see https://golangci-lint.run)"; \
	fi

## tidy: tidy go.mod / go.sum
tidy:
	$(GO) mod tidy

## check: fmt, vet, lint and test — what CI should run
check: fmt vet lint test

## dist: cross-compile release binaries into dist/
dist:
	@mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -o $(DIST)/$(BINARY)-$(VERSION)-linux-amd64 $(MAIN)
	GOOS=linux GOARCH=arm64 $(GO) build -ldflags '$(LDFLAGS)' -o $(DIST)/$(BINARY)-$(VERSION)-linux-arm64 $(MAIN)

## install: install the binary and icon under $(PREFIX)
install: build
	install -d $(BINDIR) $(ICONDIR)
	install -m 0755 $(BINARY) $(BINDIR)/$(BINARY)
	install -m 0644 icons/nordvpn.svg $(ICONDIR)/nordvpn.svg
	@echo "installed $(BINDIR)/$(BINARY)"

## uninstall: remove the installed binary and icon
uninstall:
	rm -f $(BINDIR)/$(BINARY) $(ICONDIR)/nordvpn.svg

## clean: remove build artifacts
clean:
	rm -f $(BINARY) coverage.out coverage.html
	rm -rf $(DIST)

## help: list available targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /' | sort
