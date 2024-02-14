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

	treeView *gtk.TreeView
	model    *walletModel
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
	treeView := getTreeViewObj(builder, "id_treeview_addresses")
	labelName := getLabelObj(builder, "id_label_wallet_name")
	labelLocation := getLabelObj(builder, "id_label_wallet_location")
	labelEncrypted := getLabelObj(builder, "id_label_wallet_encrypted")

	getToolButtonObj(builder, "id_button_new_address").SetIconWidget(AddIcon())
	getToolButtonObj(builder, "id_button_change_password").SetIconWidget(PasswordIcon())
	getToolButtonObj(builder, "id_button_show_seed").SetIconWidget(SeedIcon())

	labelName.SetText(model.wallet.Name())
	labelLocation.SetText(model.wallet.Path())
	if model.wallet.IsEncrypted() {
		labelEncrypted.SetText("Yes")
	} else {
		labelEncrypted.SetText("No")
	}

	colNo := createColumn("No", IDAddressesColumnNo)
	colAddress := createColumn("Address", IDAddressesColumnAddress)
	colLabel := createColumn("Label", IDAddressesColumnLabel)
	colBalance := createColumn("Balance", IDAddressesColumnBalance)
	colStake := createColumn("Stake", IDAddressesColumnStake)
	colScore := createColumn("Availability Score", IDAddressesColumnAvailabilityScore)

	treeView.AppendColumn(colNo)
	treeView.AppendColumn(colAddress)
	treeView.AppendColumn(colLabel)
	treeView.AppendColumn(colBalance)
	treeView.AppendColumn(colStake)
	treeView.AppendColumn(colScore)
	treeView.SetModel(model.ToTreeModel())

	w := &widgetWallet{
		Box:      box,
		treeView: treeView,
		model:    model,
	}

	menu, err := gtk.MenuNew()
	fatalErrorCheck(err)

	// "Update label" menu item
	item, err := gtk.MenuItemNewWithLabel("Update _Label")
	fatalErrorCheck(err)

	item.SetUseUnderline(true)
	item.Show()
	item.Connect("activate", func(item *gtk.MenuItem) bool {
		w.onUpdateLabel()

		return false
	})
	menu.Append(item)

	// "Address details" menu item
	item, err = gtk.MenuItemNewWithLabel("_Details")
	fatalErrorCheck(err)

	item.SetUseUnderline(true)
	item.Show()
	item.Connect("activate", func(item *gtk.MenuItem) bool {
		w.onShowDetails()

		return false
	})
	menu.Append(item)

	// "Private key" menu item
	item, err = gtk.MenuItemNewWithLabel("_Private key")
	fatalErrorCheck(err)

	item.SetUseUnderline(true)
	item.Show()
	item.Connect("activate", func(item *gtk.MenuItem) bool {
		w.onShowPrivateKey()

		return false
	})
	menu.Append(item)

	treeView.Connect("button-press-event",
		func(treeView *gtk.TreeView, event *gdk.Event) bool {
			eventButton := gdk.EventButtonNewFromEvent(event)
			if eventButton.Type() == gdk.EVENT_BUTTON_PRESS &&
				eventButton.Button() == gdk.BUTTON_SECONDARY {
				menu.PopupAtPointer(event)
			}

			return false
		})

	signals := map[string]interface{}{
		"on_new_address":     w.onNewAddress,
		"on_change_password": w.onChangePassword,
		"on_show_seed":       w.onShowSeed,
	}
	builder.ConnectSignals(signals)

	glib.TimeoutAdd(15000, w.timeout) // each 15 seconds

	return w, nil
}

func (ww *widgetWallet) onChangePassword() {
	changePassword(ww.model.wallet)
}

func (ww *widgetWallet) onNewAddress() {
	createAddress(ww)
}

func (ww *widgetWallet) onShowSeed() {
	password, ok := getWalletPassword(ww.model.wallet)
	if !ok {
		return
	}

	seed, err := ww.model.wallet.Mnemonic(password)
	errorCheck(err)

	showSeed(seed)
}

func (ww *widgetWallet) timeout() bool {
	ww.model.rebuildModel()

	return true
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
	selection, err := ww.treeView.GetSelection()
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
