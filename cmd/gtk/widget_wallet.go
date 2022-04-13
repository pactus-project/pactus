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

func buildWidgetWallet() *widgetWallet {
	builder, err := gtk.BuilderNewFromString(string(uiWidgetWallet))
	errorCheck(err)

	objBox, err := builder.GetObject("id_box_wallet")
	errorCheck(err)

	box, err := isBox(objBox)
	errorCheck(err)

	objTreeView, err := builder.GetObject("id_treeview_addresses")
	errorCheck(err)

	treeView, err := isTreeView(objTreeView)
	errorCheck(err)

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

	w := &widgetWallet{
		Box:      box,
		treeView: treeView,
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

func (ww *widgetWallet) SetModel(model *walletModel) {
	ww.model = model
	ww.treeView.SetModel(model.ToTreeModel())
}
