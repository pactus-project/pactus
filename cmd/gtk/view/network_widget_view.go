//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type NetworkWidgetView struct {
	ViewBuilder

	Box *gtk.Box

	LabelNetworkName    *gtk.Label
	LabelConnectedPeers *gtk.Label

	TreeViewPeers *gtk.TreeView
	listStore     *gtk.ListStore
}

func NewNetworkWidgetView() *NetworkWidgetView {
	builder := NewViewBuilder(assets.NetworkWidgetUI)

	treeViewPeers := builder.GetTreeViewObj("id_treeview_peers")

	view := &NetworkWidgetView{
		ViewBuilder: builder,
		Box:         builder.GetBoxObj("id_box_network"),

		LabelNetworkName:    builder.GetLabelObj("id_label_network_name"),
		LabelConnectedPeers: builder.GetLabelObj("id_label_connected_peers"),

		TreeViewPeers: treeViewPeers,
	}

	listStore, err := gtk.ListStoreNew(
		glib.TYPE_STRING, // No
		glib.TYPE_STRING, // Moniker
		glib.TYPE_STRING, // Address
		glib.TYPE_STRING, // Peer ID
		glib.TYPE_STRING, // Height
		glib.TYPE_STRING, // Agent
		glib.TYPE_STRING, // Direction
	)
	gtkutil.FatalErrorCheck(err)

	view.listStore = listStore
	view.TreeViewPeers.SetModel(listStore.ToTreeModel())

	colNo := createTextColumn("No", 0)
	colMoniker := createTextColumn("Moniker", 1)
	colAddress := createTextColumn("Address", 2)
	colPeerID := createTextColumn("Peer ID", 3)
	colHeight := createTextColumn("Height", 4)
	colAgent := createTextColumn("Agent", 5)
	colDirection := createTextColumn("Direction", 6)

	view.TreeViewPeers.AppendColumn(colNo)
	view.TreeViewPeers.AppendColumn(colMoniker)
	view.TreeViewPeers.AppendColumn(colAddress)
	view.TreeViewPeers.AppendColumn(colPeerID)
	view.TreeViewPeers.AppendColumn(colHeight)
	view.TreeViewPeers.AppendColumn(colAgent)
	view.TreeViewPeers.AppendColumn(colDirection)

	return view
}

func (view *NetworkWidgetView) ClearRows() {
	view.listStore.Clear()
}

func (view *NetworkWidgetView) AppendRow(cols []int, values []any) {
	iter := view.listStore.Append()
	_ = view.listStore.Set(iter, cols, values)
}
