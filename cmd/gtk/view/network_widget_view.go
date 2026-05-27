//go111:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
)

type NetworkWidgetView struct {
	ViewBuilder

	Box *gtk.Box

	LabelNetworkName    *gtk.Label
	LabelConnectedPeers *gtk.Label

	ColViewPeers *gtk.ColumnView
}

func NewNetworkWidgetView() *NetworkWidgetView {
	builder := NewViewBuilder(assets.NetworkWidgetUI)

	colViewPeers := builder.GetColumnViewObj("id_columnview_peers")

	view := &NetworkWidgetView{
		ViewBuilder: builder,
		Box:         builder.GetBoxObj("id_box_network"),

		LabelNetworkName:    builder.GetLabelObj("id_label_network_name"),
		LabelConnectedPeers: builder.GetLabelObj("id_label_connected_peers"),

		ColViewPeers: colViewPeers,
	}

	// listStore := gtk.NewListStore([]glib.Type{
	// 	glib.TypeString, // No
	// 	glib.TypeString, // Moniker
	// 	glib.TypeString, // Address
	// 	glib.TypeString, // Peer ID
	// 	glib.TypeString, // Height
	// 	glib.TypeString, // Agent
	// 	glib.TypeString, // Direction
	// })

	// view.listStore = listStore
	// view.TreeViewPeers.SetModel(&listStore.TreeModel)

	// colNo := createTextColumn("No", 0)
	// colMoniker := createTextColumn("Moniker", 1)
	// colAddress := createTextColumn("Address", 2)
	// colPeerID := createTextColumn("Peer ID", 3)
	// colHeight := createTextColumn("Height", 4)
	// colAgent := createTextColumn("Agent", 5)
	// colDirection := createTextColumn("Direction", 6)

	// view.TreeViewPeers.AppendColumn(colNo)
	// view.TreeViewPeers.AppendColumn(colMoniker)
	// view.TreeViewPeers.AppendColumn(colAddress)
	// view.TreeViewPeers.AppendColumn(colPeerID)
	// view.TreeViewPeers.AppendColumn(colHeight)
	// view.TreeViewPeers.AppendColumn(colAgent)
	// view.TreeViewPeers.AppendColumn(colDirection)

	return view
}

// func (view *NetworkWidgetView) ClearRows() {
// 	view.listStore.Clear()
// }

// func (view *NetworkWidgetView) AppendRow(cols []int, values []any) {
// 	gtkutil.AppendRowToListStore(view.listStore, cols, values)
// }
