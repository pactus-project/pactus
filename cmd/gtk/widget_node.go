//go:build gtk

package main

import (
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/util"
)

//go:embed assets/ui/widget_node.ui
var uiWidgetNode []byte

type widgetNode struct {
	*gtk.Box

	genesisTime          time.Time // TODO: move this logic to the state
	model                *nodeModel
	labelLastBlockTime   *gtk.Label
	labelLastBlockHeight *gtk.Label
	labelBlocksLeft      *gtk.Label
	labelCommitteeSize   *gtk.Label
	labelInCommittee     *gtk.Label
	labelCommitteeStake  *gtk.Label
	labelTotalStake      *gtk.Label
	progressBarSynced    *gtk.ProgressBar
}

func buildWidgetNode(model *nodeModel, genesisTime time.Time) (*widgetNode, error) {
	builder, err := gtk.BuilderNewFromString(string(uiWidgetNode))
	if err != nil {
		return nil, err
	}

	box := getBoxObj(builder, "id_box_node")
	labelLocation := getLabelObj(builder, "id_label_working_directory")

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	labelLocation.SetText(cwd)

	w := &widgetNode{
		Box:                  box,
		model:                model,
		genesisTime:          genesisTime,
		labelLastBlockTime:   getLabelObj(builder, "id_label_last_block_time"),
		labelLastBlockHeight: getLabelObj(builder, "id_label_last_block_height"),
		labelBlocksLeft:      getLabelObj(builder, "id_label_blocks_left"),
		progressBarSynced:    getProgressBarObj(builder, "id_progress_synced"),
		labelCommitteeSize:   getLabelObj(builder, "id_label_committee_size"),
		labelInCommittee:     getLabelObj(builder, "id_label_in_committee"),
		labelCommitteeStake:  getLabelObj(builder, "id_label_committee_stake"),
		labelTotalStake:      getLabelObj(builder, "id_label_total_stake"),
	}

	signals := map[string]interface{}{}
	builder.ConnectSignals(signals)

	glib.TimeoutAdd(1000, w.timeout1)
	glib.TimeoutAdd(10000, w.timeout10)

	// Update widget for the first time
	w.timeout1()
	w.timeout10()
	return w, nil
}

func (wn *widgetNode) timeout1() bool {
	// updating gui in another thread, this will fix "Not Responding" issue on Windows
	go func() {
		lastBlockTime := wn.model.node.State().LastBlockTime()
		lastBlockHeight := wn.model.node.State().LastBlockHeight()

		// Fixing sudden panic
		// https://github.com/gotk3/gotk3/issues/686
		glib.IdleAdd(func() bool {
			wn.labelLastBlockTime.SetText(lastBlockTime.Format("02 Jan 06 15:04:05 MST"))
			wn.labelLastBlockHeight.SetText(strconv.FormatInt(int64(lastBlockHeight), 10))

			// TODO move this logic to state
			nowUnix := time.Now().Unix()
			lastBlockTimeUnix := lastBlockTime.Unix()
			genTimeUnix := wn.genesisTime.Unix()

			percentage := float64(lastBlockTimeUnix-genTimeUnix) / float64(nowUnix-genTimeUnix)
			wn.progressBarSynced.SetFraction(percentage)
			wn.progressBarSynced.SetText(fmt.Sprintf("%s %%",
				strconv.FormatFloat(percentage*100, 'f', 2, 64)))

			blocksLeft := (nowUnix - lastBlockTimeUnix) / 10
			wn.labelBlocksLeft.SetText(strconv.FormatInt(blocksLeft, 10))

			return false
		})
	}()

	return true
}

func (wn *widgetNode) timeout10() bool {
	go func() {
		committeeSize := wn.model.node.State().Params().CommitteeSize
		committeeStake := wn.model.node.State().CommitteePower()
		totalStake := wn.model.node.State().TotalPower()
		isInCommittee := "No"
		// if wn.model.node.State().IsInCommittee(
		// 	wn.model.node.State().ValidatorAddress()) {
		// 	isInCommittee = "Yes"
		// }

		glib.IdleAdd(func() bool {
			wn.labelCommitteeSize.SetText(fmt.Sprintf("%v", committeeSize))
			wn.labelCommitteeStake.SetText(util.ChangeToString(committeeStake))
			wn.labelTotalStake.SetText(util.ChangeToString(totalStake))
			wn.labelInCommittee.SetText(fmt.Sprintf("%v", isInCommittee))

			return false
		})
	}()

	return true
}
