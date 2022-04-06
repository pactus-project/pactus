#!/bin/bash


OS=`uname | tr [:upper:] [:lower:]`
MACH=`uname -m`
TAG=`git describe --abbrev=0 --tags`

make release

if [[ "$OS" == "mingw"* ]]; then
    OS="windows"
    7z a zarb-windows-$TAG-$MACH.zip zarb-daemon.exe
else
    tar -czvf zarb-$OS-$TAG-$MACH.tar.gz ./zarb-daemon
fi
