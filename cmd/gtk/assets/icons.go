//go:build gtk

package assets

import (
	_ "embed"

	"github.com/gotk3/gotk3/gdk"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

var (
	//go:embed icons/add.svg
	iconAddData     []byte
	IconAddPixbuf16 *gdk.Pixbuf

	//go:embed icons/ok.svg
	iconOKData     []byte
	IconOkPixbuf16 *gdk.Pixbuf

	//go:embed icons/cancel.svg
	iconCancelData     []byte
	IconCancelPixbuf16 *gdk.Pixbuf

	//go:embed icons/password.svg
	iconPasswordData     []byte
	IconPasswordPixbuf16 *gdk.Pixbuf

	//go:embed icons/seed.svg
	iconSeedData     []byte
	IconSeedPixbuf16 *gdk.Pixbuf

	//go:embed icons/close.svg
	iconCloseData     []byte
	IconClosePixbuf16 *gdk.Pixbuf

	//go:embed icons/send.svg
	iconSendData     []byte
	IconSendPixbuf16 *gdk.Pixbuf

	//go:embed icons/fee.svg
	iconFeeData     []byte
	IconFeePixbuf16 *gdk.Pixbuf

	//go:embed icons/refresh.svg
	iconRefreshData     []byte
	IconRefreshPixbuf16 *gdk.Pixbuf

	//go:embed icons/prev.svg
	iconPrevData     []byte
	IconPrevPixbuf16 *gdk.Pixbuf

	//go:embed icons/next.svg
	iconNextData     []byte
	IconNextPixbuf16 *gdk.Pixbuf
)

func initIcons() {
	toPixbuf := func(data []byte) *gdk.Pixbuf {
		pixbuf, err := gtkutil.PixbufFromBytes(data, gtkutil.WithSize(16, 16))
		if err != nil {
			return missingPixbuf(16)
		}

		return pixbuf
	}

	IconAddPixbuf16 = toPixbuf(iconAddData)
	IconOkPixbuf16 = toPixbuf(iconOKData)
	IconCancelPixbuf16 = toPixbuf(iconCancelData)
	IconPasswordPixbuf16 = toPixbuf(iconPasswordData)
	IconSeedPixbuf16 = toPixbuf(iconSeedData)
	IconClosePixbuf16 = toPixbuf(iconCloseData)
	IconSendPixbuf16 = toPixbuf(iconSendData)
	IconFeePixbuf16 = toPixbuf(iconFeeData)
	IconRefreshPixbuf16 = toPixbuf(iconRefreshData)
	IconPrevPixbuf16 = toPixbuf(iconPrevData)
	IconNextPixbuf16 = toPixbuf(iconNextData)
}
