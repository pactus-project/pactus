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
}

func NewWalletCreateAddressDialogController(
	view *view.WalletCreateAddressDialogView,
	model *model.WalletModel,
) *WalletCreateAddressDialogController {
	return &WalletCreateAddressDialogController{view: view, model: model}
}

func (c *WalletCreateAddressDialogController) Run() {
	combo := c.view.AddressTypeCombo
	combo.Append(crypto.AddressTypeValidator.String(), "Validator")
	combo.Append(crypto.AddressTypeBLSAccount.String(), "BLS Account")
	combo.Append(crypto.AddressTypeEd25519Account.String(), "Ed25519 Account")
	combo.Append(crypto.AddressTypeSecp256k1Account.String(), "Secp256k1 Account")

	combo.SetActive(2) // Edd25519 ia active.

	onOk := func() {
		defer c.view.Dialog.Close()

		label := gtkutil.GetEntryText(c.view.LabelEntry)
		typ, err := crypto.AddressTypeFromString(combo.GetActiveID())
		if err != nil {
			gtkutil.ShowError(err)

			return
		}

		password := ""
		if typ == crypto.AddressTypeEd25519Account ||
			typ == crypto.AddressTypeSecp256k1Account {
			pwd, ok := PasswordProvider(c.model)
			if !ok {
				return
			}
			password = pwd
		}

		_, err = c.model.NewAddress(typ, label, password)
		if err != nil {
			gtkutil.ShowError(err)

			return
		}
	}

	onCancel := func() { c.view.Dialog.Close() }

	c.view.ConnectSignals(map[string]any{
		"on_ok":     onOk,
		"on_cancel": onCancel,
	})

	c.view.Dialog.SetModal(true)
	gtkutil.RunDialog(c.view.Dialog)
}
