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
	gtkutil.IdleAddSync(func() {
		dlg := view.NewAboutDialog()
		dlg.SetVersion(version.NodeVersion().StringWithAlias())
		gtkutil.RunDialog(&dlg.Dialog)
	})
}

func (*Navigator) ShowAboutGtk() {
	gtkutil.IdleAddSync(func() {
		dlg := view.NewAboutGTKDialog()
		dlg.Dialog.SetModal(true)
		gtkutil.RunDialog(&dlg.Dialog)
	})
}

func (n *Navigator) ShowWalletShowSeed() {
	gtkutil.IdleAddSync(func() {
		dlgView := view.NewWalletSeedDialogView()
		dlgCtrl := NewWalletSeedDialogController(dlgView, n.walletModel)
		dlgCtrl.Run()
	})
}

func (n *Navigator) ShowWalletNewAddress() {
	gtkutil.IdleAddSync(func() {
		dlgView := view.NewWalletCreateAddressDialogView()
		dlgCtrl := NewWalletCreateAddressDialogController(dlgView, n.walletModel)
		dlgCtrl.Run()
	})

	go func() {
		n.walletCtrl.RefreshAddresses()
	}()
}

func (n *Navigator) ShowWalletSetDefaultFee() {
	gtkutil.IdleAddSync(func() {
		dlgView := view.NewWalletDefaultFeeDialogView()
		dlgCtrl := NewWalletDefaultFeeDialogController(dlgView, n.walletModel)
		dlgCtrl.Run()
	})

	go func() {
		n.walletCtrl.RefreshInfo()
	}()
}

func (n *Navigator) ShowWalletChangePassword() {
	gtkutil.IdleAddSync(func() {
		dlgView := view.NewWalletChangePasswordDialogView()
		dlgCtrl := NewWalletChangePasswordDialogController(dlgView, n.walletModel)
		dlgCtrl.Run()
	})

	go func() {
		n.walletCtrl.RefreshInfo()
	}()
}

func (n *Navigator) ShowTransactionTransfer() {
	gtkutil.IdleAddSync(func() {
		dialogView := view.NewTxTransferDialogView()
		ctrl := NewTxTransferDialogController(dialogView, n.walletModel)
		ctrl.Run()
	})

	go func() {
		n.walletCtrl.RefreshTransactions()
	}()
}

func (n *Navigator) ShowTransactionBond() {
	gtkutil.IdleAddSync(func() {
		dialogView := view.NewTxBondDialogView()
		ctrl := NewTxBondDialogController(dialogView, n.walletModel)
		ctrl.Run()
	})

	go func() {
		n.walletCtrl.RefreshTransactions()
	}()
}

func (n *Navigator) ShowTransactionUnbond() {
	gtkutil.IdleAddSync(func() {
		dialogView := view.NewTxUnbondDialogView()
		ctrl := NewTxUnbondDialogController(dialogView, n.walletModel)
		ctrl.Run()
	})

	go func() {
		n.walletCtrl.RefreshTransactions()
	}()
}

func (n *Navigator) ShowTransactionWithdraw() {
	gtkutil.IdleAddSync(func() {
		dialogView := view.NewTxWithdrawDialogView()
		ctrl := NewTxWithdrawDialogController(dialogView, n.walletModel)
		ctrl.Run()
	})

	go func() {
		n.walletCtrl.RefreshTransactions()
	}()
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
