//go:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type CommitteeWidgetView struct {
	ViewBuilder

	Box *gtk.Box

	LabelCommitteeSize    *gtk.Label
	LabelCommitteePower   *gtk.Label
	LabelTotalPower       *gtk.Label
	LabelProtocolVersions *gtk.Label
	ColViewMembers        *gtk.ColumnView
}

func NewCommitteeWidgetView() *CommitteeWidgetView {
	builder := NewViewBuilder(assets.CommitteeWidgetUI)

	view := &CommitteeWidgetView{
		ViewBuilder:           builder,
		Box:                   builder.GetBoxObj("id_box_committee"),
		LabelCommitteeSize:    builder.GetLabelObj("id_label_committee_size"),
		LabelCommitteePower:   builder.GetLabelObj("id_label_committee_power"),
		LabelTotalPower:       builder.GetLabelObj("id_label_total_power"),
		LabelProtocolVersions: builder.GetLabelObj("id_label_protocol_versions"),
		ColViewMembers:        builder.GetColumnViewObj("id_columnview_committee_members"),
	}

	gtkutil.ColumnViewSetDefaultProperties(view.ColViewMembers)

	return view
}
