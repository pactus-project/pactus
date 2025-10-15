//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// IDs to access the tree view columns.
const (
	IDAddressesColumnNo = iota
	IDAddressesColumnAddress
	IDAddressesColumnLabel
	IDAddressesColumnBalance
	IDAddressesColumnStake
	IDAddressesColumnAvailabilityScore
)

//go:embed assets/ui/widget_wallet.ui
var uiWidgetWallet []byte

type widgetWallet struct {
	*gtk.Box

	treeViewWallet    *gtk.TreeView
	labelTotalBalance *gtk.Label
	labelDefaultFee   *gtk.Label
	labelTotalStake   *gtk.Label
	labelEncrypted    *gtk.Label
	model             *walletModel

	// Timeout ID for cleanup
	timeoutID glib.SourceHandle
}

// Add a column to the tree view (during the initialization of the tree view).
func createColumn(title string, id int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	fatalErrorCheck(err)

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	fatalErrorCheck(err)

	column.SetResizable(true)

	return column
}

func buildWidgetWallet(model *walletModel) (*widgetWallet, error) {
	builder, err := gtk.BuilderNewFromString(string(uiWidgetWallet))
	if err != nil {
		return nil, err
	}

	box := getBoxObj(builder, "id_box_wallet")
	treeViewWallet := getTreeViewObj(builder, "id_treeview_addresses")
	labelName := getLabelObj(builder, "id_label_wallet_name")
	labelLocation := getLabelObj(builder, "id_label_wallet_location")
	labelEncrypted := getLabelObj(builder, "id_label_wallet_encrypted")
	labelTotalBalance := getLabelObj(builder, "id_label_wallet_total_balance")
	labelTotalStake := getLabelObj(builder, "id_label_wallet_total_stake")
	labelDefaultFee := getLabelObj(builder, "id_label_wallet_default_fee")

	getToolButtonObj(builder, "id_button_new_address").SetIconWidget(AddIcon())
	getToolButtonObj(builder, "id_button_set_default_fee").SetIconWidget(FeeIcon())
	getToolButtonObj(builder, "id_button_change_password").SetIconWidget(PasswordIcon())
	getToolButtonObj(builder, "id_button_show_seed").SetIconWidget(SeedIcon())

	labelName.SetText(model.wallet.Name())
	labelLocation.SetText(model.wallet.Path())

	colNo := createColumn("No", IDAddressesColumnNo)
	colAddress := createColumn("Address", IDAddressesColumnAddress)
	colLabel := createColumn("Label", IDAddressesColumnLabel)
	colBalance := createColumn("Balance", IDAddressesColumnBalance)
	colStake := createColumn("Stake", IDAddressesColumnStake)
	colScore := createColumn("Availability Score", IDAddressesColumnAvailabilityScore)

	treeViewWallet.AppendColumn(colNo)
	treeViewWallet.AppendColumn(colAddress)
	treeViewWallet.AppendColumn(colLabel)
	treeViewWallet.AppendColumn(colBalance)
	treeViewWallet.AppendColumn(colStake)
	treeViewWallet.AppendColumn(colScore)
	treeViewWallet.SetModel(model.ToTreeModel())

	wdgWallet := &widgetWallet{
		Box:               box,
		treeViewWallet:    treeViewWallet,
		labelTotalBalance: labelTotalBalance,
		labelTotalStake:   labelTotalStake,
		labelDefaultFee:   labelDefaultFee,
		labelEncrypted:    labelEncrypted,
		model:             model,
	}

	menu, err := gtk.MenuNew()
	fatalErrorCheck(err)

	// "Update label" menu item
	item, err := gtk.MenuItemNewWithLabel("Update _Label")
	fatalErrorCheck(err)

	item.SetUseUnderline(true)
	item.Show()
	item.Connect("activate", func(_ *gtk.MenuItem) bool {
		wdgWallet.onUpdateLabel()

		return false
	})
	menu.Append(item)

	// "Address details" menu item
	item, err = gtk.MenuItemNewWithLabel("_Details")
	fatalErrorCheck(err)

	item.SetUseUnderline(true)
	item.Show()
	item.Connect("activate", func(_ *gtk.MenuItem) bool {
		wdgWallet.onShowDetails()

		return false
	})
	menu.Append(item)

	// "Private key" menu item
	item, err = gtk.MenuItemNewWithLabel("_Private key")
	fatalErrorCheck(err)

	item.SetUseUnderline(true)
	item.Show()
	item.Connect("activate", func(_ *gtk.MenuItem) bool {
		wdgWallet.onShowPrivateKey()

		return false
	})
	menu.Append(item)

	treeViewWallet.Connect("button-press-event",
		func(_ *gtk.TreeView, event *gdk.Event) bool {
			eventButton := gdk.EventButtonNewFromEvent(event)
			if eventButton.Type() == gdk.EVENT_BUTTON_PRESS &&
				eventButton.Button() == gdk.BUTTON_SECONDARY {
				menu.PopupAtPointer(event)
			}

			return false
		})

	// Add double-click handler to open address details dialog
	treeViewWallet.Connect("row-activated",
		func(_ *gtk.TreeView, _ *gtk.TreePath, _ *gtk.TreeViewColumn) {
			wdgWallet.onShowDetails()
		})

	signals := map[string]any{
		"on_new_address":     wdgWallet.onNewAddress,
		"on_set_default_fee": wdgWallet.onSetDefaultFee,
		"on_change_password": wdgWallet.onChangePassword,
		"on_show_seed":       wdgWallet.onShowSeed,
	}
	builder.ConnectSignals(signals)

	wdgWallet.timeoutID = glib.TimeoutAdd(15000, wdgWallet.timeout) // each 15 seconds

	wdgWallet.rebuild()

	return wdgWallet, nil
}

func (ww *widgetWallet) onChangePassword() {
	changePassword(ww)
}

func (ww *widgetWallet) onNewAddress() {
	createAddress(ww)
}

func (ww *widgetWallet) onSetDefaultFee() {
	setDefaultFee(ww)
}

func (ww *widgetWallet) onShowSeed() {
	password, ok := getWalletPassword(ww.model.wallet)
	if !ok {
		return
	}

	seed, err := ww.model.wallet.Mnemonic(password)
	if err != nil {
		showError(err)

		return
	}

	showSeed(seed)
}

func (ww *widgetWallet) timeout() bool {
	ww.rebuild()

	return true
}

func (ww *widgetWallet) rebuild() {
	ww.model.rebuildModel()

	totalBalance, _ := ww.model.wallet.TotalBalance()
	totalStake, _ := ww.model.wallet.TotalStake()

	// Get wallet info and set default fee
	walletInfo := ww.model.wallet.Info()

	ww.labelTotalBalance.SetText(totalBalance.String())
	ww.labelTotalStake.SetText(totalStake.String())
	ww.labelDefaultFee.SetText(walletInfo.DefaultFee.String())

	if walletInfo.Encrypted {
		ww.labelEncrypted.SetText("Yes")
	} else {
		ww.labelEncrypted.SetText("No")
	}
}

func (ww *widgetWallet) onUpdateLabel() {
	addr := ww.getSelectedAddress()
	if addr != "" {
		oldLabel := ww.model.wallet.Label(addr)
		newLabel, ok := getAddressLabel(oldLabel)
		if ok {
			err := ww.model.wallet.SetLabel(addr, newLabel)
			fatalErrorCheck(err)

			err = ww.model.wallet.Save()
			fatalErrorCheck(err)

			ww.model.rebuildModel()
		}
	}
}

func (ww *widgetWallet) onShowDetails() {
	addr := ww.getSelectedAddress()
	if addr != "" {
		showAddressDetails(ww.model.wallet, addr)
	}
}

func (ww *widgetWallet) onShowPrivateKey() {
	addr := ww.getSelectedAddress()
	if addr != "" {
		showAddressPrivateKey(ww.model.wallet, addr)
	}
}

func (ww *widgetWallet) getSelectedAddress() string {
	selection, err := ww.treeViewWallet.GetSelection()
	fatalErrorCheck(err)

	if selection != nil {
		model, iter, ok := selection.GetSelected()
		if ok {
			path, err := model.(*gtk.TreeModel).GetValue(iter, IDAddressesColumnAddress)
			fatalErrorCheck(err)

			addr, err := path.GetString()
			fatalErrorCheck(err)

			return addr
		}
	}

	return ""
}

// cleanup cancels all timeouts to prevent memory leaks and potential panics.
func (ww *widgetWallet) cleanup() {
	if ww.timeoutID != 0 {
		glib.SourceRemove(ww.timeoutID)
		ww.timeoutID = 0
	}
}
