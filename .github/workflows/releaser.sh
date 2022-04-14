#!/bin/bash


OS=`uname | tr [:upper:] [:lower:]`
MACH=`uname -m`
TAG=`git describe --abbrev=0 --tags`

make release

## Releasing GUI app
if [[ "$OS" == "Linux" ]]; then
    sudo apt install libgtk-3-dev libcairo2-dev libglib2.0-dev
    go build -ldflags "-s -w" -tags glib_2_66,gtk -o ./build/zarb-gui ./cmd/gtk
elif [[ "$OS" == "Darwin" ]]; then
    brew install gobject-introspection gtk+3 adwaita-icon-theme
    go build -ldflags "-s -w" -tags gtk -o ./build/zarb-gui ./cmd/gtk
elif [[ "$OS" == "mingw"* ]]; then
    pacman -Sy --noconfirm git \
        mingw-w64-x86_64-gtk3 \
        mingw-w64-x86_64-toolchain \
        mingw-w64-x86_64-go \
        mingw-w64-x86_64-pkg-config \
        mingw-w64-x86_64-gcc \
        base-devel \
        glib2-devel

    go build -ldflags "-s -w" -tags glib_2_66,gtk -o ./build/zarb-gui.exe ./cmd/gtk
fi

if [[ "$OS" == "mingw"* ]]; then
    OS="windows"
    7z a zarb-windows-$TAG-$MACH.zip zarb-daemon.exe zarb-wallet.exe zarb-gui.exe
else
    tar -czvf zarb-$OS-$TAG-$MACH.tar.gz ./build/zarb-daemon ./build/zarb-wallet ./build/zarb-gui
fi
