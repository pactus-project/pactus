#  Updating Project Dependencies

This document is about how to update the Pactus project repository to latest version.

### Packages

First of all you need to update golang dependencies to latest version using this commands:

```sh
go get -u ./...
go mod tidy
```
Once all packages got updated, make sure you run `make build`, `make test` and `make build_gui` commands to make sure
none of previous behaviors are broken. If any packages had breaking changes or some of them are deprecated, you need to
update the code and use new methods or use another package.


### Dev tools

After packages, you need to update dev tools such as `golangci-lint`, `buf`, etc.

You can go to root [make file](../Makefile) and find all dev tools on:

```makefile
### Tools needed for development
devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@vx
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@vx
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@vx
	go install google.golang.org/protobuf/cmd/protoc-gen-go@vx
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@vx
	go install github.com/NathanBaulch/protoc-gen-cobra@vx
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@vx
	go install github.com/bufbuild/buf/cmd/buf@vx
	go install mvdan.cc/gofumpt@vx
	go install github.com/rakyll/statik@vx
	go install github.com/pacviewer/jrpc-gateway/protoc-gen-jrpc-gateway@vx
	go install ....
```

You have to find latest version of dev tools and replace them here.

> Note: consider breaking changes and deprecated packages for devtools too.

### Go version

You have to update the go version to latest release in [go.mod](../go.mod).
Make sure you are updating version of Golang on [Dockerfile](../Dockerfile).

> Note: you must run `make build` and `make build_gui` after this change to make sure everything works smoothly.

### CI/CD and GitHub workflows

You need to go to [workflows](../.github/workflows) directory and update old GitHub actions to latest version.
You can find the latest version by searching the action name on GitHub.

### Example Pull Request

Here is an example pull request to find out what you need to update and how to set commit message:
https://github.com/pactus-project/pactus/pull/1202

