//go:build gtk

package assets

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func InitAssets() {
	initIcons()
	initImages()
}

// missingPixbuf tries to return an icon-theme "missing image" pixbuf and falls
// back to a simple solid square if the theme isn't available.
func missingPixbuf(size int) *gdk.Pixbuf {
	theme, err := gtk.IconThemeGetDefault()
	if err == nil && theme != nil {
		pixbuf, err := theme.LoadIcon("image-missing", size, 0)
		if err == nil || pixbuf != nil {
			return pixbuf
		}
	}

	// Last resort: a tiny gray square (ARGB32).
	pixbuf, err := gdk.PixbufNew(gdk.COLORSPACE_RGB, true, 8, size, size)
	if err == nil && pixbuf != nil {
		// 0xAARRGGBB
		pixbuf.Fill(0xFF666666)

		return pixbuf
	}

	return nil
}
