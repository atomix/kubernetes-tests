export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ATOMIX_TESTS_VERSION := latest

all: build

build: # @HELP build the source code
build: build-peer
	GOOS=linux GOARCH=amd64 go build -o build/_output/kubernetes-tests ./cmd/kubernetes-tests

build-peer:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/peer/_output/bin/atomix-peer ./cmd/atomix-peer

test: # @HELP run the unit tests and source code validation
test: build license_check linters
	go test github.com/atomix/kubernetes-tests/...

linters: # @HELP examines Go source code and reports coding problems
	golangci-lint run

license_check: # @HELP examine and ensure license headers exist
	./build/licensing/boilerplate.py -v

proto: # @HELP build Protobuf/gRPC generated types
proto:
	docker run -it -v `pwd`:/go/src/github.com/atomix/kubernetes-tests \
		-w /go/src/github.com/atomix/kubernetes-tests \
		--entrypoint build/bin/compile_protos.sh \
		onosproject/protoc-go:stable

images: # @HELP build kubernetes-tests Docker image
images: image-peer

image-peer: build-peer
	docker build . -f build/peer/Dockerfile -t atomix/test-peer:${ATOMIX_TESTS_VERSION}

kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image atomix/test-peer:${ATOMIX_TESTS_VERSION}

push: # @HELP push kubernetes-tests Docker image
	docker push atomix/kubernetes-tests:${ATOMIX_TESTS_VERSION}
