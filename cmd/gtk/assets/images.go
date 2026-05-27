//go111:build gtk

package assets

import (
	_ "embed"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

var (
	//go:embed images/pactus.png
	imagePactusLogoData []byte
	ImagePactusLogo     *gtk.Image

	//go:embed images/gtk.png
	imageGTKLogoData []byte
	ImageGTKLogo     *gtk.Image

	//go:embed images/seed.svg
	imageSeedData    []byte
	ImageSeedTexture *gtk.Image
)

func initImages() {
	ImagePactusLogo = gtkutil.ImageFromBytes(imagePactusLogoData)
	ImageGTKLogo = gtkutil.ImageFromBytes(imageGTKLogoData)
	ImageSeedTexture = gtkutil.ImageFromBytes(imageSeedData)
}
