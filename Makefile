UNAME := $(shell uname)
GOTOOLS = \
	zombiezen.com/go/capnproto2/...


PACKAGES=$(shell go list ./... | grep -v '/vendor/')
TAGS=-tags 'zarb'
LDFLAGS= -ldflags "-X github.com/zarbchain/zarb-go/version.GitCommit=`git rev-parse --short=8 HEAD`"
CAPNP_INC = -I$(GOPATH)/src/zombiezen.com/go/capnproto2/std
BLS_PATH= $(shell pwd)/.bls


all: tools build install test

########################################
### Tools & dependencies
tools:
	@echo "Installing tools"
	GO111MODULE=off go get $(GOTOOLS)


########################################
### Build zarb
build:
	go build $(LDFLAGS) $(TAGS) -o build/zarb ./cmd/zarb/

install:
	go install $(LDFLAGS) -ldflags '-s' $(TAGS) ./cmd/zarb

########################################
### Testing
test:
	go test $(PACKAGES)


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
	@go fmt ./...

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: build install test
.PHONY: tools deps bls
.PHONY: fmt metalinter docker
