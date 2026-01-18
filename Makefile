BINARY_NAME=devkit
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X devkit/pkg/version.Version=${VERSION} -X devkit/pkg/version.BuildTime=${BUILD_TIME}"

.PHONY: build build-all clean test lint install

build:
	@mkdir -p bin
	go build ${LDFLAGS} -o bin/${BINARY_NAME} .

build-macos:
	@mkdir -p bin
	@echo "Building for macOS (current architecture)..."
	go build ${LDFLAGS} -o bin/${BINARY_NAME}-macos .
	@cp bin/${BINARY_NAME}-macos bin/devcli
	@chmod +x bin/devcli
	@echo "Build complete: bin/${BINARY_NAME}-macos and bin/devcli"

build-macos-all:
	@mkdir -p bin
	@echo "Building for macOS (amd64 and arm64)..."
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-arm64 .
	@echo "Build complete: bin/${BINARY_NAME}-darwin-amd64 and bin/${BINARY_NAME}-darwin-arm64"

build-all:
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-windows-amd64.exe .

clean:
	rm -rf bin/

test:
	go test -v ./...

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

install:
	go install ${LDFLAGS} .
