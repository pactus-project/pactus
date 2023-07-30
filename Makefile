PACKAGES=$(shell go list ./... | grep -v 'tests')
BUILD_LDFLAGS= -ldflags "-X github.com/pactus-project/pactus/version.build=`git rev-parse --short=8 HEAD`"

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
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.12
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.12
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install github.com/bufbuild/buf/cmd/buf@v1.25.0
	go install github.com/rakyll/statik@v0.1

########################################
### Building
build:
	go build $(BUILD_LDFLAGS) -o ./build/pactus-daemon$(EXE) ./cmd/daemon
	go build $(BUILD_LDFLAGS) -o ./build/pactus-wallet$(EXE) ./cmd/wallet

build_gui:
	go build $(BUILD_LDFLAGS) -tags gtk -o ./build/pactus-gui$(EXE) ./cmd/gtk

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
	cd www/grpc/ && $(RM) gen && buf generate proto

# Generate static assets for Swagger-UI
	cd www/grpc/ && statik -m -f -src swagger-ui/

########################################
### Formatting, linting, and vetting
fmt:
	gofmt -s -w .

check:
	golangci-lint run \
		--build-tags "${BUILD_TAG}" \
		--timeout=20m0s \
		--enable=gofmt \
		--enable=unconvert \
		--enable=unparam \
		--enable=asciicheck \
		--enable=misspell \
		--enable=revive \
		--enable=decorder \
		--enable=reassign \
		--enable=usestdlibvars \
		--enable=nilerr \
		--enable=gosec \
		--enable=exportloopref \
		--enable=whitespace \
		--enable=goimports \
		--enable=gocyclo \
		--enable=lll

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: build build_gui
.PHONY: test unit_test test_race
.PHONY: devtools proto
.PHONY: fmt check docker
