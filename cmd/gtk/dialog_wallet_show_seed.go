//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

var (
	//go:embed assets/ui/dialog_wallet_show_seed.ui
	uiWalletShowSeedDialog []byte

	//go:embed assets/images/seed.svg
	imageSeed []byte
)

func showSeed(seed string) {
	builder, err := gtk.BuilderNewFromString(string(uiWalletShowSeedDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_wallet_show_seed")
	textViewSeed := getTextViewObj(builder, "id_textview_seed")
	setTextViewContent(textViewSeed, seed)

	getButtonObj(builder, "id_button_close").SetImage(CloseIcon())

	pixbuf, _ := gdk.PixbufNewFromDataOnly(imageSeed)
	getImageObj(builder, "id_image_seed").SetFromPixbuf(pixbuf)

	onClose := func() {
		dlg.Close()
	}

	signals := map[string]any{
		"on_close": onClose,
	}
	builder.ConnectSignals(signals)

	dlg.SetModal(true)

	runDialog(dlg)
}
