//go:build gtk

package controller

import (
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/autostart"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/version"
)

// Navigator owns dialog creation and UI flows (open -> run -> post-actions).
type Navigator struct {
	walletModel *model.WalletModel
	walletCtrl  *WalletWidgetController
	configModel *model.ConfigModel
	gtkApp      *gtk.Application

	autostart          *autostart.Entry
	autostartAction    *gio.SimpleAction
	onAutostartChanged func(enabled bool)

	themeAction    *gio.SimpleAction
	onThemeChanged func(mode gtkutil.ThemeMode)
}

func NewNavigator(gtkApp *gtk.Application,
	walletModel *model.WalletModel, walletCtrl *WalletWidgetController, configModel *model.ConfigModel,
	autostartEntry *autostart.Entry,
) *Navigator {
	return &Navigator{
		walletModel: walletModel,
		walletCtrl:  walletCtrl,
		gtkApp:      gtkApp,
		configModel: configModel,
		autostart:   autostartEntry,
	}
}

func (*Navigator) ShowAbout() {
	gtkutil.IdleAddSync(func() {
		dlg := view.NewAboutDialog()
		dlg.SetVersion(version.NodeVersion().StringWithAlias())
		gtkutil.ShowModalWindow(&dlg.Window)
	})
}

func (*Navigator) ShowAboutGtk() {
	gtkutil.IdleAddSync(func() {
		dlg := view.NewAboutGTKDialog()
		gtkutil.ShowModalWindow(&dlg.Window)
	})
}

func (n *Navigator) ShowConfigEditor() {
	gtkutil.IdleAddSync(func() {
		dlgView := view.NewConfigEditorDialogView()
		dlgCtrl := NewConfigEditorDialogController(dlgView, n.configModel)
		dlgCtrl.Show()
	})
}

func (n *Navigator) ShowWalletShowSeed() {
	gtkutil.IdleAddSync(func() {
		dlgView := view.NewWalletSeedDialogView()
		dlgCtrl := NewWalletSeedDialogController(dlgView, n.walletModel)
		dlgCtrl.Show()
	})
}

func (n *Navigator) ShowWalletNewAddress() {
	gtkutil.IdleAddSync(func() {
		dlgView := view.NewWalletCreateAddressDialogView()
		dlgCtrl := NewWalletCreateAddressDialogController(dlgView, n.walletModel)
		dlgCtrl.Show(func() {
			go n.walletCtrl.RefreshAddresses()
		})
	})
}

func (n *Navigator) ShowWalletSetDefaultFee() {
	gtkutil.IdleAddSync(func() {
		dlgView := view.NewWalletDefaultFeeDialogView()
		dlgCtrl := NewWalletDefaultFeeDialogController(dlgView, n.walletModel)
		dlgCtrl.Run(func() {
			go n.walletCtrl.RefreshInfo()
		})
	})
}

func (n *Navigator) ShowWalletChangePassword() {
	gtkutil.IdleAddSync(func() {
		dlgView := view.NewWalletChangePasswordDialogView()
		dlgCtrl := NewWalletChangePasswordDialogController(dlgView, n.walletModel)
		dlgCtrl.Show()
	})

	go func() {
		n.walletCtrl.RefreshInfo()
	}()
}

func (n *Navigator) ShowTransactionTransfer() {
	gtkutil.IdleAddSync(func() {
		dialogView := view.NewTxTransferDialogView()
		ctrl := NewTxTransferDialogController(dialogView, n.walletModel)
		ctrl.Show()
	})

	go func() {
		n.walletCtrl.RefreshTransactions()
	}()
}

func (n *Navigator) ShowTransactionBond() {
	gtkutil.IdleAddSync(func() {
		dialogView := view.NewTxBondDialogView()
		ctrl := NewTxBondDialogController(dialogView, n.walletModel)
		ctrl.Show()
	})

	go func() {
		n.walletCtrl.RefreshTransactions()
	}()
}

func (n *Navigator) ShowTransactionUnbond() {
	gtkutil.IdleAddSync(func() {
		dialogView := view.NewTxUnbondDialogView()
		ctrl := NewTxUnbondDialogController(dialogView, n.walletModel)
		ctrl.Show()
	})

	go func() {
		n.walletCtrl.RefreshTransactions()
	}()
}

func (n *Navigator) ShowTransactionWithdraw() {
	gtkutil.IdleAddSync(func() {
		dialogView := view.NewTxWithdrawDialogView()
		ctrl := NewTxWithdrawDialogController(dialogView, n.walletModel)
		ctrl.Show()
	})

	go func() {
		n.walletCtrl.RefreshTransactions()
	}()
}

func (n *Navigator) ShowTransactionAnchor() {
	gtkutil.IdleAddSync(func() {
		dialogView := view.NewTxAnchorDialogView()
		ctrl := NewTxAnchorDialogController(dialogView, n.walletModel)
		ctrl.Show()
	})

	go func() {
		n.walletCtrl.RefreshTransactions()
	}()
}

func (n *Navigator) BrowseWebsite() {
	n.openWebsite("https://pactus.org/")
}

func (n *Navigator) BrowseExplorer() {
	n.openWebsite("https://pactusscan.com/")
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

// AutostartEnabled reports whether Pactus GUI starts on system startup.
func (n *Navigator) AutostartEnabled() bool {
	return n.autostart.IsEnabled()
}

// OnAutostartChanged sets a function that is called on the UI thread after
// the autostart option changes, so other menus can show the new state.
func (n *Navigator) OnAutostartChanged(callback func(enabled bool)) {
	n.onAutostartChanged = callback
}

// SetAutostart turns starting Pactus GUI on system startup on or off.
// It must be called on the UI thread.
func (n *Navigator) SetAutostart(enabled bool) {
	if err := n.autostart.SetEnabled(enabled); err != nil {
		gtkutil.ShowErrorDialog(n.gtkApp.ActiveWindow(), err.Error(), nil)
	}

	enabled = n.autostart.IsEnabled()
	if n.autostartAction != nil {
		n.autostartAction.SetState(glib.NewVariantBoolean(enabled))
	}
	if n.onAutostartChanged != nil {
		n.onAutostartChanged(enabled)
	}
}

// OnThemeChanged sets a function that is called on the UI thread after the
// theme changes, so the title bar can show the new choice.
func (n *Navigator) OnThemeChanged(callback func(mode gtkutil.ThemeMode)) {
	n.onThemeChanged = callback
}

// SetThemeMode applies and persists the light, dark or system appearance.
// It must be called on the UI thread.
func (n *Navigator) SetThemeMode(mode gtkutil.ThemeMode) {
	gtkutil.ApplyThemeMode(mode)
	gtkutil.SaveThemeMode(mode)

	if n.themeAction != nil {
		n.themeAction.SetState(glib.NewVariantString(string(mode)))
	}
	if n.onThemeChanged != nil {
		n.onThemeChanged(mode)
	}
}

// ThemeMenu returns the Light / Dark / System choices as radio items. It is
// shared by the File menu and the title bar button, so both stay in sync.
func ThemeMenu() *gio.Menu {
	menu := gio.NewMenu()
	menu.Append("_Light", "app.theme::"+string(gtkutil.ThemeLight))
	menu.Append("_Dark", "app.theme::"+string(gtkutil.ThemeDark))
	menu.Append("_System", "app.theme::"+string(gtkutil.ThemeSystem))

	return menu
}

func (n *Navigator) CreateMenu(isLocal bool) *gio.Menu {
	// Helper to create an app action that calls a given function
	createAppAction := func(name string, callback func()) {
		action := gio.NewSimpleAction(name, nil)
		action.ConnectActivate(func(*glib.Variant) {
			callback()
		})
		action.SetEnabled(true)
		n.gtkApp.AddAction(action)
	}

	createAppAction("quit", n.Quit)
	if isLocal {
		createAppAction("config", n.ShowConfigEditor)
	}
	createAppAction("newaddress", n.ShowWalletNewAddress)
	createAppAction("changepassword", n.ShowWalletChangePassword)
	createAppAction("showseed", n.ShowWalletShowSeed)
	createAppAction("defaultfee", n.ShowWalletSetDefaultFee)
	createAppAction("transfer", n.ShowTransactionTransfer)
	createAppAction("bond", n.ShowTransactionBond)
	createAppAction("unbond", n.ShowTransactionUnbond)
	createAppAction("withdraw", n.ShowTransactionWithdraw)
	createAppAction("anchor", n.ShowTransactionAnchor)
	createAppAction("website", n.BrowseWebsite)
	createAppAction("explorer", n.BrowseExplorer)
	createAppAction("docs", n.BrowseDocs)
	createAppAction("about", n.ShowAbout)
	createAppAction("aboutgtk", n.ShowAboutGtk)

	// A stateful boolean action is shown as a check item in the menu.
	n.autostartAction = gio.NewSimpleActionStateful("autostart", nil,
		glib.NewVariantBoolean(n.AutostartEnabled()))
	n.autostartAction.ConnectActivate(func(*glib.Variant) {
		n.SetAutostart(!n.AutostartEnabled())
	})
	n.gtkApp.AddAction(n.autostartAction)

	// A stateful string action is shown as a group of radio items.
	n.themeAction = gio.NewSimpleActionStateful("theme", glib.NewVariantType("s"),
		glib.NewVariantString(string(gtkutil.LoadThemeMode())))
	n.themeAction.ConnectActivate(func(param *glib.Variant) {
		if param != nil {
			n.SetThemeMode(gtkutil.ThemeMode(param.String()))
		}
	})
	n.gtkApp.AddAction(n.themeAction)

	menubar := gio.NewMenu()

	// --- File Menu ---
	fileMenu := gio.NewMenu()
	fileMenu.Append("Edit _Config", "app.config")
	fileMenu.Append("Start on System _Startup", "app.autostart")
	fileMenu.AppendSubmenu("_Appearance", ThemeMenu())
	fileMenu.AppendSection("", gio.NewMenu()) // separator
	fileMenu.Append("_Quit", "app.quit")

	fileItem := gio.NewMenuItem("_File", "")
	fileItem.SetSubmenu(fileMenu)
	menubar.AppendItem(fileItem)

	// --- Wallet Menu ---
	walletMenu := gio.NewMenu()
	walletMenu.Append("_New Address", "app.newaddress")
	walletMenu.Append("Change _Password", "app.changepassword")
	walletMenu.Append("Show _Seed", "app.showseed")
	walletMenu.Append("Change Default _Fee", "app.defaultfee")

	walletItem := gio.NewMenuItem("_Wallet", "")
	walletItem.SetSubmenu(walletMenu)
	menubar.AppendItem(walletItem)

	// --- Transaction Menu ---
	transactionMenu := gio.NewMenu()
	transactionMenu.Append("_Transfer", "app.transfer")
	transactionMenu.Append("_Bond", "app.bond")
	transactionMenu.Append("_Unbond", "app.unbond")
	transactionMenu.Append("_Withdraw", "app.withdraw")
	transactionMenu.Append("_Anchor", "app.anchor")

	transactionItem := gio.NewMenuItem("_Transaction", "")
	transactionItem.SetSubmenu(transactionMenu)
	menubar.AppendItem(transactionItem)

	// --- Help Menu ---
	helpMenu := gio.NewMenu()
	helpMenu.Append("_Website", "app.website")
	helpMenu.Append("_Explorer", "app.explorer")
	helpMenu.Append("_Documentation", "app.docs")
	helpMenu.Append("About _Pactus", "app.about")
	helpMenu.Append("About _GTK", "app.aboutgtk")

	helpItem := gio.NewMenuItem("_Help", "")
	helpItem.SetSubmenu(helpMenu)
	menubar.AppendItem(helpItem)

	return menubar
}
