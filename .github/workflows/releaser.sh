#!/bin/sh


OS=`uname | tr [:upper:] [:lower:]`
MACH=`uname -m`
TAG=`git describe --abbrev=0 --tags`

make bls && make build_bls_release

if [[ "$OS" == "mingw"* ]]; then
    OS="windows"
    Compress-Archive ./zarb.exe zarb-windows-$TAG-$MACH.zip
else
    tar -czvf zarb-$OS-$TAG-$MACH.tar.gz ./zarb
fi