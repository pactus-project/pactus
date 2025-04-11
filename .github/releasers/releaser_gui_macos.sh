#!/bin/bash

set -e

ROOT_DIR="$(pwd)"
VERSION="$(echo `git -C ${ROOT_DIR} describe --abbrev=0 --tags` | sed 's/^.//')" # "v1.2.3" -> "1.2.3"
BUILD_DIR="${ROOT_DIR}/build"
PACKAGE_NAME="pactus-gui_${VERSION}"
PACKAGE_DIR="${ROOT_DIR}/${PACKAGE_NAME}"

# Check the architecture
ARC="$(uname -m)"

if [ "${ARC}" = "x86_64" ]; then
    FILE_NAME="${PACKAGE_NAME}_darwin_amd64"
elif [ "${ARC}" = "arm64" ]; then
    FILE_NAME="${PACKAGE_NAME}_darwin_arm64"
else
    echo "Unsupported architecture: ${ARC}"
    exit 1
fi

mkdir ${PACKAGE_DIR}

echo "Building the binaries"

cd ${ROOT_DIR}
CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/pactus-daemon ./cmd/daemon
CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/pactus-wallet ./cmd/wallet
CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/pactus-shell ./cmd/shell
go build -ldflags "-s -w" -trimpath -tags gtk -o ${BUILD_DIR}/pactus-gui ./cmd/gtk


echo "Installing gtk-mac-bundler"
git clone https://gitlab.gnome.org/GNOME/gtk-mac-bundler.git
cd gtk-mac-bundler

# A workaround to make bundle without builing GTK+ using jhbuild.
rm bundler/run-install-name-tool-change.sh
cp ${ROOT_DIR}/.github/releasers/macos/run-install-name-tool-change.sh bundler/run-install-name-tool-change.sh
chmod +x bundler/run-install-name-tool-change.sh

# make sure launcher is executable
chmod +x ${ROOT_DIR}/.github/releasers/macos/gtk3-launcher.sh

make install

export PATH=${PATH}:${HOME}/.bin:${HOME}/local/bin
BUNDLER=$(which gtk-mac-bundler)

echo "gtk-mac-bundler found at ${BUNDLER}"

cd -

echo "Bundling the GUI package"
GUI_BUNDLE=${ROOT_DIR}/gui-bundle
mkdir ${GUI_BUNDLE}

cp ${BUILD_DIR}/pactus-gui                ${GUI_BUNDLE}
cp ${ROOT_DIR}/.github/releasers/macos/*  ${GUI_BUNDLE}

# Icon
cp ${ROOT_DIR}/.github/releasers/pactus.icns  ${GUI_BUNDLE}

# https://stackoverflow.com/questions/21242932/sed-i-may-not-be-used-with-stdin-on-mac-os-x
sed -i '' "s/%SHORTVERSION%/${VERSION}/"     ${GUI_BUNDLE}/Info.plist
sed -i '' "s/%VERSION%/Version ${VERSION}/"  ${GUI_BUNDLE}/Info.plist

export GUI_BUNDLE
export ROOT_DIR

${BUNDLER} ${GUI_BUNDLE}/gui.bundle

# Removing Cellar as workaround
rm -rf ${ROOT_DIR}/pactus-gui.app/Contents/Resources/Cellar

echo "Creating dmg"
# https://github.com/create-dmg/create-dmg
create-dmg \
  --volname "Pactus GUI" \
  "${FILE_NAME}.dmg" \
  "${ROOT_DIR}/pactus-gui.app"

echo "Creating archive"
cp ${BUILD_DIR}/pactus-daemon     ${PACKAGE_DIR}
cp ${BUILD_DIR}/pactus-wallet     ${PACKAGE_DIR}
cp ${BUILD_DIR}/pactus-shell      ${PACKAGE_DIR}
cp -R ${ROOT_DIR}/pactus-gui.app  ${PACKAGE_DIR}

tar -czvf ${ROOT_DIR}/${FILE_NAME}.tar.gz -p ${PACKAGE_NAME}
