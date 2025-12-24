//go:build gtk

package controller

import (
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/view"
)

type AddressLabelModel interface {
	AddressLabel(address string) string
	SetAddressLabel(address string, label string) error
}

type AddressLabelDialogController struct {
	view  *view.AddressLabelDialogView
	model AddressLabelModel
}

func NewAddressLabelDialogController(
	view *view.AddressLabelDialogView,
	model AddressLabelModel,
) *AddressLabelDialogController {
	return &AddressLabelDialogController{view: view, model: model}
}

func (c *AddressLabelDialogController) Run(address string) {
	oldLabel := c.model.AddressLabel(address)
	c.view.LabelEntry.SetText(oldLabel)

	onOk := func() {
		newLabel := gtkutil.GetEntryText(c.view.LabelEntry)
		if err := c.model.SetAddressLabel(address, newLabel); err != nil {
			gtkutil.ShowError(err)

			return
		}
		c.view.Dialog.Close()
	}
	onCancel := func() {
		c.view.Dialog.Close()
	}

	c.view.ConnectSignals(map[string]any{
		"on_ok":     onOk,
		"on_cancel": onCancel,
	})

	c.view.Dialog.SetModal(true)
	gtkutil.RunDialog(c.view.Dialog)
}
