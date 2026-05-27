//go111:build gtk

package assets

import (
	_ "embed"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

var (
	//go:embed icons/add.svg
	iconAddData []byte
	IconAdd16   *gtk.Image

	//go:embed icons/ok.svg
	iconOKData []byte
	IconOk16   *gtk.Image

	//go:embed icons/cancel.svg
	iconCancelData []byte
	IconCancel16   *gtk.Image

	//go:embed icons/password.svg
	iconPasswordData []byte
	IconPassword16   *gtk.Image

	//go:embed icons/seed.svg
	iconSeedData []byte
	IconSeed16   *gtk.Image

	//go:embed icons/close.svg
	iconCloseData []byte
	IconClose16   *gtk.Image

	//go:embed icons/send.svg
	iconSendData []byte
	IconSend16   *gtk.Image

	//go:embed icons/fee.svg
	iconFeeData []byte
	IconFee16   *gtk.Image

	//go:embed icons/refresh.svg
	iconRefreshData []byte
	IconRefresh16   *gtk.Image

	//go:embed icons/prev.svg
	iconPrevData []byte
	IconPrev16   *gtk.Image

	//go:embed icons/next.svg
	iconNextData []byte
	IconNext16   *gtk.Image

	//go:embed icons/save.svg
	iconSaveData []byte
	IconSave16   *gtk.Image
)

func initIcons() {
	toImage := func(data []byte) *gtk.Image {
		return gtkutil.ImageFromBytes(data, gtkutil.WithImageSize(16, 16))
	}

	IconAdd16 = toImage(iconAddData)
	IconOk16 = toImage(iconOKData)
	IconCancel16 = toImage(iconCancelData)
	IconPassword16 = toImage(iconPasswordData)
	IconSeed16 = toImage(iconSeedData)
	IconClose16 = toImage(iconCloseData)
	IconSend16 = toImage(iconSendData)
	IconFee16 = toImage(iconFeeData)
	IconRefresh16 = toImage(iconRefreshData)
	IconPrev16 = toImage(iconPrevData)
	IconNext16 = toImage(iconNextData)
	IconSave16 = toImage(iconSaveData)
}
