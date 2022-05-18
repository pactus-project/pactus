//go:build gtk

package main

import (
	_ "embed"
	"log"

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

	treeView.AppendColumn(colNo)
	treeView.AppendColumn(colAddress)
	treeView.AppendColumn(colLabel)
	treeView.AppendColumn(colBalance)
	treeView.AppendColumn(colStake)
	treeView.SetModel(model.ToTreeModel())

	w := &widgetWallet{
		Box:      box,
		treeView: treeView,
		model:    model,
	}

	menu, err := gtk.MenuNew()
	fatalErrorCheck(err)

	item, err := gtk.MenuItemNewWithLabel("Update _Label")
	fatalErrorCheck(err)

	item.SetUseUnderline(true)
	item.Show()
	item.Connect("activate", func(item *gtk.MenuItem) bool {
		w.onUpdateLabel()
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
		"on_new_address": w.onNewAddress,
	}
	builder.ConnectSignals(signals)

	glib.TimeoutAdd(10000, w.timeout)

	return w, nil
}

func (ww *widgetWallet) onNewAddress() {
	password, ok := getWalletPassword(ww.model.wallet)
	if !ok {
		return
	}

	err := ww.model.createAddress(password)
	fatalErrorCheck(err)
}

func (ww *widgetWallet) timeout() bool {
	err := ww.model.rebuildModel()
	fatalErrorCheck(err)

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

			err = ww.model.rebuildModel()
			fatalErrorCheck(err)
		}
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
			log.Printf("treeSelectionChangedCB: selected path: %s\n", addr)

			return addr
		}
	}
	return ""
}
