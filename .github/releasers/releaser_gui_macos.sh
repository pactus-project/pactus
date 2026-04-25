#!/bin/bash

set -e

ROOT_DIR="$(pwd)"
VERSION="$(echo `git -C ${ROOT_DIR} describe --abbrev=0 --tags` | sed 's/^.//')" # "v1.2.3" -> "1.2.3"
BUILD_DIR="${ROOT_DIR}/build"
PACKAGE_NAME="pactus-gui_${VERSION}"
PACKAGE_DIR="${ROOT_DIR}/${PACKAGE_NAME}"

# Ensure GTK prefix is provided for bundling assets
if [ -z "${LIB_HOME}" ]; then
    echo "LIB_HOME is not set. Set it to your GTK prefix (e.g. /opt/homebrew or /usr/local)."
    exit 1
fi

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

mkdir -p ${PACKAGE_DIR}

echo "Building the binaries for macOS ${ARC} architecture"

cd ${ROOT_DIR}
CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/pactus-daemon ./cmd/daemon
CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/pactus-wallet ./cmd/wallet
CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/pactus-shell ./cmd/shell
CGO_ENABLED=1 go build -ldflags "-s -w -extldflags -headerpad_max_install_names" -trimpath -tags gtk -o ${BUILD_DIR}/pactus-gui ./cmd/gtk


echo "Installing gtk-mac-bundler"
git clone https://gitlab.gnome.org/GNOME/gtk-mac-bundler.git
cd gtk-mac-bundler

# A workaround to make bundle without building GTK+ using jhbuild.
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
mkdir -p ${GUI_BUNDLE}

cp ${BUILD_DIR}/pactus-gui                ${GUI_BUNDLE}
cp ${ROOT_DIR}/.github/releasers/macos/*  ${GUI_BUNDLE}

# https://stackoverflow.com/questions/21242932/sed-i-may-not-be-used-with-stdin-on-mac-os-x
sed -i '' "s/%SHORTVERSION%/${VERSION}/"     ${GUI_BUNDLE}/Info.plist
sed -i '' "s/%VERSION%/Version ${VERSION}/"  ${GUI_BUNDLE}/Info.plist

export GUI_BUNDLE
export ROOT_DIR

${BUNDLER} ${GUI_BUNDLE}/gui.bundle

# Removing Cellar as workaround
rm -rf ${ROOT_DIR}/pactus-gui.app/Contents/Resources/Cellar


# After gtk-mac-bundler and your fix-install-names script...

if [ ! -z "${MACOS_CERT_IDENTITY}" ]; then
    echo "=== Step 1: Signing all .so and .dylib files inside the app bundle..."

    # Use find to locate all .so and .dylib files, regardless of permissions.
    find ${ROOT_DIR}/pactus-gui.app/Contents \( -name "*.dylib" -o -name "*.so" \) -print0 | while IFS= read -r -d '' binary; do
        echo "Signing binary: $binary"
        codesign --force --timestamp --options runtime --verbose --sign "${MACOS_CERT_IDENTITY}" "$binary"
    done

    echo "=== Step 2: Signing all executable binaries inside the app bundle..."

    # Use find to locate all regular files with executable bits, check for Mach-O, and sign them.
    find ${ROOT_DIR}/pactus-gui.app/Contents -type f -perm +111 -print0 | while IFS= read -r -d '' binary; do
        if file "$binary" | grep -q "Mach-O"; then
            echo "Signing executable: $binary"
            codesign --force --timestamp --options runtime --verbose --sign "${MACOS_CERT_IDENTITY}" "$binary"
        fi
    done

    echo "=== Step 3: Signing standalone binaries in build directory..."
    for bin in pactus-daemon pactus-wallet pactus-shell pactus-gui; do
        if [ -f "${BUILD_DIR}/${bin}" ]; then
            codesign --force --timestamp --options runtime --verbose --sign "${MACOS_CERT_IDENTITY}" "${BUILD_DIR}/${bin}"
        fi
    done

    echo "=== Step 4: Finally, signing the whole .app bundle..."
    codesign --force --timestamp --options runtime --verbose --sign "${MACOS_CERT_IDENTITY}" ${ROOT_DIR}/pactus-gui.app
fi

# if [ ! -z "${MACOS_CERT_IDENTITY}" ]; then
#     echo "=== Signing artifacts..."
#     codesign --force --options runtime --timestamp --sign "${MACOS_CERT_IDENTITY}" ${BUILD_DIR}/pactus-daemon
#     codesign --force --options runtime --timestamp --sign "${MACOS_CERT_IDENTITY}" ${BUILD_DIR}/pactus-wallet
#     codesign --force --options runtime --timestamp --sign "${MACOS_CERT_IDENTITY}" ${BUILD_DIR}/pactus-shell
#     codesign --force --options runtime --timestamp --sign "${MACOS_CERT_IDENTITY}" ${BUILD_DIR}/pactus-gui

#     echo "=== Signing app bundle..."
#     codesign --force --options runtime --timestamp --sign "${MACOS_CERT_IDENTITY}" ${ROOT_DIR}/pactus-gui.app
# fi

echo "Creating dmg"
# https://github.com/create-dmg/create-dmg
create-dmg --version
create-dmg --skip-jenkins \
  --volname "Pactus GUI" \
  "${FILE_NAME}.dmg" \
  "${ROOT_DIR}/pactus-gui.app"



if [ ! -z "${APPLE_ID}" ]; then
    echo "=== Submitting for notarization..."

    # Capture submission ID and check final status
    SUBMISSION_ID=$(xcrun notarytool submit "${FILE_NAME}.dmg" \
        --apple-id "${APPLE_ID}" \
        --password "${APPLE_PASSWORD}" \
        --team-id "${APPLE_TEAM_ID}" \
        --wait --output-format json | jq -r '.id')

    STATUS=$(xcrun notarytool info "$SUBMISSION_ID" \
        --apple-id "${APPLE_ID}" \
        --password "${APPLE_PASSWORD}" \
        --team-id "${APPLE_TEAM_ID}" \
        --output-format json | jq -r '.status')

    if [ "$STATUS" != "Accepted" ]; then
        echo "Notarization failed with status: $STATUS"
        xcrun notarytool log "$SUBMISSION_ID" \
            --apple-id "${APPLE_ID}" \
            --password "${APPLE_PASSWORD}" \
            --team-id "${APPLE_TEAM_ID}" \
            notarization.log
        cat notarization.log
        exit 1
    fi



    # xcrun notarytool submit "${FILE_NAME}.dmg" \
    #     --apple-id "${APPLE_ID}" \
    #     --password "${APPLE_PASSWORD}" \
    #     --team-id "${APPLE_TEAM_ID}" \
    #     --wait

    # echo "Stapling DMG only (the app inside gets the ticket automatically)..."
    # # ✅ FIX: Only staple the DMG – the .app was not notarized separately, so stapling it would cause error 65.
    # xcrun stapler staple "${FILE_NAME}.dmg"

    # # ❌ REMOVED: Stapling the standalone .app
    # # xcrun stapler staple "${ROOT_DIR}/pactus-gui.app"
fi

echo "Creating tar.gz archive"
cp ${BUILD_DIR}/pactus-daemon     ${PACKAGE_DIR}
cp ${BUILD_DIR}/pactus-wallet     ${PACKAGE_DIR}
cp ${BUILD_DIR}/pactus-shell      ${PACKAGE_DIR}
cp -R ${ROOT_DIR}/pactus-gui.app  ${PACKAGE_DIR}

tar -czvf ${ROOT_DIR}/${FILE_NAME}.tar.gz -p ${PACKAGE_NAME}
