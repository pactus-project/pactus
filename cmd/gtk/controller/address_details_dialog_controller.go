//go:build gtk

package controller

import (
	"errors"

	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/wallet/vault"
)

type AddressDetailsModel interface {
	AddressInfo(addr string) *vault.AddressInfo
}

type AddressDetailsDialogController struct {
	view *view.AddressDetailsDialogView
}

func NewAddressDetailsDialogController(view *view.AddressDetailsDialogView) *AddressDetailsDialogController {
	return &AddressDetailsDialogController{view: view}
}

func (c *AddressDetailsDialogController) Run(model AddressDetailsModel, addr string) error {
	info := model.AddressInfo(addr)
	if info == nil {
		gtkutil.ShowErrorDialog(nil, "address not found")

		return errors.New("address not found")
	}

	c.view.AddressEntry.SetText(info.Address)
	c.view.PubKeyEntry.SetText(info.PublicKey)
	c.view.PathEntry.SetText(info.Path)

	c.view.ConnectSignals(map[string]any{
		"on_close": func() { c.view.Dialog.Close() },
	})

	c.view.Dialog.SetModal(true)
	gtkutil.RunDialog(c.view.Dialog)

	return nil
}
