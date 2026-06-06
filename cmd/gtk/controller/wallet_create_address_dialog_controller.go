//go:build gtk

package controller

import (
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/crypto"
)

type WalletCreateAddressDialogController struct {
	view  *view.WalletCreateAddressDialogView
	model *model.WalletModel

	lsTypes []crypto.AddressType
}

func NewWalletCreateAddressDialogController(
	view *view.WalletCreateAddressDialogView,
	model *model.WalletModel,
) *WalletCreateAddressDialogController {
	return &WalletCreateAddressDialogController{view: view, model: model}
}

func (c *WalletCreateAddressDialogController) Show(onUpdate func()) {
	gtkutil.DropDownSetup(c.view.AddressTypeDrop, []string{
		"Validator",
		"BLS Account",
		"Ed25519 Account",
		"Secp256k1 Account",
	})

	c.lsTypes = []crypto.AddressType{
		crypto.AddressTypeValidator,
		crypto.AddressTypeBLSAccount,
		crypto.AddressTypeEd25519Account,
		crypto.AddressTypeSecp256k1Account,
	}

	c.view.AddressTypeDrop.SetSelected(2) // Edd25519 ia active.

	onOK := func() {
		label := gtkutil.EntryGetText(c.view.LabelEntry)
		typ := gtkutil.DropDownGetSelectedItem(c.view.AddressTypeDrop, c.lsTypes)

		if typ == crypto.AddressTypeEd25519Account ||
			typ == crypto.AddressTypeSecp256k1Account {
			PasswordProvider(c.model, func(password string, ok bool) {
				if !ok {
					return
				}

				_, err := c.model.NewAddress(typ, label, password)
				if err != nil {
					gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

					return
				}
				c.view.Window.Close()
				onUpdate()
			})
		} else {
			_, err := c.model.NewAddress(typ, label, "")
			if err != nil {
				gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

				return
			}
			c.view.Window.Close()
			onUpdate()
		}
	}

	onCancel := func() {
		c.view.Window.Close()
	}

	gtkutil.ConnectButtonSignal(c.view.ButtonOK, onOK)
	gtkutil.ConnectButtonSignal(c.view.ButtonCancel, onCancel)

	gtkutil.ShowModalWindow(c.view.Window)
}
