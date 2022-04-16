#!/bin/bash

set -e

ROOT_DIR="$(pwd)"
VERSION="$(echo `git -C ${ROOT_DIR} describe --abbrev=0 --tags` | sed 's/^.//')" # "v1.2.3" -> "1.2.3"
BUILD_DIR="${ROOT_DIR}/build"
PACKAGE_NAME="zarb-${VERSION}"
PACKAGE_DIR="${ROOT_DIR}/${PACKAGE_NAME}"

echo "VERSION: $VERSION"
echo "ROOT_DIR: $ROOT_DIR"
echo "PACKAGE_DIR: $PACKAGE_DIR"

mkdir ${PACKAGE_DIR}

echo "Building the binaries"

cd ${ROOT_DIR}
make herumi
export CGO_LDFLAGS="-L.herumi/bls/lib -lbls384_256 -lm -lstdc++ -g -O2"
go build -ldflags "-s -w" -o ${BUILD_DIR}/zarb-daemon ./cmd/daemon
go build -ldflags "-s -w" -o ${BUILD_DIR}/zarb-wallet ./cmd/wallet
go build -ldflags "-s -w" -tags gtk -o ${BUILD_DIR}/zarb-gui ./cmd/gtk

# Moving binaries to package directory
echo "Moving binaries"
mv ${BUILD_DIR}/zarb-gui     ${PACKAGE_DIR}/zarb-gui
mv ${BUILD_DIR}/zarb-wallet  ${PACKAGE_DIR}/zarb-wallet
mv ${BUILD_DIR}/zarb-daemon  ${PACKAGE_DIR}/zarb-daemon

echo "Creating archive"
tar -czvf ${ROOT_DIR}/${PACKAGE_NAME}-linux-x86_64.tar.gz -p ${PACKAGE_NAME}
