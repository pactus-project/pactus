//go:build gtk

package controller

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/version"
)

// Navigator owns dialog creation and UI flows (open -> run -> post-actions).
type Navigator struct {
	walletModel *model.WalletModel
	walletCtrl  *WalletWidgetController
	gtkApp      *gtk.Application
}

func NewNavigator(gtkApp *gtk.Application,
	walletModel *model.WalletModel, walletCtrl *WalletWidgetController,
) *Navigator {
	return &Navigator{
		walletModel: walletModel,
		walletCtrl:  walletCtrl,
		gtkApp:      gtkApp,
	}
}

func (*Navigator) ShowAbout() {
	dlg := view.NewAboutDialog()
	dlg.SetVersion(version.NodeVersion().StringWithAlias())
	gtkutil.RunDialog(&dlg.Dialog)
}

func (*Navigator) ShowAboutGtk() {
	dlg := view.NewAboutGTKDialog()
	dlg.Dialog.SetModal(true)
	gtkutil.RunDialog(&dlg.Dialog)
}

func (n *Navigator) ShowWalletShowSeed() {
	password, ok := PasswordProvider(n.walletModel)
	if !ok {
		return
	}
	seed, err := n.walletModel.Mnemonic(password)
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

func (n *Navigator) ShowWalletNewAddress() {
	dlgView := view.NewWalletCreateAddressDialogView()
	dlgCtrl := NewWalletCreateAddressDialogController(dlgView, n.walletModel)
	dlgCtrl.Run()

	n.walletCtrl.RefreshAddresses()
}

func (n *Navigator) ShowWalletSetDefaultFee() {
	dlgView := view.NewWalletDefaultFeeDialogView()
	dlgCtrl := NewWalletDefaultFeeDialogController(dlgView, n.walletModel)
	dlgCtrl.Run()

	n.walletCtrl.RefreshInfo()
}

func (n *Navigator) ShowWalletChangePassword() {
	dlgView := view.NewWalletChangePasswordDialogView()
	dlgCtrl := NewWalletChangePasswordDialogController(dlgView, n.walletModel)
	dlgCtrl.Run()

	n.walletCtrl.RefreshInfo()
}

func (n *Navigator) ShowTransactionTransfer() {
	dialogView := view.NewTxTransferDialogView()
	ctrl := NewTxTransferDialogController(dialogView, n.walletModel)
	ctrl.Run()

	n.walletCtrl.RefreshTransactions()
}

func (n *Navigator) ShowTransactionBond() {
	dialogView := view.NewTxBondDialogView()
	ctrl := NewTxBondDialogController(dialogView, n.walletModel)
	ctrl.Run()

	n.walletCtrl.RefreshTransactions()
}

func (n *Navigator) ShowTransactionUnbond() {
	dialogView := view.NewTxUnbondDialogView()
	ctrl := NewTxUnbondDialogController(dialogView, n.walletModel)
	ctrl.Run()

	n.walletCtrl.RefreshTransactions()
}

func (n *Navigator) ShowTransactionWithdraw() {
	dialogView := view.NewTxWithdrawDialogView()
	ctrl := NewTxWithdrawDialogController(dialogView, n.walletModel)
	ctrl.Run()

	n.walletCtrl.RefreshTransactions()
}

func (n *Navigator) BrowseWebsite() {
	n.openWebsite("https://pactus.org/")
}

func (n *Navigator) BrowseExplorer() {
	n.openWebsite("https://pacviewer.com/")
}

func (n *Navigator) BrowseDocs() {
	n.openWebsite("https://docs.pactus.org/")
}

func (*Navigator) openWebsite(url string) {
	_ = gtkutil.OpenURLInBrowser(url)
}

func (n *Navigator) Quit() {
	n.gtkApp.Quit()
}
