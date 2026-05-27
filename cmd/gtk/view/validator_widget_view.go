//go111:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
)

type ValidatorWidgetView struct {
	ViewBuilder

	Box *gtk.Box

	ColViewValidators *gtk.ColumnView
}

func NewValidatorWidgetView() *ValidatorWidgetView {
	builder := NewViewBuilder(assets.ValidatorWidgetUI)

	colViewValidators := builder.GetColumnViewObj("id_columnview_validators")

	view := &ValidatorWidgetView{
		ViewBuilder: builder,
		Box:         builder.GetBoxObj("id_box_validator"),

		ColViewValidators: colViewValidators,
	}

	// // Build list store for validator table.
	// listStore := gtk.NewListStore([]glib.Type{
	// 	glib.TypeString, // no
	// 	glib.TypeString, // address
	// 	glib.TypeString, // number
	// 	glib.TypeString, // stake
	// 	glib.TypeString, // last bonding height
	// 	glib.TypeString, // last sortition height
	// 	glib.TypeString, // unbonding height
	// 	glib.TypeString, // availability score
	// })

	// view.listStore = listStore
	// view.TreeViewValidators.SetModel(&listStore.TreeModel)

	// // Columns.
	// colNo := createTextColumn("No", 0)
	// colAddress := createTextColumn("Address", 1)
	// colNumber := createTextColumn("Number", 2)
	// colStake := createTextColumn("Stake", 3)
	// colBondingHeight := createTextColumn("Bonding Height", 4)
	// colSortitionHeight := createTextColumn("Last Sortition Height", 5)
	// colUnbondingHeight := createTextColumn("Unbonding Height", 6)
	// colScore := createTextColumn("Availability Score", 7)

	// view.TreeViewValidators.AppendColumn(colNo)
	// view.TreeViewValidators.AppendColumn(colAddress)
	// view.TreeViewValidators.AppendColumn(colNumber)
	// view.TreeViewValidators.AppendColumn(colStake)
	// view.TreeViewValidators.AppendColumn(colBondingHeight)
	// view.TreeViewValidators.AppendColumn(colSortitionHeight)
	// view.TreeViewValidators.AppendColumn(colUnbondingHeight)
	// view.TreeViewValidators.AppendColumn(colScore)

	return view
}

// func (view *ValidatorWidgetView) ClearRows() {
// 	view.listStore.Clear()
// }

// func (view *ValidatorWidgetView) AppendRow(cols []int, values []any) {
// 	gtkutil.AppendRowToListStore(view.listStore, cols, values)
// }
