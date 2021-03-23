UNAME := $(shell uname)
PACKAGES=$(shell go list ./... | grep -v 'tests')
TAGS=-tags 'zarb'
HERUMI= $(shell pwd)/.herumi
CGO_LDFLAGS=CGO_LDFLAGS="-L$(HERUMI)/bls/lib -lbls384_256 -lm -lstdc++"
LDFLAGS= -ldflags "-X github.com/zarbchain/zarb-go/version.GitCommit=`git rev-parse --short=8 HEAD`"
CAPNP_INC = -I$(GOPATH)/src/zombiezen.com/go/capnproto2/std
PROTO_INC = -I. -I$(GOPATH)/src/github.com/googleapis/googleapis


all: tools build install test

########################################
### Tools & dependencies
tools:
	@echo "Installing tools"
	GO111MODULE=off go get zombiezen.com/go/capnproto2/...
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.38.0
	go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.3.0
	go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.3.0
	go get google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	go get github.com/bufbuild/buf/cmd/buf@v0.39.1
	go get github.com/rakyll/statik@v0.1.7
	go mod tidy

bls:
	@echo "Compiling bls"
	rm -rf $(HERUMI)
	git clone --recursive git://github.com/herumi/bls.git $(HERUMI)/bls && cd $(HERUMI)/bls && make minimized_static


########################################
### Build zarb
build:
	go build $(LDFLAGS) $(TAGS) -o build/zarb ./cmd/zarb/

install:
	go install $(LDFLAGS) $(TAGS) ./cmd/zarb

build_with_bls:
	$(CGO_LDFLAGS) go build $(LDFLAGS) $(TAGS) -o build/zarb ./cmd/zarb/

########################################
### Testing
unit_test:
	go test $(PACKAGES)

test:
	go test ./... -covermode=atomic

test_race:
	go test ./... --race

test_with_bls:
	$(CGO_LDFLAGS) go test ./...

########################################
### Docker
docker:
	docker build --tag zarb -f ./containers/Dockerfile .

########################################
### capnp and proto
capnp:
	capnp compile $(CAPNP_INC) -ogo ./www/capnp/zarb.capnp

proto:
	cd www/grpc/ && buf generate  --path ./proto/zarb.proto
	# Generate static assets for OpenAPI UI
	cd www/grpc/ && statik -m -f -src third_party/OpenAPI/

########################################
### Formatting, linting, and vetting
fmt:
	@go vet ./...
	@gofmt -s -w .
	@golangci-lint run -e "SA1019"

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: build install test
.PHONY: tools deps bls capnp proto
.PHONY: fmt metalinter docker
