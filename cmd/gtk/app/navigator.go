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
// It is intentionally UI-level (it depends on view/controller/gtkutil).
type Navigator struct {
	model *model.WalletModel
}

func NewNavigator(model *model.WalletModel) *Navigator {
	return &Navigator{model: model}
}

func (n *Navigator) PasswordProvider() controller.TxPasswordProvider {
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
	dlg.SetModal(true)
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

func (n *Navigator) ShowCreateAddress(afterCreate func()) {
	dlgView := view.NewWalletCreateAddressDialogView()
	dlgCtrl := controller.NewWalletCreateAddressDialogController(dlgView, n.model, n.PasswordProvider())
	dlgCtrl.Run(afterCreate)
}

func (n *Navigator) ShowSetDefaultFee(afterChange func()) {
	dlgView := view.NewWalletDefaultFeeDialogView()
	dlgCtrl := controller.NewWalletDefaultFeeDialogController(dlgView)
	dlgCtrl.Run(n.model, afterChange)
}

func (n *Navigator) ShowChangePassword(afterChange func()) {
	dlgView := view.NewWalletChangePasswordDialogView()
	dlgCtrl := controller.NewWalletChangePasswordDialogController(dlgView)
	dlgCtrl.Run(n.model, afterChange)
}

func (n *Navigator) ShowUpdateLabel(address string, afterChange func()) {
	oldLabel := n.model.Label(address)
	dlgView := view.NewAddressLabelDialogView()
	dlgCtrl := controller.NewAddressLabelDialogController(dlgView)
	newLabel, ok := dlgCtrl.Run(oldLabel)
	if !ok {
		return
	}
	if err := n.model.SetLabel(address, newLabel); err != nil {
		gtkutil.ShowError(err)

		return
	}
	if afterChange != nil {
		afterChange()
	}
}

func (n *Navigator) ShowAddressDetails(address string) {
	dlgView := view.NewAddressDetailsDialogView()
	dlgCtrl := controller.NewAddressDetailsDialogController(dlgView)
	_ = dlgCtrl.Run(n.model, address)
}

func (n *Navigator) ShowPrivateKey(address string) {
	dlgView := view.NewAddressPrivateKeyDialogView()
	dlgCtrl := controller.NewAddressPrivateKeyDialogController(dlgView, n.model, n.PasswordProvider())
	dlgCtrl.Run(address)
}

func (n *Navigator) ShowTransferTx() {
	dialogView := view.NewTxTransferDialogView()
	ctrl := controller.NewTransferTxController(dialogView, n.model, n.PasswordProvider())
	ctrl.BindAndRun()
}

func (n *Navigator) ShowBondTx() {
	dialogView := view.NewTxBondDialogView()
	ctrl := controller.NewBondTxController(dialogView, n.model, n.PasswordProvider())
	ctrl.BindAndRun()
}

func (n *Navigator) ShowUnbondTx() {
	dialogView := view.NewTxUnbondDialogView()
	ctrl := controller.NewUnbondTxController(dialogView, n.model, n.PasswordProvider())
	ctrl.BindAndRun()
}

func (n *Navigator) ShowWithdrawTx() {
	dialogView := view.NewTxWithdrawDialogView()
	ctrl := controller.NewWithdrawTxController(dialogView, n.model, n.PasswordProvider())
	ctrl.BindAndRun()
}
