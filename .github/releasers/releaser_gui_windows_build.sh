#!/bin/bash

set -e

ROOT_DIR="$(pwd)"
BUILD_DIR="${ROOT_DIR}/build"

# Copy Windows resources file for application icon
cp ${ROOT_DIR}/.github/releasers/windows/rsrc_windows_amd64.syso ${ROOT_DIR}/cmd/gtk/rsrc_windows_amd64.syso

# This fixes a bug in pkgconfig: invalid flag in pkg-config --libs: -Wl,-luuid
sed -i -e 's/-Wl,-luuid/-luuid/g' /mingw64/lib/pkgconfig/gdk-3.0.pc

CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/unsigned/pactus-daemon.exe        ./cmd/daemon
CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/unsigned/pactus-wallet.exe        ./cmd/wallet
CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/unsigned/pactus-shell.exe         ./cmd/shell
go build -ldflags "-s -w -H windowsgui" -trimpath -tags gtk -o ${BUILD_DIR}/unsigned/pactus-gui.exe ./cmd/gtk
