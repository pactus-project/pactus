//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

var (
	//go:embed assets/images/add.svg
	imgAddIcon []byte
)

func addIcon() *gtk.Image {
	pixbuf, _ := gdk.PixbufNewFromDataOnly(imgAddIcon)
	resized, _ := pixbuf.ScaleSimple(12, 12, gdk.INTERP_NEAREST)
	image, _ := gtk.ImageNewFromPixbuf(resized)
	image.ShowAll()

	return image
}
