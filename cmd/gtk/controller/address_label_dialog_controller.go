//go:build gtk

package controller

import (
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/view"
)

type AddressLabelDialogController struct {
	view *view.AddressLabelDialogView
}

func NewAddressLabelDialogController(view *view.AddressLabelDialogView) *AddressLabelDialogController {
	return &AddressLabelDialogController{view: view}
}

func (c *AddressLabelDialogController) Run(oldLabel string) (string, bool) {
	c.view.LabelEntry.SetText(oldLabel)

	newLabel := ""
	ok := false

	onOk := func() {
		newLabel = gtkutil.GetEntryText(c.view.LabelEntry)
		ok = true
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

	return newLabel, ok
}
