#!/bin/bash

set -e

ROOT_DIR="$(pwd)"
VERSION="$(echo `git -C ${ROOT_DIR} describe --abbrev=0 --tags` | sed 's/^.//')" # "v1.2.3" -> "1.2.3"
BUILD_DIR="${ROOT_DIR}/build"
PACKAGE_NAME="pactus-${VERSION}"
PACKAGE_DIR="${ROOT_DIR}/${PACKAGE_NAME}"

echo "VERSION: $VERSION"
echo "ROOT_DIR: $ROOT_DIR"
echo "PACKAGE_DIR: $PACKAGE_DIR"

mkdir ${PACKAGE_DIR}

echo "Building the binaries"

cd ${ROOT_DIR}
make herumi
export CGO_LDFLAGS="-L.herumi/bls/lib -lbls384_256 -lm -g -O2"
go build -ldflags "-s -w" -o ${BUILD_DIR}/pactus-daemon ./cmd/daemon
go build -ldflags "-s -w" -o ${BUILD_DIR}/pactus-wallet ./cmd/wallet
go build -ldflags "-s -w" -tags gtk -o ${BUILD_DIR}/pactus-gui ./cmd/gtk

# Moving binaries to package directory
echo "Moving binaries"
mv ${BUILD_DIR}/pactus-gui     ${PACKAGE_DIR}/pactus-gui
mv ${BUILD_DIR}/pactus-wallet  ${PACKAGE_DIR}/pactus-wallet
mv ${BUILD_DIR}/pactus-daemon  ${PACKAGE_DIR}/pactus-daemon

echo "Creating archive"
tar -czvf ${ROOT_DIR}/${PACKAGE_NAME}-linux-x86_64.tar.gz -p ${PACKAGE_NAME}
