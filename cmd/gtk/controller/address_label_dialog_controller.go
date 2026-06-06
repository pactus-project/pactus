//go:build gtk

package controller

import (
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
)

type AddressLabelDialogController struct {
	view  *view.AddressLabelDialogView
	model *model.WalletModel
}

func NewAddressLabelDialogController(
	view *view.AddressLabelDialogView,
	model *model.WalletModel,
) *AddressLabelDialogController {
	return &AddressLabelDialogController{view: view, model: model}
}

func (c *AddressLabelDialogController) Show(address string, onUpdate func()) {
	oldLabel := c.model.AddressLabel(address)
	c.view.LabelEntry.SetText(oldLabel)

	onOk := func() {
		newLabel := gtkutil.EntryGetText(c.view.LabelEntry)
		if err := c.model.SetAddressLabel(address, newLabel); err != nil {
			gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

			return
		}
		c.view.Window.Close()
		onUpdate()
	}
	onCancel := func() {
		c.view.Window.Close()
	}

	gtkutil.ConnectButtonSignal(c.view.ButtonOK, onOk)
	gtkutil.ConnectButtonSignal(c.view.ButtonCancel, onCancel)

	gtkutil.ShowModalWindow(c.view.Window)
}
