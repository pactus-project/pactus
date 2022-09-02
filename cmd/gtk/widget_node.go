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
	progressBarSynced    *gtk.ProgressBar
}

func buildWidgetNode(model *nodeModel, genesisTime time.Time) (*widgetNode, error) {
	builder, err := gtk.BuilderNewFromString(string(uiWidgetNode))
	if err != nil {
		return nil, err
	}

	box := getBoxObj(builder, "id_box_node")
	labelLocation := getLabelObj(builder, "id_label_working_directory")
	labelValidatorAddress := getLabelObj(builder, "id_label_validator_address")
	labelRewardAddress := getLabelObj(builder, "id_label_reward_address")

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	labelLocation.SetText(cwd)
	labelValidatorAddress.SetText(model.node.State().ValidatorAddress().String())
	labelRewardAddress.SetText(model.node.State().RewardAddress().String())

	w := &widgetNode{
		Box:                  box,
		model:                model,
		genesisTime:          genesisTime,
		labelLastBlockTime:   getLabelObj(builder, "id_label_last_block_time"),
		labelLastBlockHeight: getLabelObj(builder, "id_label_last_block_height"),
		labelBlocksLeft:      getLabelObj(builder, "id_label_blocks_left"),
		progressBarSynced:    getProgressBarObj(builder, "id_progress_synced"),
	}

	signals := map[string]interface{}{}
	builder.ConnectSignals(signals)

	glib.TimeoutAdd(1000, w.timeout)

	// Update widget for the first time
	w.timeout()
	return w, nil
}

func (wn *widgetNode) timeout() bool {
	lastBlockTime := wn.model.node.State().LastBlockTime()
	lastBlockHeight := wn.model.node.State().LastBlockHeight()
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

	return true
}
