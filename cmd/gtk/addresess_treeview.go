package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type addressesTreeView struct {
	*gtk.TreeView
}

// IDs to access the tree view columns by
const (
	ADDRESSES_COLUMN_NO = iota
	ADDRESSES_COLUMN_ADDRESS
	ADDRESSES_COLUMN_LABEL
	ADDRESSES_COLUMN_BALANCE
	ADDRESSES_COLUMN_STAKE
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

	return column
}

func buildAddressesTreeView(builder *gtk.Builder) *addressesTreeView {
	objTreeView, err := builder.GetObject("addresses_treeview")
	errorCheck(err)

	treeView, err := isTreeView(objTreeView)
	errorCheck(err)

	colNo := createColumn("No", ADDRESSES_COLUMN_NO)
	colAddress := createColumn("Address", ADDRESSES_COLUMN_ADDRESS)
	colLabel := createColumn("label", ADDRESSES_COLUMN_LABEL)
	colBalance := createColumn("balance", ADDRESSES_COLUMN_BALANCE)
	colStake := createColumn("Stake", ADDRESSES_COLUMN_STAKE)

	treeView.AppendColumn(colNo)
	treeView.AppendColumn(colAddress)
	treeView.AppendColumn(colLabel)
	treeView.AppendColumn(colBalance)
	treeView.AppendColumn(colStake)

	return &addressesTreeView{
		TreeView: treeView,
	}
}
