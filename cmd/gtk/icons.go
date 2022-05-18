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

	//go:embed assets/images/ok.svg
	imgOkIcon []byte

	//go:embed assets/images/cancel.svg
	imgCancelIcon []byte
)

func pixbufToIcon16(pixbuf *gdk.Pixbuf) *gtk.Image {
	resized, _ := pixbuf.ScaleSimple(16, 16, gdk.INTERP_NEAREST)
	image, _ := gtk.ImageNewFromPixbuf(resized)
	image.ShowAll()

	image.SetMarginEnd(2)

	return image
}

func AddIcon() *gtk.Image {
	pixbuf, _ := gdk.PixbufNewFromDataOnly(imgAddIcon)
	return pixbufToIcon16(pixbuf)
}

func OkIcon() *gtk.Image {
	pixbuf, _ := gdk.PixbufNewFromDataOnly(imgOkIcon)
	return pixbufToIcon16(pixbuf)
}

func CancelIcon() *gtk.Image {
	pixbuf, _ := gdk.PixbufNewFromDataOnly(imgCancelIcon)
	return pixbufToIcon16(pixbuf)
}
