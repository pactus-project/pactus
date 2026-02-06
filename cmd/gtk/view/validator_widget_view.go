//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type ValidatorWidgetView struct {
	ViewBuilder

	Box *gtk.Box

	TreeViewValidators *gtk.TreeView
	listStore          *gtk.ListStore
}

func NewValidatorWidgetView() *ValidatorWidgetView {
	builder := NewViewBuilder(assets.ValidatorWidgetUI)

	treeViewValidators := builder.GetTreeViewObj("id_treeview_validators")

	view := &ValidatorWidgetView{
		ViewBuilder: builder,
		Box:         builder.GetBoxObj("id_box_validator"),

		TreeViewValidators: treeViewValidators,
	}

	// Build list store for validator table.
	listStore, err := gtk.ListStoreNew(
		glib.TYPE_STRING, // no
		glib.TYPE_STRING, // address
		glib.TYPE_STRING, // number
		glib.TYPE_STRING, // stake
		glib.TYPE_STRING, // last bonding height
		glib.TYPE_STRING, // last sortition height
		glib.TYPE_STRING, // unbonding height
		glib.TYPE_STRING, // availability score
	)
	gtkutil.FatalErrorCheck(err)

	view.listStore = listStore
	view.TreeViewValidators.SetModel(listStore.ToTreeModel())

	// Columns.
	colNo := createTextColumn("No", 0)
	colAddress := createTextColumn("Address", 1)
	colNumber := createTextColumn("Number", 2)
	colStake := createTextColumn("Stake", 3)
	colBondingHeight := createTextColumn("Bonding Height", 4)
	colSortitionHeight := createTextColumn("Last Sortition Height", 5)
	colUnbondingHeight := createTextColumn("Unbonding Height", 6)
	colScore := createTextColumn("Availability Score", 7)

	view.TreeViewValidators.AppendColumn(colNo)
	view.TreeViewValidators.AppendColumn(colAddress)
	view.TreeViewValidators.AppendColumn(colNumber)
	view.TreeViewValidators.AppendColumn(colStake)
	view.TreeViewValidators.AppendColumn(colBondingHeight)
	view.TreeViewValidators.AppendColumn(colSortitionHeight)
	view.TreeViewValidators.AppendColumn(colUnbondingHeight)
	view.TreeViewValidators.AppendColumn(colScore)

	return view
}

func (view *ValidatorWidgetView) ClearRows() {
	view.listStore.Clear()
}

func (view *ValidatorWidgetView) AppendRow(cols []int, values []any) {
	iter := view.listStore.Append()
	_ = view.listStore.Set(iter, cols, values)
}
