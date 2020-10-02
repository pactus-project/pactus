UNAME := $(shell uname)
GOTOOLS = \
	zombiezen.com/go/capnproto2/...


PACKAGES=$(shell go list ./... | grep -v '/vendor/')
TAGS=-tags 'zarb'
LDFLAGS= -ldflags "-X gitlab.com/zarb-chain/zarb-go/version.GitCommit=`git rev-parse --short=8 HEAD`"
CAPNP_INC = -I$(GOPATH)/src/zombiezen.com/go/capnproto2/std
BLS_PATH= $(shell pwd)/.bls


all: tools build install test test_release

########################################
### Tools & dependencies
tools:
	@echo "Installing tools"
	go get $(GOTOOLS)

bls:
	@echo "Compiling herumi/bls"
	rm -rf .bls && rm -rf .bls
	git clone git://github.com/herumi/mcl.git .bls/mcl && cd .bls/mcl && make -j4
	git clone git://github.com/herumi/bls.git .bls/bls && cd .bls/bls && make lib/libbls384_256.a


########################################
### Build zarb
build:
	go build $(LDFLAGS) $(TAGS) -o build/zarb ./cmd/zarb/

install:
	go install $(LDFLAGS) $(TAGS) ./cmd/zarb

########################################
### Testing
test:
	go test $(PACKAGES)

test_release:
	go test -tags release $(PACKAGES)


########################################
### Docker
docker:
	docker build --tag zarb -f ./containers/Dockerfile .


########################################
### capnp
capnp:
	capnp compile $(CAPNP_INC) -ogo ./www/capnp/zarb.capnp

########################################
### Formatting, linting, and vetting
fmt:
	@go fmt ./...

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: build install test test_release
.PHONY: tools deps bls
.PHONY: fmt metalinter docker
