//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
)

type NodeWidgetView struct {
	ViewBuilder

	Box *gtk.Box

	LabelConnectionType  *gtk.Label
	LabelConnectionValue *gtk.Label
	LabelNetwork         *gtk.Label
	LabelNetworkID       *gtk.Label
	LabelAgent           *gtk.Label
	LabelMoniker         *gtk.Label
	LabelIsPrune         *gtk.Label

	LabelClockOffset     *gtk.Label
	LabelLastBlockTime   *gtk.Label
	LabelLastBlockHeight *gtk.Label
	LabelBlocksLeft      *gtk.Label
	ProgressBarSynced    *gtk.ProgressBar
	LabelCommitteeSize   *gtk.Label
	LabelActiveValidator *gtk.Label
	LabelInCommittee     *gtk.Label
	LabelCommitteeStake  *gtk.Label
	LabelTotalStake      *gtk.Label
	LabelNumConnections  *gtk.Label
	LabelReachability    *gtk.Label
}

func NewNodeWidgetView() *NodeWidgetView {
	builder := NewViewBuilder(assets.NodeWidgetUI)

	view := &NodeWidgetView{
		ViewBuilder: builder,
		Box:         builder.GetBoxObj("id_box_node"),

		LabelConnectionType:  builder.GetLabelObj("id_label_connection_type"),
		LabelConnectionValue: builder.GetLabelObj("id_label_connection_value"),
		LabelNetwork:         builder.GetLabelObj("id_label_network"),
		LabelNetworkID:       builder.GetLabelObj("id_label_network_id"),
		LabelAgent:           builder.GetLabelObj("id_label_agent"),
		LabelMoniker:         builder.GetLabelObj("id_label_moniker"),
		LabelIsPrune:         builder.GetLabelObj("id_label_is_prune"),

		LabelClockOffset:     builder.GetLabelObj("id_label_clock_offset"),
		LabelLastBlockTime:   builder.GetLabelObj("id_label_last_block_time"),
		LabelLastBlockHeight: builder.GetLabelObj("id_label_last_block_height"),
		LabelBlocksLeft:      builder.GetLabelObj("id_label_blocks_left"),
		ProgressBarSynced:    builder.GetProgressBarObj("id_progress_synced"),
		LabelCommitteeSize:   builder.GetLabelObj("id_label_committee_size"),
		LabelActiveValidator: builder.GetLabelObj("id_label_active_validators"),
		LabelInCommittee:     builder.GetLabelObj("id_label_in_committee"),
		LabelCommitteeStake:  builder.GetLabelObj("id_label_committee_power"),
		LabelTotalStake:      builder.GetLabelObj("id_label_total_power"),
		LabelNumConnections:  builder.GetLabelObj("id_label_num_connections"),
		LabelReachability:    builder.GetLabelObj("id_label_reachability"),
	}

	return view
}
