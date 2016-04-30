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
