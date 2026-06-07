//go:build gtk

package assets

import (
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gdkpixbuf/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
)

func InitAssets() {
	initIcons()
	initImages()
}

// missingTexture creates a Picture widget that displays a solid gray square
// as a placeholder for a missing image. It returns nil if the pixbuf cannot be created.
func missingTexture(size int) *gdk.Texture {
	// Create a gray square pixbuf
	pixbuf := gdkpixbuf.NewPixbuf(gdkpixbuf.ColorspaceRGB, true, 8, size, size)
	if pixbuf == nil {
		return nil
	}
	pixbuf.Fill(0xeeeeee)

	return gdk.NewTextureForPixbuf(pixbuf)
}

func TextureFromBytes(data []byte) *gdk.Texture {
	bytes := glib.NewBytes(data)
	texture, err := gdk.NewTextureFromBytes(bytes)
	if err != nil {
		return missingTexture(16)
	}

	return texture
}
