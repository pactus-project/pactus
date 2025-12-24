//go:build gtk

package controller

import (
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/types"
)

type WalletCreateAddressModel interface {
	NewAddress(addressType crypto.AddressType, label string, opts ...wallet.NewAddressOption) (*types.AddressInfo, error)
}

type WalletCreateAddressDialogController struct {
	view   *view.WalletCreateAddressDialogView
	model  WalletCreateAddressModel
	getPwd PasswordProvider
}

func NewWalletCreateAddressDialogController(
	view *view.WalletCreateAddressDialogView,
	model WalletCreateAddressModel,
	getPwd PasswordProvider,
) *WalletCreateAddressDialogController {
	return &WalletCreateAddressDialogController{view: view, model: model, getPwd: getPwd}
}

func (c *WalletCreateAddressDialogController) Run() {
	combo := c.view.AddressTypeCombo
	combo.Append(crypto.AddressTypeEd25519Account.String(), "ED25519 Account")
	combo.Append(crypto.AddressTypeBLSAccount.String(), "BLS Account")
	combo.Append(crypto.AddressTypeValidator.String(), "Validator")
	combo.SetActive(0)

	onOk := func() {
		c.view.ButtonOK.SetSensitive(false)
		defer c.view.ButtonOK.SetSensitive(true)

		label := gtkutil.GetEntryText(c.view.LabelEntry)
		typ := combo.GetActiveID()

		var err error
		switch typ {
		case crypto.AddressTypeEd25519Account.String():
			password, ok := c.getPwd()
			if !ok {
				return
			}
			_, err = c.model.NewAddress(
				crypto.AddressTypeEd25519Account,
				label,
				wallet.WithPassword(password),
			)
		case crypto.AddressTypeBLSAccount.String():
			_, err = c.model.NewAddress(crypto.AddressTypeBLSAccount, label)
		case crypto.AddressTypeValidator.String():
			_, err = c.model.NewAddress(crypto.AddressTypeValidator, label)
		default:
			return
		}

		if err != nil {
			gtkutil.ShowError(err)

			return
		}

		c.view.Dialog.Close()
	}

	onCancel := func() { c.view.Dialog.Close() }

	c.view.ConnectSignals(map[string]any{
		"on_ok":     onOk,
		"on_cancel": onCancel,
	})

	c.view.Dialog.SetModal(true)
	gtkutil.RunDialog(c.view.Dialog)
}
