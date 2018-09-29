.DEFAULT_GOAL := build
timestamp = $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
git_revision = $(shell git describe --always --abbrev=10 --dirty)
add_metadata = -ldflags "-X main.buildTime=$(timestamp) -X main.revision=$(git_revision)"

build: fetch_dependencies ## Build
	go build $(add_metadata) -v github.com/mp4096/bart/cmd/bart

install: fetch_dependencies ## Build and install
	go install $(add_metadata) -v github.com/mp4096/bart/cmd/bart

xcompile_linux: fetch_dependencies ## Cross-compile for Linux x64
	env GOOS=linux GOARCH=amd64 go build $(add_metadata) -v github.com/mp4096/bart/cmd/bart

xcompile_win: fetch_dependencies ## Cross-compile for Windows x64
	env GOOS=windows GOARCH=amd64 go build $(add_metadata) -v github.com/mp4096/bart/cmd/bart

xcompile_mac: fetch_dependencies ## Cross-compile for macOS x64
	env GOOS=darwin GOARCH=amd64 go build $(add_metadata) -v github.com/mp4096/bart/cmd/bart

fetch_dependencies: ## Fetch all dependencies
	go get -t ./...

fmt: fetch_dependencies ## Call go fmt in all directories
	go fmt ./...
	gofmt -w -s ./..

delete_previews: ## Delete previews
	find . -type f -name 'bart_preview_*' -delete

vet: fetch_dependencies ## Call go vet in all directories
	go vet ./...

release_binaries: ## Compile binaries for Linux, macOS and Windows; generate digests
	./make_binaries.sh

.PHONY: build install xcompile_linux xcompile_win xcompile_mac \
	fmt delete_previews help vet fetch_dependencies release_binaries

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'
