PACKAGES=$(shell go list ./... | grep -v 'tests' | grep -v 'grpc/gen')

ifneq (,$(filter $(OS),Windows_NT MINGW64))
EXE = .exe
endif

# Handle sed differences between macOS and Linux
ifeq ($(shell uname),Darwin)
  SED_CMD = sed -i ''
else
  SED_CMD = sed -i
endif

all: build test

########################################
### Tools needed for development
devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.2
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.27
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.27
	go install github.com/NathanBaulch/protoc-gen-cobra@v1.2.1
	go install github.com/pactus-project/protoc-gen-doc/cmd/protoc-gen-doc@v0.0.0-20250428081805-49e8e8103faa
	go install github.com/bufbuild/buf/cmd/buf@v1.55
	go install mvdan.cc/gofumpt@latest
	go install github.com/pacviewer/jrpc-gateway/protoc-gen-jrpc-gateway@v0.6

########################################
### Building
build:
	go build -o ./build/pactus-daemon$(EXE) ./cmd/daemon
	go build -o ./build/pactus-wallet$(EXE) ./cmd/wallet
	go build -o ./build/pactus-shell$(EXE)  ./cmd/shell

build_race:
	go build -race -o ./build/pactus-daemon$(EXE) ./cmd/daemon
	go build -race -o ./build/pactus-wallet$(EXE) ./cmd/wallet

build_gui:
	go build -tags gtk -o ./build/pactus-gui$(EXE) ./cmd/gtk

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
	docker build --tag pactus .

########################################
### proto

# This target works only on Unix-like terminals.
proto:
	rm -rf www/grpc/gen
	cd www/grpc && buf generate --template ./buf/buf.gen.yaml --config ./buf/buf.yaml ./proto

proto-check:
	cd www/grpc && buf lint --config ./buf/buf.yaml

proto-format:
	cd www/grpc && buf format --config ./buf/buf.yaml -w


########################################
### Formatting the code
fmt:
	gofumpt -l -w .

check:
	golangci-lint run --build-tags "${BUILD_TAG}" --timeout=20m0s

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: build build_gui
.PHONY: test unit_test test_race
.PHONY: proto proto-format proto-check
.PHONY: devtools fmt check docker
