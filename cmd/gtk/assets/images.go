//go:build gtk

package assets

import (
	_ "embed"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
)

var (
	//go:embed images/pactus-tray-dark-16.png
	ImagePactusTrayDark16Data []byte

	//go:embed images/pactus-tray-light-16.png
	ImagePactusTrayLight16Data []byte

	//go:embed images/pactus-tray-dark-24.png
	ImagePactusTrayDark24Data []byte

	//go:embed images/pactus-tray-light-24.png
	ImagePactusTrayLight24Data []byte

	//go:embed images/pactus-tray-dark-32.png
	ImagePactusTrayDark32Data []byte

	//go:embed images/pactus-tray-light-32.png
	ImagePactusTrayLight32Data []byte

	//go:embed images/pactus.png
	ImagePactusLogoData    []byte
	ImagePactusLogoTexture *gdk.Texture

	//go:embed images/gtk.png
	ImageGTKLogoData    []byte
	ImageGTKLogoTexture *gdk.Texture
)

func initImages() {
	ImagePactusLogoTexture = TextureFromBytes(ImagePactusLogoData)
	ImageGTKLogoTexture = TextureFromBytes(ImageGTKLogoData)
}
