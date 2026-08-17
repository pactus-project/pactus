//go:build gtk

package gtkutil

import (
	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
)

func ColumnViewAppendTextColumn[T any](colView *gtk.ColumnView, title string, extractor func(T) string) {
	column := ColumnViewCreateTextColumn(title, extractor)
	column.SetResizable(true)
	column.SetExpand(false)

	colView.AppendColumn(column)
}

// ColumnViewAppendTextColumnEx appends a text column with control over the cell
// text alignment (xalign 0 = left, 1 = right), whether the column expands to
// fill remaining width, and an optional CSS class applied to the cell label.
func ColumnViewAppendTextColumnEx[T any](colView *gtk.ColumnView, title string,
	xalign float32, expand bool, cssClass string, extractor func(T) string,
) {
	factory := gtk.NewSignalListItemFactory()
	factory.ConnectSetup(func(obj *glib.Object) {
		cell := obj.Cast().(*gtk.ColumnViewCell)
		label := gtk.NewLabel("")
		label.SetHExpand(true)
		label.SetXAlign(xalign)
		if cssClass != "" {
			label.AddCSSClass(cssClass)
		}
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
	column.SetExpand(expand)
	column.SetResizable(true)

	colView.AppendColumn(column)
}

// ColumnViewAppendEllipsizedColumn appends an expanding column whose long text
// is ellipsized at the end, with the full value available on hover.
func ColumnViewAppendEllipsizedColumn[T any](colView *gtk.ColumnView, title string, extractor func(T) string) {
	factory := gtk.NewSignalListItemFactory()
	factory.ConnectSetup(func(obj *glib.Object) {
		cell := obj.Cast().(*gtk.ColumnViewCell)
		label := gtk.NewLabel("")
		label.SetHExpand(true)
		label.SetXAlign(0)
		label.SetEllipsize(pango.EllipsizeEnd)
		cell.SetChild(label)
	})
	factory.ConnectBind(func(obj *glib.Object) {
		cell := obj.Cast().(*gtk.ColumnViewCell)
		row := gioutil.ObjectValue[T](cell.Item())
		label := cell.Child().(*gtk.Label)
		value := extractor(row)
		label.SetText(value)
		label.SetTooltipText(value)
	})

	column := gtk.NewColumnViewColumn(title, &factory.ListItemFactory)
	column.SetTitle(title)
	column.SetExpand(true)
	column.SetResizable(true)

	colView.AppendColumn(column)
}

// ColumnViewAppendAddressColumn appends an expanding column that middle-
// ellipsizes long identifiers such as addresses and shows the full value on
// hover, so the table never needs horizontal scrolling.
func ColumnViewAppendAddressColumn[T any](colView *gtk.ColumnView, title string, extractor func(T) string) {
	factory := gtk.NewSignalListItemFactory()
	factory.ConnectSetup(func(obj *glib.Object) {
		cell := obj.Cast().(*gtk.ColumnViewCell)
		label := gtk.NewLabel("")
		label.SetHExpand(true)
		label.SetXAlign(0)
		label.SetEllipsize(pango.EllipsizeMiddle)
		label.AddCSSClass("cell-mono")
		cell.SetChild(label)
	})
	factory.ConnectBind(func(obj *glib.Object) {
		cell := obj.Cast().(*gtk.ColumnViewCell)
		row := gioutil.ObjectValue[T](cell.Item())
		label := cell.Child().(*gtk.Label)
		value := extractor(row)
		label.SetText(value)
		label.SetTooltipText(value)
	})

	column := gtk.NewColumnViewColumn(title, &factory.ListItemFactory)
	column.SetTitle(title)
	column.SetExpand(true)
	column.SetResizable(true)

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
	colView.SetShowRowSeparators(false)
	colView.SetShowColumnSeparators(false)
	colView.SetSingleClickActivate(true)
	colView.AddCSSClass("data-table")
}
