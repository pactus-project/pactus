UNAME := $(shell uname)
GOTOOLS = \
	zombiezen.com/go/capnproto2/... \
	github.com/golangci/golangci-lint/cmd/golangci-lint@v1.33.0


PACKAGES=$(shell go list ./... | grep -v 'tests')
TAGS=-tags 'zarb'
HERUMI= $(shell pwd)/.herumi
CGO_LDFLAGS=CGO_LDFLAGS="-L$(HERUMI)/mcl/lib -L$(HERUMI)/bls/lib -lmcl -lbls384_256 -lm -lstdc++"
LDFLAGS= -ldflags "-X github.com/zarbchain/zarb-go/version.GitCommit=`git rev-parse --short=8 HEAD`"
CAPNP_INC = -I$(GOPATH)/src/zombiezen.com/go/capnproto2/std



all: tools build install test

########################################
### Tools & dependencies
tools:
	@echo "Installing tools"
	go get $(GOTOOLS)


bls:
	@echo "Compiling bls"
	rm -rf $(HERUMI)
	git clone git://github.com/herumi/mcl.git $(HERUMI)/mcl && cd $(HERUMI)/mcl && make lib/libmcl.a
	git clone git://github.com/herumi/bls.git $(HERUMI)/bls && cd $(HERUMI)/bls && make minimized_static


########################################
### Build zarb
build:
	go build $(LDFLAGS) $(TAGS) -o build/zarb ./cmd/zarb/

install: fmt
	go install $(LDFLAGS) $(TAGS) ./cmd/zarb

build_with_bls:
	$(CGO_LDFLAGS) go build $(LDFLAGS) $(TAGS) -o build/zarb ./cmd/zarb/

########################################
### Testing
unit_test:
	go test $(PACKAGES)

test:
	go test ./...

test_with_bls:
	$(CGO_LDFLAGS) go test ./...

########################################
### Docker
docker:
	docker build --tag zarb -f ./containers/Dockerfile .


########################################
### capnp
capnp: tools
	capnp compile $(CAPNP_INC) -ogo ./www/capnp/zarb.capnp

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
.PHONY: tools deps bls
.PHONY: fmt metalinter docker
