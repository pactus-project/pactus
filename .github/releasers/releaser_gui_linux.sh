#!/bin/bash

# Set –e is used within the Bash to stop execution instantly as a query exits
# while having a non-zero status.
set -e

ROOT_DIR="$(pwd)"
VERSION="$(echo `git -C ${ROOT_DIR} describe --abbrev=0 --tags` | sed 's/^.//')" # "v1.2.3" -> "1.2.3"
BUILD_DIR="${ROOT_DIR}/build"
PACKAGE_NAME="pactus-gui_${VERSION}"
PACKAGE_DIR="${ROOT_DIR}/${PACKAGE_NAME}"

mkdir ${PACKAGE_DIR}

echo "Building the binaries"

cd ${ROOT_DIR}
CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/pactus-daemon ./cmd/daemon
CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/pactus-wallet ./cmd/wallet
go build -ldflags "-s -w" -trimpath -tags gtk -o ${BUILD_DIR}/pactus-gui ./cmd/gtk

# Moving binaries to package directory
echo "Moving binaries"
mv ${BUILD_DIR}/pactus-gui     ${PACKAGE_DIR}/pactus-gui
mv ${BUILD_DIR}/pactus-wallet  ${PACKAGE_DIR}/pactus-wallet
mv ${BUILD_DIR}/pactus-daemon  ${PACKAGE_DIR}/pactus-daemon

echo "Creating archive"
tar -czvf ${ROOT_DIR}/${PACKAGE_NAME}_linux_amd64.tar.gz -p ${PACKAGE_NAME}
