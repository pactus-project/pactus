//go:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
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

	view := &NetworkWidgetView{
		ViewBuilder:         builder,
		Box:                 builder.GetBoxObj("id_box_network"),
		LabelNetworkName:    builder.GetLabelObj("id_label_network_name"),
		LabelConnectedPeers: builder.GetLabelObj("id_label_connected_peers"),
		ColViewPeers:        builder.GetColumnViewObj("id_columnview_peers"),
	}

	gtkutil.ColumnViewSetDefaultProperties(view.ColViewPeers)

	return view
}
