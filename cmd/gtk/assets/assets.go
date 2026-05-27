//go111:build gtk

package assets

import (
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
)

func InitAssets() {
	initIcons()
	initImages()
}

// missingTexture allocates raw bytes to create a solid color texture.
func missingTexture(size int) *gdk.Texture {
	// 4 bytes per pixel (RGBA)
	stride := size * 4
	byteCount := stride * size
	pixels := make([]byte, byteCount)

	// Fill the byte array with 0xFF666666 (R:0x66, G:0x66, B:0x66, A:0xFF)
	for i := 0; i < byteCount; i += 4 {
		pixels[i] = 0x66   // Red
		pixels[i+1] = 0x66 // Green
		pixels[i+2] = 0x66 // Blue
		pixels[i+3] = 0xFF // Alpha
	}

	// Wrap bytes into a GLib Bytes object
	gBytes := glib.NewBytes(pixels)

	// Create the GTK4 memory texture
	texture := gdk.NewMemoryTexture(
		size,
		size,
		gdk.MemoryR8G8B8A8,
		gBytes,
		uint(stride),
	)

	return &texture.Texture
}
