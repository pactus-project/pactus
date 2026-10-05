//go:build gtk

package assets

import (
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
)

func InitAssets() {
	initIcons()
	initImages()
}

// missingTexture creates a solid gray square texture used as a placeholder for a
// missing image. It returns nil if size is not positive.
func missingTexture(size int) *gdk.Texture {
	if size <= 0 {
		return nil
	}

	// A memory texture needs raw pixel data. Build an opaque gray square using
	// the R8G8B8 format (three bytes per pixel).
	pixels := make([]byte, size*size*3)
	for i := range pixels {
		pixels[i] = 0xee
	}

	texture := gdk.NewMemoryTexture(size, size, gdk.MemoryR8G8B8,
		glib.NewBytes(pixels), uint(size*3))

	return &texture.Texture
}

func TextureFromBytes(data []byte) *gdk.Texture {
	bytes := glib.NewBytes(data)
	texture, err := gdk.NewTextureFromBytes(bytes)
	if err != nil {
		return missingTexture(16)
	}

	return texture
}
