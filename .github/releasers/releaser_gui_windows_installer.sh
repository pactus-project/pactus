#!/bin/bash

set -e

ROOT_DIR="$(pwd)"
VERSION="$(echo `git -C ${ROOT_DIR} describe --abbrev=0 --tags` | sed 's/^.//')"
BUILD_DIR="${ROOT_DIR}/build"
PACKAGE_NAME="pactus-gui_${VERSION}"
PACKAGE_DIR="${ROOT_DIR}/${PACKAGE_NAME}"
FILE_NAME="${PACKAGE_NAME}_windows_amd64"
INNO_PATH="/c/Program Files (x86)/Inno Setup 6"

echo "🚀 Starting Pactus GUI Windows packaging..."

# Create package directory
mkdir -p "${PACKAGE_DIR}/pactus-gui"

# Bundle GTK application using Python bundler
python3 "${ROOT_DIR}/.github/releasers/gtk-win-bundler.py" \
    "${BUILD_DIR}/signed/pactus-gui.exe" \
    "${PACKAGE_DIR}/pactus-gui"

# Move other binaries
cp ${BUILD_DIR}/signed/pactus-daemon.exe  ${PACKAGE_DIR}/pactus-daemon.exe
cp ${BUILD_DIR}/signed/pactus-wallet.exe  ${PACKAGE_DIR}/pactus-wallet.exe
cp ${BUILD_DIR}/signed/pactus-shell.exe   ${PACKAGE_DIR}/pactus-shell.exe
cp ${BUILD_DIR}/signed/pactus-gui.exe     ${PACKAGE_DIR}/pactus-gui/pactus-gui.exe

# Create archive
7z a ${ROOT_DIR}/${FILE_NAME}.zip ${PACKAGE_DIR}

# Create installer
cat << EOF > ${ROOT_DIR}/inno.iss
[Setup]
AppId=Pactus
AppName=Pactus
AppVersion=${VERSION}
AppPublisher=Pactus
AppPublisherURL=https://pactus.org/
DefaultDirName={autopf}/Pactus
DefaultGroupName=Pactus
SetupIconFile=.github/releasers/pactus.ico
LicenseFile=LICENSE
Uninstallable=yes
UninstallDisplayIcon={app}\\pactus-gui\\pactus-gui.exe

[Files]
Source:"${PACKAGE_NAME}/*"; DestDir:"{app}"; Flags: recursesubdirs

[Icons]
Name:"{group}\\Pactus"; Filename:"{app}\\pactus-gui\\pactus-gui.exe"
Name:"{commondesktop}\\Pactus"; Filename:"{app}\\pactus-gui\\pactus-gui.exe"

[Run]
Filename:"{app}\\pactus-gui\\pactus-gui.exe"; Description:"Launch Pactus"; Flags: postinstall nowait
EOF

# Build installer
INNO_DIR=$(cygpath -w -s "${INNO_PATH}")
"${INNO_DIR}/ISCC.exe" "${ROOT_DIR}/inno.iss"
mv "Output/mysetup.exe" "${BUILD_DIR}/unsigned/${FILE_NAME}_installer.exe"

echo "🎉 Build complete! Package: ${BUILD_DIR}/unsigned/${FILE_NAME}_installer.exe"
