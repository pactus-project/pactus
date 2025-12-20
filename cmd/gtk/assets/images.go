//go:build gtk

package assets

import (
	_ "embed"

	"github.com/gotk3/gotk3/gdk"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

var (
	//go:embed images/pactus.png
	imagePactusLogoData   []byte
	ImagePactusLogoPixbuf *gdk.Pixbuf

	//go:embed images/gtk.png
	imageGTKLogoData   []byte
	ImageGTKLogoPixbuf *gdk.Pixbuf

	//go:embed images/seed.svg
	imageSeedData   []byte
	ImageSeedPixbuf *gdk.Pixbuf
)

func initImages() {
	toPixbuf := func(data []byte) *gdk.Pixbuf {
		pixbuf, err := gtkutil.PixbufFromBytes(data)
		if err != nil {
			return missingPixbuf(128)
		}

		return pixbuf
	}

	ImagePactusLogoPixbuf = toPixbuf(imagePactusLogoData)
	ImageGTKLogoPixbuf = toPixbuf(imageGTKLogoData)
	ImageSeedPixbuf = toPixbuf(imageSeedData)
}
