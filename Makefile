PACKAGES=$(shell go list ./... | grep -v 'tests')
HERUMI= $(shell pwd)/.herumi
CGO_LDFLAGS=CGO_LDFLAGS="-L$(HERUMI)/bls/lib -lbls384_256 -lm -lstdc++ -g -O2"
BUILD_LDFLAGS= -ldflags "-X github.com/zarbchain/zarb-go/version.build=`git rev-parse --short=8 HEAD`"
RELEASE_LDFLAGS= -ldflags "-s -w"
RELEASE_GUI_LDFLAGS= -ldflags "-s -w -H windowsgui"

ifneq (,$(filter $(OS),Windows_NT MINGW64))
EXE = .exe
endif

all: install test

########################################
### Tools needed for development
devtools:
	@echo "Installing devtools"
	go install zombiezen.com/go/capnproto2/capnpc-go@v2.18
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.10
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.10
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install github.com/bufbuild/buf/cmd/buf@v1.3
	go install github.com/rakyll/statik@v0.1
	go install github.com/gordonklaus/ineffassign@latest

herumi:
	@if [ ! -d $(HERUMI) ]; then \
		git clone --recursive https://github.com/herumi/bls.git $(HERUMI)/bls && cd $(HERUMI)/bls && make minimized_static; \
	fi

########################################
### Building
build:
	go build $(BUILD_LDFLAGS) -o ./build/zarb-daemon$(EXE) ./cmd/daemon
	go build $(BUILD_LDFLAGS) -o ./build/zarb-wallet$(EXE) ./cmd/wallet

release: herumi
	$(CGO_LDFLAGS) go build $(RELEASE_LDFLAGS) -o ./build/zarb-daemon$(EXE) ./cmd/daemon
	$(CGO_LDFLAGS) go build $(RELEASE_LDFLAGS) -o ./build/zarb-wallet$(EXE) ./cmd/wallet

build_gui:
	go build $(BUILD_LDFLAGS) -tags gtk -o ./build/zarb-gui$(EXE) ./cmd/gtk

release_gui: herumi
	$(CGO_LDFLAGS) go build $(RELEASE_GUI_LDFLAGS) -tags gtk -o ./build/zarb-gui$(EXE) ./cmd/gtk


########################################
### Testing
unit_test:
	go test $(PACKAGES)

test:
	go test ./... -covermode=atomic

test_race:
	go test ./... --race

########################################
### Docker
docker:
	docker build --tag zarb .

########################################
### capnp and proto
capnp:
	capnp compile \
		-ogo ./www/capnp/zarb.capnp

proto:
	cd www/grpc/ && buf generate --path ./proto/zarb.proto --path proto/payloads.proto
	# Generate static assets for OpenAPI UI
	cd www/grpc/ && statik -m -f -src third_party/OpenAPI/

########################################
### Formatting, linting, and vetting
fmt:
	gofmt -s -w .
	golangci-lint run -e "SA1019" \
		--timeout=5m0s \
		--enable=gofmt \
		--enable=unconvert \
		--enable=unparam \
		--enable=revive \
		--enable=asciicheck \
		--enable=misspell \
		--enable=asciicheck \
		--enable=decorder \
		--enable=depguard \
		--enable=nilerr \
		--enable=gosec \
		--enable=gocyclo
	ineffassign ./...

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: build install release
.PHONY: test unit_test test_race
.PHONY: devtools herumi capnp proto
.PHONY: fmt docker
