package main

import (
	_ "embed"
	"log"

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
	ID_ADDRESSES_COLUMN_NO = iota
	ID_ADDRESSES_COLUMN_ADDRESS
	ID_ADDRESSES_COLUMN_LABEL
	ID_ADDRESSES_COLUMN_BALANCE
	ID_ADDRESSES_COLUMN_STAKE
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

func buildWidgetWallet(model *walletModel) *widgetWallet {
	builder, err := gtk.BuilderNewFromString(string(uiWidgetWallet))
	errorCheck(err)

	box := getBoxObj(builder, "id_box_wallet")
	treeView := getTreeViewObj(builder, "id_treeview_addresses")
	nameLabel := getLabelObj(builder, "id_wallet_name")
	nameLocation := getLabelObj(builder, "id_wallet_location")
	nameEncrypted := getLabelObj(builder, "id_wallet_encrypted")

	nameLabel.SetText(model.wallet.Name())
	nameLocation.SetText(model.wallet.Path())
	if model.wallet.IsEncrypted() {
		nameEncrypted.SetText("Yes")
	} else {
		nameEncrypted.SetText("No")
	}

	colNo := createColumn("No", ID_ADDRESSES_COLUMN_NO)
	colAddress := createColumn("Address", ID_ADDRESSES_COLUMN_ADDRESS)
	colLabel := createColumn("Label", ID_ADDRESSES_COLUMN_LABEL)
	colBalance := createColumn("Balance", ID_ADDRESSES_COLUMN_BALANCE)
	colStake := createColumn("Stake", ID_ADDRESSES_COLUMN_STAKE)

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

	return w
}

func (ww *widgetWallet) onNewAddress() {
	password, ok := getWalletPassword(nil, ww.model.wallet)
	if !ok {
		return
	}

	ww.model.createAddress(password)

}
