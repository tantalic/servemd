VERSION := $(shell git describe --always --dirty)
GOVERSION := 1.9.1
PKG := "toolhouse.com/projects-site"

.PHONY: help
.DEFAULT_GOAL := help

run: ## Build and run
	go-bindata-assetfs -ignore ".DS_Store" -debug assets/...
	go run *.go

build: bindata-assetfs ## Build for the current arch with the local go version installation
	go build

docker: linux ## Builds docker image
	docker build -t tantalic/servemd:$(VERSION) .

bindata-assetfs:
	go-bindata-assetfs -ignore ".DS_Store" assets/...

install: bindata-assetfs ## Install on the local machine
	go install

osx: bindata-assetfs
	docker run --env GOOS=darwin --env GOARCH=amd64 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/servemd-darwin_amd64
	docker run --env GOOS=darwin --env GOARCH=386   --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/servemd-darwin_386

linux: bindata-assetfs
	docker run --env GOOS=linux --env GOARCH=amd64 --env CGO_ENABLED=0 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -a -tags netgo -ldflags '-w' -o build/servemd-linux_amd64
	docker run --env GOOS=linux --env GOARCH=386   --env CGO_ENABLED=0 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -a -tags netgo -ldflags '-w' -o build/servemd-linux_386
	docker run --env GOOS=linux --env GOARCH=arm   --env CGO_ENABLED=0 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -a -tags netgo -ldflags '-w' -o build/servemd-linux_arm
	docker run --env GOOS=linux --env GOARCH=arm64 --env CGO_ENABLED=0 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -a -tags netgo -ldflags '-w' -o build/servemd-linux_arm64	

freebsd: bindata-assetfs
	docker run --env GOOS=freebsd --env GOARCH=amd64 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/servemd-freebsd_amd64
	docker run --env GOOS=freebsd --env GOARCH=386   --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/servemd-freebsd_386
	docker run --env GOOS=freebsd --env GOARCH=arm   --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/servemd-freebsd_arm

openbsd: bindata-assetfs
	docker run --env GOOS=openbsd --env GOARCH=amd64 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/servemd-openbsd_amd64
	docker run --env GOOS=openbsd --env GOARCH=386   --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/servemd-openbsd_386
	docker run --env GOOS=openbsd --env GOARCH=arm   --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/servemd-openbsd_arm


netbsd: bindata-assetfs
	docker run --env GOOS=netbsd --env GOARCH=amd64 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/servemd-netbsd_amd64
	docker run --env GOOS=netbsd --env GOARCH=386   --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/servemd-netbsd_386
	docker run --env GOOS=netbsd --env GOARCH=arm   --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/servemd-netbsd_arm

dragonfly: bindata-assetfs
	docker run --env GOOS=dragonfly --env GOARCH=amd64 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/servemd-dragonfly_amd64

windows: bindata-assetfs
	docker run --env GOOS=windows --env GOARCH=amd64 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/servemd-windows_amd64
	docker run --env GOOS=windows --env GOARCH=386   --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/servemd-windows_386

all: osx linux freebsd openbsd netbsd dragonfly windows ## Build for all supported platforms

help: ## List available make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
