# Update Dependencies

This document outlines the steps to update dependencies in the Pactus project repository
to their latest versions.

## Update Go

First, update the Go version to the [latest release](https://go.dev/doc/install) in [go.mod](../go.mod).
Ensure the Golang version is also updated in the [Dockerfile](../Dockerfile).

## Update Dependencies

To update Go dependencies to their latest versions, use the following commands:

```sh
go get -u ./...
go mod tidy
```

Once all packages are updated, run `make build`, `make test`, and `make build_gui`
to ensure no existing functionality is broken.
If any packages introduce breaking changes or are deprecated,
update the code to use the new methods or replace the package entirely.

## Update Dev Tools

Next, update development tools such as `golangci-lint`, `buf`, and others.
Refer to the [Makefile](../Makefile) and locate the `devtools` target.
Replace each tool with its latest version.

## Update GitHub Workflows

Navigate to the [workflows](../.github/workflows) directory and
update outdated GitHub actions to their latest versions.
You can find the latest versions by searching for the action name on GitHub.

## Buf Dependencies

We use [buf](https://buf.build/explore) to generate code from proto files.
Update the buf plugins in [buf.gen.yaml](../www/grpc/buf/buf.gen.yaml) to their latest versions.

## Example Pull Request

For reference, see this example pull request for updates and commit format:
https://github.com/pactus-project/pactus/pull/1202
