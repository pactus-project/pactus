//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

var (
	//go:embed assets/icons/add.svg
	iconAdd []byte

	//go:embed assets/icons/ok.svg
	iconOK []byte

	//go:embed assets/icons/cancel.svg
	iconCancel []byte

	//go:embed assets/icons/password.svg
	iconPassword []byte

	//go:embed assets/icons/seed.svg
	iconSeed []byte

	//go:embed assets/icons/close.svg
	iconClose []byte

	//go:embed assets/icons/send.svg
	iconSend []byte
)

func pixbufToIcon16(pixbuf *gdk.Pixbuf) *gtk.Image {
	if pixbuf == nil {
		// Return empty image if pixbuf failed to load
		image, _ := gtk.ImageNew()
		return image
	}
	resized, _ := pixbuf.ScaleSimple(16, 16, gdk.INTERP_NEAREST)
	image, _ := gtk.ImageNewFromPixbuf(resized)
	image.ShowAll()

	image.SetMarginEnd(2)

	return image
}

func AddIcon() *gtk.Image {
	pixbuf, _ := gdk.PixbufNewFromDataOnly(iconAdd)

	return pixbufToIcon16(pixbuf)
}

func OkIcon() *gtk.Image {
	pixbuf, _ := gdk.PixbufNewFromDataOnly(iconOK)

	return pixbufToIcon16(pixbuf)
}

func CancelIcon() *gtk.Image {
	pixbuf, _ := gdk.PixbufNewFromDataOnly(iconCancel)

	return pixbufToIcon16(pixbuf)
}

func CloseIcon() *gtk.Image {
	pixbuf, _ := gdk.PixbufNewFromDataOnly(iconClose)

	return pixbufToIcon16(pixbuf)
}

func PasswordIcon() *gtk.Image {
	pixbuf, _ := gdk.PixbufNewFromDataOnly(iconPassword)

	return pixbufToIcon16(pixbuf)
}

func SeedIcon() *gtk.Image {
	pixbuf, _ := gdk.PixbufNewFromDataOnly(iconSeed)

	return pixbufToIcon16(pixbuf)
}

func SendIcon() *gtk.Image {
	pixbuf, _ := gdk.PixbufNewFromDataOnly(iconSend)

	return pixbufToIcon16(pixbuf)
}
