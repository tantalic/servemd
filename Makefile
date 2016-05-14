build: bindata-assetfs
	go build

docker: bindata-assetfs
	docker run --rm -v "$(GOPATH)":/go -w /go/src/tantalic.com/servemd blang/golang-alpine go build -v
	docker build -t servemd .
	rm servemd

bindata-assetfs:
	go-bindata-assetfs -ignore ".DS_Store" assets/...

run:
	go-bindata-assetfs -ignore ".DS_Store" -debug assets/...
	go run *.go

install: bindata-assetfs
	go install

osx: bindata-assetfs
	env GOOS=darwin GOARCH=amd64 go build -o build/servemd-darwin_amd64
	env GOOS=darwin GOARCH=386   go build -o build/servemd-darwin_386

linux: bindata-assetfs
	env GOOS=linux GOARCH=amd64 go build -o build/servemd-linux_amd64
	env GOOS=linux GOARCH=386   go build -o build/servemd-linux_386
	env GOOS=linux GOARCH=arm   go build -o build/servemd-linux_arm
	env GOOS=linux GOARCH=arm64 go build -o build/servemd-linux_arm64

freebsd: bindata-assetfs
	env GOOS=freebsd GOARCH=amd64 go build -o build/servemd-freebsd_amd64
	env GOOS=freebsd GOARCH=386   go build -o build/servemd-freebsd_386
	env GOOS=freebsd GOARCH=arm   go build -o build/servemd-freebsd_arm

openbsd: bindata-assetfs
	env GOOS=openbsd GOARCH=amd64 go build -o build/servemd-openbsd_amd64
	env GOOS=openbsd GOARCH=386   go build -o build/servemd-openbsd_386
	env GOOS=openbsd GOARCH=arm   go build -o build/servemd-openbsd_arm

netbsd: bindata-assetfs
	env GOOS=netbsd GOARCH=amd64 go build -o build/servemd-netbsd_amd64
	env GOOS=netbsd GOARCH=386   go build -o build/servemd-netbsd_386
	env GOOS=netbsd GOARCH=arm   go build -o build/servemd-netbsd_arm

dragonfly: bindata-assetfs
	env GOOS=dragonfly GOARCH=amd64 go build -o build/servemd-dragonfly_amd64

windows: bindata-assetfs
	env GOOS=windows GOARCH=amd64 go build -o build/servemd-windows_amd64
	env GOOS=windows GOARCH=386   go build -o build/servemd-windows_386

all: osx linux freebsd openbsd netbsd dragonfly windows
