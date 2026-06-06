//go:build gtk

package gtkutil

import (
	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func ColumnViewAppendTextColumn[T any](colView *gtk.ColumnView, title string, extractor func(T) string) {
	column := ColumnViewCreateTextColumn(title, extractor)
	column.SetResizable(true)
	column.SetExpand(false)

	colView.AppendColumn(column)
}

func ColumnViewCreateTextColumn[T any](title string, extractor func(T) string) *gtk.ColumnViewColumn {
	factory := gtk.NewSignalListItemFactory()
	factory.ConnectSetup(func(obj *glib.Object) {
		cell := obj.Cast().(*gtk.ColumnViewCell)

		label := gtk.NewLabel("")
		cell.SetChild(label)
	})
	factory.ConnectBind(func(obj *glib.Object) {
		cell := obj.Cast().(*gtk.ColumnViewCell)
		row := gioutil.ObjectValue[T](cell.Item())

		label := cell.Child().(*gtk.Label)
		label.SetText(extractor(row))
	})

	column := gtk.NewColumnViewColumn(title, &factory.ListItemFactory)
	column.SetTitle(title)
	column.SetExpand(true)
	column.SetResizable(true)

	return column
}

func ColumnViewSetup[T any](colView *gtk.ColumnView, listModel *gioutil.ListModel[T]) {
	factory := gtk.NewSignalListItemFactory()

	factory.ConnectSetup(func(obj *glib.Object) {
		listItem := obj.Cast().(*gtk.ListItem)
		label := gtk.NewLabel("")
		listItem.SetChild(label)
	})

	factory.ConnectBind(func(obj *glib.Object) {
		listItem := obj.Cast().(*gtk.ListItem)
		label := listItem.Child().(*gtk.Label)

		// Get the row position and fetch the corresponding data
		position := listItem.Position()
		rowData := listModel.At(int(position))

		label.SetObjectProperty("row-data", rowData)
		label.SetObjectProperty("row-index", int(position))
	})

	columnViewSelectOnRightClick(colView)
}

func columnViewSelectOnRightClick(*gtk.ColumnView) {
	// TODO: complete me!
}

// ColumnViewGetSelectedItem gets the selection item from the ColumnView.
func ColumnViewGetSelectedItem[T any](colView *gtk.ColumnView, model *gioutil.ListModel[T]) T {
	selectionModel := colView.Model().Cast().(*gtk.SingleSelection)
	selectedPos := selectionModel.Selected()
	if selectedPos == gtk.InvalidListPosition {
		Logf("No item selected in ColumnView")
		var zeroValue T

		return zeroValue
	}
	Logf("Selected position in ColumnView: %d", selectedPos)

	return model.At(int(selectedPos))
}

func ColumnViewSetDefaultProperties(colView *gtk.ColumnView) {
	colView.SetShowRowSeparators(true)
	colView.SetShowColumnSeparators(true)
	colView.SetSingleClickActivate(true)
}
