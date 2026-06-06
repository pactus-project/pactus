//go:build gtk

package assets

import (
	_ "embed"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
)

var (
	//go:embed images/pactus.png
	imagePactusLogoData    []byte
	ImagePactusLogoTexture *gdk.Texture

	//go:embed images/gtk.png
	imageGTKLogoData    []byte
	ImageGTKLogoTexture *gdk.Texture
)

func initImages() {
	ImagePactusLogoTexture = TextureFromBytes(imagePactusLogoData)
	ImageGTKLogoTexture = TextureFromBytes(imageGTKLogoData)
}
