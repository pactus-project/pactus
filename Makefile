PACKAGES=$(shell go list ./... | grep -v 'tests' | grep -v 'grpc/gen')

ifneq (,$(filter $(OS),Windows_NT MINGW64))
EXE = .exe
RM = del /q
else
RM = rm -rf
endif


all: build test

########################################
### Tools needed for development
devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.20
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.20
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4
	go install github.com/NathanBaulch/protoc-gen-cobra@v1.2
	go install github.com/pactus-project/protoc-gen-doc/cmd/protoc-gen-doc@master
	go install github.com/bufbuild/buf/cmd/buf@v1.32
	go install mvdan.cc/gofumpt@latest
	go install github.com/rakyll/statik@v0.1.7
	go install github.com/pacviewer/jrpc-gateway/protoc-gen-jrpc-gateway@v0.3.2
	go install github.com/pacviewer/jrpc-gateway/protoc-gen-jrpc-doc/cmd/protoc-gen-jrpc-doc@v0.1.7

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
proto:
	$(RM) www/grpc/gen
	cd www/grpc/buf && buf generate --template buf.gen.yaml ../proto
# Generate static assets for Swagger-UI
	cd www/grpc/ && statik -m -f -src swagger-ui/

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
.PHONY: devtools proto
.PHONY: fmt check docker
