#!/bin/bash

# Set –e is used within the Bash to stop execution instantly as a query exits
# while having a non-zero status.
set -e

ROOT_DIR="$(pwd)"
VERSION="$(echo `git -C ${ROOT_DIR} describe --abbrev=0 --tags` | sed 's/^.//')" # "v1.2.3" -> "1.2.3"
PACKAGE_NAME="pactus-cli_${VERSION}"


# https://go.dev/doc/install/source#environment

for OS_ARCH in \
     "linux amd64" "linux arm64" \
     "android arm64" \
     "freebsd amd64" "freebsd arm" \
     "darwin amd64" "darwin arm64" \
     "windows 386" "windows amd64"; do

    PAIR=($OS_ARCH);
    OS=${PAIR[0]};
    ARCH=${PAIR[1]};

    cd ${ROOT_DIR}

    PACKAGE_NAME_OS=${PACKAGE_NAME}_${OS}_${ARCH}
    BUILD_DIR=${ROOT_DIR}/build/${PACKAGE_NAME_OS}

    if [ $OS = "windows" ]; then
        EXE=".exe"
    fi

    CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/${PACKAGE_NAME}/pactus-daemon${EXE} ./cmd/daemon
    CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/${PACKAGE_NAME}/pactus-wallet${EXE} ./cmd/wallet
    CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/${PACKAGE_NAME}/pactus-shell${EXE} ./cmd/shell

    cd ${BUILD_DIR}
    if [ $OS = "windows" ]; then
        zip -r ${PACKAGE_NAME_OS}.zip ${PACKAGE_NAME}
        mv ${PACKAGE_NAME_OS}.zip ${ROOT_DIR}
    else
        tar -czvf ${PACKAGE_NAME_OS}.tar.gz -p ${PACKAGE_NAME}
        mv ${PACKAGE_NAME_OS}.tar.gz ${ROOT_DIR}
    fi
done
