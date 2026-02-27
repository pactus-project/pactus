//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
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

	TreeViewMembers *gtk.TreeView
	listStore       *gtk.ListStore
}

func NewCommitteeWidgetView() *CommitteeWidgetView {
	builder := NewViewBuilder(assets.CommitteeWidgetUI)

	treeViewMembers := builder.GetTreeViewObj("id_treeview_committee_members")

	view := &CommitteeWidgetView{
		ViewBuilder: builder,
		Box:         builder.GetBoxObj("id_box_committee"),

		LabelCommitteeSize:    builder.GetLabelObj("id_label_committee_size"),
		LabelCommitteePower:   builder.GetLabelObj("id_label_committee_power"),
		LabelTotalPower:       builder.GetLabelObj("id_label_total_power"),
		LabelProtocolVersions: builder.GetLabelObj("id_label_protocol_versions"),

		TreeViewMembers: treeViewMembers,
	}

	// Build list store for committee members table.
	listStore, err := gtk.ListStoreNew(
		glib.TYPE_STRING, // No
		glib.TYPE_STRING, // Address
		glib.TYPE_STRING, // Number
		glib.TYPE_STRING, // Stake
		glib.TYPE_STRING, // Last Bonding Height
		glib.TYPE_STRING, // Last Sortition Height
		glib.TYPE_STRING, // Protocol Version
		glib.TYPE_STRING, // Availability Score
	)
	gtkutil.FatalErrorCheck(err)

	view.listStore = listStore
	view.TreeViewMembers.SetModel(listStore.ToTreeModel())

	colNo := createTextColumn("No", 0)
	colAddress := createTextColumn("Address", 1)
	colNumber := createTextColumn("Number", 2)
	colStake := createTextColumn("Stake", 3)
	colBondingHeight := createTextColumn("Bonding Height", 4)
	colSortitionHeight := createTextColumn("Sortition Height", 5)
	colProtocolVersion := createTextColumn("Protocol", 6)
	colScore := createTextColumn("Availability", 7)

	view.TreeViewMembers.AppendColumn(colNo)
	view.TreeViewMembers.AppendColumn(colAddress)
	view.TreeViewMembers.AppendColumn(colNumber)
	view.TreeViewMembers.AppendColumn(colStake)
	view.TreeViewMembers.AppendColumn(colBondingHeight)
	view.TreeViewMembers.AppendColumn(colSortitionHeight)
	view.TreeViewMembers.AppendColumn(colProtocolVersion)
	view.TreeViewMembers.AppendColumn(colScore)

	return view
}

func (view *CommitteeWidgetView) ClearRows() {
	view.listStore.Clear()
}

func (view *CommitteeWidgetView) AppendRow(cols []int, values []any) {
	iter := view.listStore.Append()
	_ = view.listStore.Set(iter, cols, values)
}
