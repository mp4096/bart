.DEFAULT_GOAL := build
THIS_FILE := $(lastword $(MAKEFILE_LIST))

build: fetch_dependencies ## Build
	go build -v github.com/mp4096/bart/cmd/bart

install: fetch_dependencies ## Build and install
	go install -v github.com/mp4096/bart/cmd/bart

xcompile_win: fetch_dependencies ## Cross-compile for Windows x64
	env GOOS=windows GOARCH=amd64 go build -v github.com/mp4096/bart/cmd/bart

xcompile_mac: fetch_dependencies ## Cross-compile for macOS x64
	env GOOS=darwin GOARCH=amd64 go build -v github.com/mp4096/bart/cmd/bart

fetch_dependencies: ## Fetch all dependencies
	go get -t ./...

fmt: ## Call go fmt in all directories
	go fmt ./...

delete_previews: ## Delete previews
	find . -type f -name 'bart_preview_*' -delete

vet: ## Call go vet in all directories
	go vet ./...

release_binaries: ## Compile binaries for Linux, macOS and Windows; generate digests
	rm -f release_info.md bart bart.exe
	echo "# Bart binaries\n" >> release_info.md
	echo "git revision:\n" >> release_info.md
	echo "\`\`\`" >> release_info.md
	git rev-parse HEAD >> release_info.md
	echo "\`\`\`\n" >> release_info.md
	echo "Go version:\n" >> release_info.md
	echo "\`\`\`" >> release_info.md
	go version >> release_info.md
	echo "\`\`\`\n" >> release_info.md
	echo "\n## Linux x64\n" >> release_info.md
	$(MAKE) -f $(THIS_FILE) build
	echo "SHA256 digest:\n" >> release_info.md
	echo "\`\`\`" >> release_info.md
	sha256sum bart >> release_info.md
	echo "\`\`\`\n" >> release_info.md
	echo "SHA512 digest:\n" >> release_info.md
	echo "\`\`\`" >> release_info.md
	sha512sum bart >> release_info.md
	echo "\`\`\`\n" >> release_info.md
	tar -cvzf bart_linux_x64.tar.gz bart
	echo "\n## macOS x64\n" >> release_info.md
	$(MAKE) -f $(THIS_FILE) xcompile_mac
	echo "SHA256 digest:\n" >> release_info.md
	echo "\`\`\`" >> release_info.md
	sha256sum bart >> release_info.md
	echo "\`\`\`\n" >> release_info.md
	echo "SHA512 digest:\n" >> release_info.md
	echo "\`\`\`" >> release_info.md
	sha512sum bart >> release_info.md
	echo "\`\`\`\n" >> release_info.md
	tar -cvzf bart_darwin_x64.tar.gz bart
	echo "\n## Windows x64\n" >> release_info.md
	$(MAKE) -f $(THIS_FILE) xcompile_win
	echo "SHA256 digest:\n" >> release_info.md
	echo "\`\`\`" >> release_info.md
	sha256sum bart.exe >> release_info.md
	echo "\`\`\`\n" >> release_info.md
	echo "SHA512 digest:\n" >> release_info.md
	echo "\`\`\`" >> release_info.md
	sha512sum bart.exe >> release_info.md
	echo "\`\`\`" >> release_info.md
	zip bart_windows_x64.zip bart.exe

.PHONY: build install xcompile_win xcompile_mac \
	fmt delete_previews help vet fetch_dependencies release_binaries

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'
