//go:build gtk

package controller

import (
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type AddressDetailsDialogController struct {
	view  *view.AddressDetailsDialogView
	model *model.WalletModel
}

func NewAddressDetailsDialogController(
	view *view.AddressDetailsDialogView,
	model *model.WalletModel,
) *AddressDetailsDialogController {
	return &AddressDetailsDialogController{view: view, model: model}
}

func (c *AddressDetailsDialogController) Show(info *pactus.AddressInfo) {
	c.view.AddressEntry.SetText(info.Address)
	c.view.PubKeyEntry.SetText(info.PublicKey)
	c.view.PathEntry.SetText(info.Path)

	onClose := func() { c.view.Window.Close() }

	gtkutil.ConnectButtonSignal(c.view.ButtonClose, onClose)

	gtkutil.ShowModalWindow(c.view.Window)
}
