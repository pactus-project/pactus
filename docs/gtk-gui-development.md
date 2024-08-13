# Pactus GUI

This document is quick guide for developing and updating [the Pactus Core GUI](../cmd/gtk/).

## Requirements

The Pactus Core GUI utilizes gtk for desktop GUI. To develop, build and test it you must have these packages installed:

### Linux

1. `libgtk-3-dev`
2. `libcairo2-dev`
3. `libglib2.0-dev`

Install using apt:

```bash
apt install libgtk-3-dev libcairo2-dev libglib2.0-dev
```

### Mac OS

1. `gtk+3`

Install using brew:

```bash
brew install gtk+3
```

### Windows

1. `glib2-devel`
2. `mingw-w64-x86_64-go`
3. `mingw-w64-x86_64-gtk3`
4. `mingw-w64-x86_64-glib2`
5. `mingw-w64-x86_64-gcc`
6. `mingw-w64-x86_64-pkg-config`


With these packages installed you can build GUI using `make build_gui` command. You can run the GUI like: `./pactus-gui`, `./pactus-gui.exe`.


The [Assets](../cmd/gtk/assets/) file includes required images, icons, ui files and custom CSS files. All [`.ui`](../cmd/gtk/assets/ui/) files are used to defined the user interface of GUI, for a proper edit and change on them make sure you have [Glade](https://glade.gnome.org/) installed on your machine.

## Running linter

When you make changes on GUI files and try to run linter using `make check`, it won't include [gtk](../cmd/gtk/) in it's checks. So make sure you add gtk build flag like this:

```bash
BUILD_TAG=gtk make check
```
