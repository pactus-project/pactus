//go:build gtk

package main

import (
	_ "embed"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type widgetWallet struct {
	*gtk.Box

	treeView *gtk.TreeView
	model    *walletModel
}

//go:embed assets/ui/widget_wallet.ui
var uiWidgetWallet []byte

// IDs to access the tree view columns by
const (
	IDAddressesColumnNo = iota
	IDAddressesColumnAddress
	IDAddressesColumnLabel
	IDAddressesColumnBalance
	IDAddressesColumnStake
)

// Add a column to the tree view (during the initialization of the tree view)
func createColumn(title string, id int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Unable to create text cell renderer:", err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}

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

	signals := map[string]interface{}{
		"on_new_address": w.onNewAddress,
	}
	builder.ConnectSignals(signals)

	glib.TimeoutAdd(10000, w.timeout)

	return w, nil
}

func (ww *widgetWallet) onNewAddress() {
	password, ok := getWalletPassword(nil, ww.model.wallet)
	if !ok {
		return
	}

	err := ww.model.createAddress(password)
	errorCheck(nil, err)
}

func (wn *widgetWallet) timeout() bool {
	err := wn.model.rebuildModel()
	if err != nil {
		errorCheck(nil, err)
	}
	return true
}
