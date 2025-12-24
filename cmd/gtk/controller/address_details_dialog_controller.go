//go:build gtk

package controller

import (
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/wallet/types"
)

type AddressDetailsModel interface {
	AddressInfo(addr string) *types.AddressInfo
}

type AddressDetailsDialogController struct {
	view  *view.AddressDetailsDialogView
	model AddressDetailsModel
}

func NewAddressDetailsDialogController(
	view *view.AddressDetailsDialogView,
	model AddressDetailsModel,
) *AddressDetailsDialogController {
	return &AddressDetailsDialogController{view: view, model: model}
}

func (c *AddressDetailsDialogController) Run(addr string) {
	info := c.model.AddressInfo(addr)
	if info == nil {
		gtkutil.ShowErrorDialog(nil, "address not found")

		return
	}

	c.view.AddressEntry.SetText(info.Address)
	c.view.PubKeyEntry.SetText(info.PublicKey)
	c.view.PathEntry.SetText(info.Path)

	onClose := func() { c.view.Dialog.Close() }

	c.view.ConnectSignals(map[string]any{
		"on_close": onClose,
	})

	c.view.Dialog.SetModal(true)
	gtkutil.RunDialog(c.view.Dialog)
}
