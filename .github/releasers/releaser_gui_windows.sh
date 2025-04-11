#!/bin/bash

set -e

ROOT_DIR="$(pwd)"
VERSION="$(echo `git -C ${ROOT_DIR} describe --abbrev=0 --tags` | sed 's/^.//')" # "v1.2.3" -> "1.2.3"
BUILD_DIR="${ROOT_DIR}/build"
PACKAGE_NAME="pactus-gui_${VERSION}"
PACKAGE_DIR="${ROOT_DIR}/${PACKAGE_NAME}"

mkdir ${PACKAGE_DIR}

echo "Building the binaries"

# This fixes a bug in pkgconfig: invalid flag in pkg-config --libs: -Wl,-luuid
sed -i -e 's/-Wl,-luuid/-luuid/g' /mingw64/lib/pkgconfig/gdk-3.0.pc

CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/pactus-daemon.exe ./cmd/daemon
CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/pactus-wallet.exe ./cmd/wallet
CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/pactus-shell.exe  ./cmd/shell
go build -ldflags "-s -w -H windowsgui" -trimpath -tags gtk -o ${BUILD_DIR}/pactus-gui.exe ./cmd/gtk

# Copying the neccesary libraries
echo "Creating GUI directory"
GUI_DIR="${PACKAGE_DIR}/pactus-gui"
mkdir ${GUI_DIR}

echo "Copying required DLLs and EXEs from ${MINGW_PREFIX}/bin"

cd "${MINGW_PREFIX}/bin"
cp -p \
  "gdbus.exe" \
  "gspawn-win64-helper.exe" \
  "gspawn-win64-helper-console.exe" \
  "libatk-1.0-0.dll" \
  "libbrotlicommon.dll" \
  "libbrotlidec.dll" \
  "libbz2-1.dll" \
  "libcairo-2.dll" \
  "libcairo-gobject-2.dll" \
  "libdatrie-1.dll" \
  "libdeflate.dll" \
  "libepoxy-0.dll" \
  "libexpat-1.dll" \
  "libffi-8.dll" \
  "libfontconfig-1.dll" \
  "libfreetype-6.dll" \
  "libfribidi-0.dll" \
  "libgcc_s_seh-1.dll" \
  "libgdk_pixbuf-2.0-0.dll" \
  "libgdk-3-0.dll" \
  "libgio-2.0-0.dll" \
  "libglib-2.0-0.dll" \
  "libgmodule-2.0-0.dll" \
  "libgobject-2.0-0.dll" \
  "libgomp-1.dll" \
  "libgraphite2.dll" \
  "libgtk-3-0.dll" \
  "libharfbuzz-0.dll" \
  "libiconv-2.dll" \
  "libintl-8.dll" \
  "libjbig-0.dll" \
  "libjpeg-8.dll" \
  "libLerc.dll" \
  "liblzma-5.dll" \
  "libpango-1.0-0.dll" \
  "libpangocairo-1.0-0.dll" \
  "libpangoft2-1.0-0.dll" \
  "libpangowin32-1.0-0.dll" \
  "libpcre2-16-0.dll" \
  "libpcre2-32-0.dll" \
  "libpcre2-8-0.dll" \
  "libpcre2-posix-3.dll" \
  "libpixman-1-0.dll" \
  "libpng16-16.dll" \
  "librsvg-2-2.dll" \
  "libstdc++-6.dll" \
  "libsystre-0.dll" \
  "libthai-0.dll" \
  "libtiff-6.dll" \
  "libtre-5.dll" \
  "libwebp-7.dll" \
  "libwinpthread-1.dll" \
  "libxml2-2.dll" \
  "libzstd.dll" \
  "zlib1.dll" \
  "libsharpyuv-0.dll" \
  "${GUI_DIR}"


echo "Copying GDK pixbuf from ${MINGW_PREFIX}/lib/gdk-pixbuf-2.0"
mkdir -p "${GUI_DIR}/lib/gdk-pixbuf-2.0"
cp -rp "${MINGW_PREFIX}/lib/gdk-pixbuf-2.0" "${GUI_DIR}/lib"

#########
### Based on this toturial: https://www.gtk.org/docs/installations/windows#building-and-distributing-your-application

echo "Step 1: gtk-3.0/"
mkdir -p "${GUI_DIR}/share/themes/Windows10/gtk-3.0/"
cp -rp "${MINGW_PREFIX}/share/gtk-3.0" "${GUI_DIR}/share/themes/Windows10/"

echo "Step 2: Adwaita icons"
mkdir -p "${GUI_DIR}/share/icons/Adwaita"
cp -rp "${MINGW_PREFIX}/share/icons/Adwaita" "${GUI_DIR}/share/icons"

echo "Step 3: hicolor icons"
mkdir -p "${GUI_DIR}/share/icons/hicolor"
cp -rp "${MINGW_PREFIX}/share/icons/hicolor" "${GUI_DIR}/share/icons"

echo "Step 4: settings.ini"
mkdir "${GUI_DIR}/share/gtk-3.0/"
echo "[Settings]
gtk-theme-name=Windows10" > "${GUI_DIR}/share/gtk-3.0/settings.ini"

echo "Step 5: glib-compile-schemas"
mkdir -p "${GUI_DIR}/share/glib-2.0/schemas"
cp -rp "${MINGW_PREFIX}/share/glib-2.0/schemas/gschemas.compiled" "${GUI_DIR}/share/glib-2.0/schemas"


# Moving binaries to package directory
cd ${ROOT_DIR}
echo "Moving binaries"
mv ${BUILD_DIR}/pactus-daemon.exe  ${PACKAGE_DIR}/pactus-daemon.exe
mv ${BUILD_DIR}/pactus-wallet.exe  ${PACKAGE_DIR}/pactus-wallet.exe
mv ${BUILD_DIR}/pactus-shell.exe   ${PACKAGE_DIR}/pactus-shell.exe
mv ${BUILD_DIR}/pactus-gui.exe     ${PACKAGE_DIR}/pactus-gui/pactus-gui.exe

echo "Archiving the package"
7z a ${ROOT_DIR}/${PACKAGE_NAME}_windows_amd64.zip ${PACKAGE_DIR}

echo "Creating installer"
echo "
[Setup]
Appname=Pactus
AppVersion=${VERSION}
DefaultDirName={autopf64}/Pactus
DefaultGroupName=Pactus
[Files]
Source:${PACKAGE_NAME}/*; DestDir:{app}; Flags: recursesubdirs
[Icons]
Name:{group}\\Pactus GUI; Filename:{app}\\pactus-gui\\pactus-gui.exe;" >> ${ROOT_DIR}/inno.iss

cd ${ROOT_DIR}
INNO_DIR=$(cygpath -w -s "/c/Program Files (x86)/Inno Setup 6")
${INNO_DIR}/ISCC.exe ${ROOT_DIR}/inno.iss
mv Output/mysetup.exe ${ROOT_DIR}/${PACKAGE_NAME}_windows_amd64_installer.exe
