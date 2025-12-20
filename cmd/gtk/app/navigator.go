//go:build gtk

package app

import (
	"github.com/pactus-project/pactus/cmd/gtk/controller"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/version"
)

// Navigator owns dialog creation and UI flows (open -> run -> post-actions).
type Navigator struct {
	model *model.WalletModel
}

func NewNavigator(model *model.WalletModel) *Navigator {
	return &Navigator{model: model}
}

func (n *Navigator) PasswordProvider() controller.PasswordProvider {
	return func() (string, bool) {
		if n.model == nil || !n.model.IsEncrypted() {
			return "", true
		}

		pwdView := view.NewWalletPasswordDialogView()
		pwdCtrl := controller.NewWalletPasswordDialogController(pwdView)

		return pwdCtrl.Run()
	}
}

func (*Navigator) ShowAbout() {
	dlg := view.NewAboutDialog()
	dlg.SetVersion(version.NodeVersion().StringWithAlias())
	gtkutil.RunDialog(&dlg.Dialog)
}

func (*Navigator) ShowAboutGTK() {
	dlg := view.NewAboutGTKDialog()
	dlg.Dialog.SetModal(true)
	gtkutil.RunDialog(&dlg.Dialog)
}

func (n *Navigator) ShowSeed() {
	password, ok := n.PasswordProvider()()
	if !ok {
		return
	}
	seed, err := n.model.Mnemonic(password)
	if err != nil {
		gtkutil.ShowError(err)

		return
	}
	dlgView := view.NewWalletSeedDialogView()
	gtkutil.SetTextViewContent(dlgView.TextView, seed)
	dlgView.ConnectSignals(map[string]any{"on_close": func() { dlgView.Dialog.Close() }})
	dlgView.Dialog.SetModal(true)
	gtkutil.RunDialog(dlgView.Dialog)
}

func (n *Navigator) ShowCreateAddress() {
	dlgView := view.NewWalletCreateAddressDialogView()
	dlgCtrl := controller.NewWalletCreateAddressDialogController(dlgView, n.model, n.PasswordProvider())
	dlgCtrl.Run()
}

func (n *Navigator) ShowSetDefaultFee() {
	dlgView := view.NewWalletDefaultFeeDialogView()
	dlgCtrl := controller.NewWalletDefaultFeeDialogController(dlgView, n.model)
	dlgCtrl.Run()
}

func (n *Navigator) ShowChangePassword() {
	dlgView := view.NewWalletChangePasswordDialogView()
	dlgCtrl := controller.NewWalletChangePasswordDialogController(dlgView, n.model)
	dlgCtrl.Run()
}

func (n *Navigator) ShowUpdateLabel(address string) {
	dlgView := view.NewAddressLabelDialogView()
	dlgCtrl := controller.NewAddressLabelDialogController(dlgView, n.model)
	dlgCtrl.Run(address)
}

func (n *Navigator) ShowAddressDetails(address string) {
	dlgView := view.NewAddressDetailsDialogView()
	dlgCtrl := controller.NewAddressDetailsDialogController(dlgView, n.model)
	dlgCtrl.Run(address)
}

func (n *Navigator) ShowPrivateKey(address string) {
	dlgView := view.NewAddressPrivateKeyDialogView()
	dlgCtrl := controller.NewAddressPrivateKeyDialogController(dlgView, n.model, n.PasswordProvider())
	dlgCtrl.Run(address)
}

func (n *Navigator) ShowTransferTx() {
	dialogView := view.NewTxTransferDialogView()
	ctrl := controller.NewTxTransferDialogController(dialogView, n.model, n.PasswordProvider())
	ctrl.Run()
}

func (n *Navigator) ShowBondTx() {
	dialogView := view.NewTxBondDialogView()
	ctrl := controller.NewTxBondDialogController(dialogView, n.model, n.PasswordProvider())
	ctrl.Run()
}

func (n *Navigator) ShowUnbondTx() {
	dialogView := view.NewTxUnbondDialogView()
	ctrl := controller.NewTxUnbondDialogController(dialogView, n.model, n.PasswordProvider())
	ctrl.Run()
}

func (n *Navigator) ShowWithdrawTx() {
	dialogView := view.NewTxWithdrawDialogView()
	ctrl := controller.NewTxWithdrawDialogController(dialogView, n.model, n.PasswordProvider())
	ctrl.Run()
}
