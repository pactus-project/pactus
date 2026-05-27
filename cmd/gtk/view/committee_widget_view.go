//go111:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
)

type CommitteeWidgetView struct {
	ViewBuilder

	Box *gtk.Box

	LabelCommitteeSize    *gtk.Label
	LabelCommitteePower   *gtk.Label
	LabelTotalPower       *gtk.Label
	LabelProtocolVersions *gtk.Label

	ColViewMembers *gtk.ColumnView
}

func NewCommitteeWidgetView() *CommitteeWidgetView {
	builder := NewViewBuilder(assets.CommitteeWidgetUI)

	colViewMembers := builder.GetColumnViewObj("id_columnview_committee_members")

	view := &CommitteeWidgetView{
		ViewBuilder: builder,
		Box:         builder.GetBoxObj("id_box_committee"),

		LabelCommitteeSize:    builder.GetLabelObj("id_label_committee_size"),
		LabelCommitteePower:   builder.GetLabelObj("id_label_committee_power"),
		LabelTotalPower:       builder.GetLabelObj("id_label_total_power"),
		LabelProtocolVersions: builder.GetLabelObj("id_label_protocol_versions"),

		ColViewMembers: colViewMembers,
	}

	// columnView = gtk.NewColumnView(nil)

	// // Build list store for committee members table.
	// listStore := gio.NewListStore(glib.Type{
	// 	glib.TypeString, // No
	// 	glib.TypeString, // Address
	// 	glib.TypeString, // Number
	// 	glib.TypeString, // Stake
	// 	glib.TypeString, // Last Bonding Height
	// 	glib.TypeString, // Last Sortition Height
	// 	glib.TypeString, // Protocol Version
	// 	glib.TypeString, // Availability Score
	// })

	// view.listStore = listStore
	// view.ListViewMembers.SetModel(&listStore.TreeModel)

	// colNo := createTextColumn("No", 0)
	// colAddress := createTextColumn("Address", 1)
	// colNumber := createTextColumn("Number", 2)
	// colStake := createTextColumn("Stake", 3)
	// colBondingHeight := createTextColumn("Bonding Height", 4)
	// colSortitionHeight := createTextColumn("Sortition Height", 5)
	// colProtocolVersion := createTextColumn("Protocol", 6)
	// colScore := createTextColumn("Availability", 7)

	// view.ListViewMembers.AppendColumn(colNo)
	// view.ListViewMembers.AppendColumn(colAddress)
	// view.ListViewMembers.AppendColumn(colNumber)
	// view.ListViewMembers.AppendColumn(colStake)
	// view.ListViewMembers.AppendColumn(colBondingHeight)
	// view.ListViewMembers.AppendColumn(colSortitionHeight)
	// view.ListViewMembers.AppendColumn(colProtocolVersion)
	// view.ListViewMembers.AppendColumn(colScore)

	return view
}

// func (view *CommitteeWidgetView) ClearRows() {
// 	view.listStore.Clear()
// }

// func (view *CommitteeWidgetView) AppendRow(cols []int, values []any) {
// 	gtkutil.AppendRowToListStore(view.listStore, cols, values)
// }
