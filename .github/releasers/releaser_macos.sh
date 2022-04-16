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