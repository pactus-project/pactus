PACKAGES=$(shell go list ./... | grep -v 'tests')

TAGS="zarb"
VTAG=$(shell git describe --tags --abbrev=0)
HERUMI= .herumi
LDFLAGS= -ldflags "-X github.com/zarbchain/zarb-go/version.build=`git rev-parse --short=8 HEAD` -X github.com/zarbchain/zarb-go/version.semVer=$(VTAG)"


_OS=$(shell go env GOOS)
_ARCH=$(shell go env GOARCH)


LIB_DIR=$(shell pwd)/$(HERUMI)/bls/lib/$(_OS)/$(_ARCH)
CGO_LDFLAGS=CGO_LDFLAGS="-L$(LIB_DIR) -lbls384_256 -lm -lstdc++ -g -O2"
CGO_LDFLAGS_Default=$(shell go env CGO_LDFLAGS)
_EXT=$(shell go env GOEXE)

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
	@echo "cloning bls"
	rm -rf $(HERUMI)
	git clone https://github.com/herumi/bls-go-binary.git $(HERUMI)


########################################
### Building
build:
	go build $(LDFLAGS) -tags $(TAGS) -o build/zarb ./cmd/zarb/

install:
	go install $(LDFLAGS) -tags $(TAGS) ./cmd/zarb

build_with_bls:
	go env -w $(CGO_LDFLAGS)
	go build $(LDFLAGS) -tags $(TAGS) -o build/zarb ./cmd/zarb/
	go env -w CGO_LDFLAGS="$(CGO_LDFLAGS_Default)"

########################################
### Testing
unit_test:
	go test $(PACKAGES)

test:
	go test ./... -covermode=atomic

test_race:
	go test ./... --race

test_bls:
	$(CGO_LDFLAGS) go test ./...

########################################
### Docker
docker:
	docker build --tag zarb -f ./containers/Dockerfile .

########################################
### capnp and proto
capnp:
	capnp compile \
		-I $(GOPATH)/src/zombiezen.com/go/capnproto2/std \
		-ogo ./www/capnp/zarb.capnp

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

fmt_bls:
	$(CGO_LDFLAGS) go vet ./...
	$(CGO_LDFLAGS) gofmt -s -w .
	$(CGO_LDFLAGS) golangci-lint run -e "SA1019"

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: build install build_bls build_bls_release
.PHONY: test unit_test test_race test_bls
.PHONY: tools bls capnp proto
.PHONY: fmt fmt_bls docker