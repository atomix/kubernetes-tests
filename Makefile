export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ATOMIX_TESTS_VERSION := latest

all: build

build: # @HELP build the source code
build: build-member build-group-member build-partition-group-member
	GOOS=linux GOARCH=amd64 go build -o build/_output/kubernetes-tests ./cmd/kubernetes-tests

build-member:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/membership/_output/bin/atomix-member ./cmd/atomix-member
build-group-member:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/group-membership/_output/bin/atomix-group-member ./cmd/atomix-group-member
build-partition-group-member:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/partition-group-membership/_output/bin/atomix-partition-group-member ./cmd/atomix-partition-group-member

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
images: image-member image-group-member image-partition-group-member

image-member: build-member
	docker build . -f build/membership/Dockerfile -t atomix/test-member:latest
image-group-member: build-group-member
	docker build . -f build/group-membership/Dockerfile -t atomix/test-group-member:latest
image-partition-group-member: build-partition-group-member
	docker build . -f build/partition-group-membership/Dockerfile -t atomix/test-partition-group-member:latest

kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image atomix/test-member:latest
	kind load docker-image atomix/test-group-member:latest
	kind load docker-image atomix/test-partition-group-member:latest

push: # @HELP push kubernetes-tests Docker image
	docker push atomix/kubernetes-tests:${ATOMIX_TESTS_VERSION}
