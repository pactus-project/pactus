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

echo "Building the binaries"

cd ${ROOT_DIR}
make herumi
export CGO_LDFLAGS="-L.herumi/bls/lib -lbls384_256 -lm -g -O2"
go build -ldflags "-s -w" -o ${BUILD_DIR}/zarb-daemon ./cmd/daemon
go build -ldflags "-s -w" -o ${BUILD_DIR}/zarb-wallet ./cmd/wallet
go build -ldflags "-s -w" -tags gtk -o ${BUILD_DIR}/zarb-gui ./cmd/gtk



echo "Installing gtk-mac-bundler"
git clone https://gitlab.gnome.org/GNOME/gtk-mac-bundler.git
cd gtk-mac-bundler
make install
BUNDLER=${HOME}/.local/bin/gtk-mac-bundler
cd -

echo "Bundling the GUI package"
GUI_BUNDLE=${ROOT_DIR}/gui-bundle
mkdir ${GUI_BUNDLE}

cp ${BUILD_DIR}/pactus-gui                  ${GUI_BUNDLE}
cp ${ROOT_DIR}/.github/releasers/macos/*  ${GUI_BUNDLE}

# https://stackoverflow.com/questions/21242932/sed-i-may-not-be-used-with-stdin-on-mac-os-x
sed -i '' "s/%SHORTVERSION%/${VERSION}/"     ${GUI_BUNDLE}/Info.plist
sed -i '' "s/%VERSION%/Version ${VERSION}/"  ${GUI_BUNDLE}/Info.plist

export GUI_BUNDLE
export ROOT_DIR

${BUNDLER} ${GUI_BUNDLE}/gui.bundle


echo "Creating dmg"
# https://github.com/create-dmg/create-dmg
create-dmg \
  --volname "Pactus GUI" \
  "${PACKAGE_NAME}-osx-64.dmg" \
  "${ROOT_DIR}/pactus-gui.app"

echo "Creating archive"
mkdir ${PACKAGE_DIR}
cp ${BUILD_DIR}/pactus-daemon     ${PACKAGE_DIR}
cp ${BUILD_DIR}/pactus-wallet     ${PACKAGE_DIR}
cp -R ${ROOT_DIR}/pactus-gui.app  ${PACKAGE_DIR}

tar -czvf ${ROOT_DIR}/${PACKAGE_NAME}-osx-64.tar.gz -p ${PACKAGE_NAME}
