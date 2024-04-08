GO_FILES = $(shell find . -name "*.go")
PKG_VERSION = $(shell jq -r .version package.json)

.PHONY: build
build: build-darwin-arm64 build-darwin-x64 build-linux-arm64 build-linux-x64 build-windows-arm64 build-windows-x64

.PHONY: build-darwin-arm64
build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 PKG_VERSION=$(PKG_VERSION) go build -o ./dist/moldable-darwin-arm64

.PHONY: build-linux-arm64
build-linux-arm64:
	GOOS=linux GOARCH=arm64 PKG_VERSION=$(PKG_VERSION) go build -o ./dist/moldable-linux-arm64

.PHONY: build-windows-arm64
build-windows-arm64:
	GOOS=windows GOARCH=arm64 PKG_VERSION=$(PKG_VERSION) go build -o ./dist/moldable-windows-arm64.exe

.PHONY: build-darwin-x64
build-darwin-x64:
	GOOS=darwin GOARCH=amd64 PKG_VERSION=$(PKG_VERSION) go build -o ./dist/moldable-darwin-x64

.PHONY: build-linux-x64
build-linux-x64:
	GOOS=linux GOARCH=amd64 PKG_VERSION=$(PKG_VERSION) go build -o ./dist/moldable-linux-x64

.PHONY: build-windows-x64
build-windows-x64:
	GOOS=windows GOARCH=amd64 PKG_VERSION=$(PKG_VERSION) go build -o ./dist/moldable-windows-x64.exe

test: $(GO_FILES) go.mod go.sum
	go test -v ./...

lint: $(GO_FILES) go.mod go.sum
	staticcheck ./...

.PHONY: clean
clean:
	go clean -testcache -r

format: $(GO_FILES) go.mod go.sum
	go fmt ./...

dev: $(GO_FILES) go.mod go.sum
	MOCKDIR=__mocks__ PKG_VERSION=$(PKG_VERSION) go run ./main.go 

.PHONY: go-install
go-install: 
	go mod download 
