//go:build gtk

package controller

import (
	"context"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/ezex-io/gopkg/scheduler"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util/logger"
)

// clockOutOfSyncThreshold is the clock offset above which we show a warning.
const clockOutOfSyncThreshold = 5 * time.Second

type NodeWidgetController struct {
	view  *view.NodeWidgetView
	model *model.NodeModel
}

type nodeWidgetSnapshot struct {
	committeeSize    int
	committeeStake   amount.Amount
	totalStake       amount.Amount
	activeValidators int32
	numConnections   string
	reachability     string
	inCommittee      bool
	clockOffset      time.Duration
	clockOffsetErr   error
}

func NewNodeWidgetController(view *view.NodeWidgetView, model *model.NodeModel) *NodeWidgetController {
	return &NodeWidgetController{view: view, model: model}
}

func (c *NodeWidgetController) BuildView(ctx context.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	c.view.LabelWorkingDirectory.SetText(cwd)

	nodeInfo, err := c.model.GetNodeInfo()
	if err != nil {
		return err
	}
	chainInfo, err := c.model.GetBlockchainInfo()
	if err != nil {
		return err
	}

	c.view.LabelNetwork.SetText(nodeInfo.NetworkName)
	c.view.LabelNetworkID.SetText(nodeInfo.PeerId)
	c.view.LabelMoniker.SetText(nodeInfo.Moniker)
	c.view.LabelIsPrune.SetText(strconv.FormatBool(chainInfo.IsPruned))

	c.view.ConnectSignals(map[string]any{})

	scheduler.Every(ctx, time.Second).Do(c.timeout1)
	scheduler.Every(ctx, 10*time.Second).Do(c.timeout10)

	// Initial refresh.
	c.timeout1()
	c.timeout10()

	return nil
}

// syncProgressWindowBlocks is the number of blocks behind that maps to 0% sync progress (~10 min at 10s/block).
const syncProgressWindowBlocks = 60

func (c *NodeWidgetController) timeout1() {
	chainInfo, err := c.model.GetBlockchainInfo()
	if err != nil {
		return
	}
	lastBlockTime := time.Unix(chainInfo.LastBlockTime, 0)
	lastBlockHeight := chainInfo.LastBlockHeight

	glib.IdleAdd(func() bool {
		c.view.LabelLastBlockTime.SetText(lastBlockTime.Format("02 Jan 06 15:04:05 MST"))
		c.view.LabelLastBlockHeight.SetText(strconv.FormatInt(int64(lastBlockHeight), 10))

		nowUnix := time.Now().Unix()
		lastBlockTimeUnix := lastBlockTime.Unix()
		blocksLeft := (nowUnix - lastBlockTimeUnix) / 10
		c.view.LabelBlocksLeft.SetText(strconv.FormatInt(blocksLeft, 10))

		// Sync progress: 100% when up-to-date, 0% when syncProgressWindowBlocks behind (no genesis time).
		percentage := 1.0 - float64(blocksLeft)/float64(syncProgressWindowBlocks)
		if percentage < 0 {
			percentage = 0
		}
		if percentage > 1 {
			percentage = 1
		}
		c.view.ProgressBarSynced.SetFraction(percentage)
		c.view.ProgressBarSynced.SetText(fmt.Sprintf("%s %%", strconv.FormatFloat(percentage*100, 'f', 2, 64)))

		return false
	})
}

func (c *NodeWidgetController) timeout10() {
	chainInfo, err := c.model.GetBlockchainInfo()
	if err != nil {
		return
	}
	committeeInfo, _ := c.model.GetCommitteeInfo()
	consensusInfo, _ := c.model.GetConsensusInfo()
	nodeInfo, _ := c.model.GetNodeInfo()

	committeeSize := 0
	if committeeInfo != nil {
		committeeSize = len(committeeInfo.Validators)
	}
	inCommittee := consensusInfo != nil && len(consensusInfo.Instances) > 0

	var clockOffset time.Duration
	var clockOffsetErr error
	if nodeInfo != nil {
		clockOffset = time.Duration(nodeInfo.ClockOffset * float64(time.Second))
	}
	var numConnections, reachability string
	if nodeInfo != nil && nodeInfo.ConnectionInfo != nil {
		ci := nodeInfo.ConnectionInfo
		numConnections = fmt.Sprintf("%v (Inbound: %v, Outbound %v)",
			ci.Connections, ci.InboundConnections, ci.OutboundConnections)
		reachability = nodeInfo.Reachability
	}

	snapshot := nodeWidgetSnapshot{
		committeeSize:    committeeSize,
		committeeStake:   amount.Amount(chainInfo.CommitteePower),
		totalStake:       amount.Amount(chainInfo.TotalPower),
		activeValidators: chainInfo.ActiveValidators,
		numConnections:   numConnections,
		reachability:     reachability,
		inCommittee:      inCommittee,
		clockOffset:      clockOffset,
		clockOffsetErr:   clockOffsetErr,
	}

	glib.IdleAdd(func() bool {
		return c.applyTimeout10Snapshot(&snapshot)
	})
}

func (c *NodeWidgetController) applyTimeout10Snapshot(snapshot *nodeWidgetSnapshot) bool {
	styleContext, err := c.view.LabelClockOffset.GetStyleContext()
	if err != nil {
		logger.Error("failed to get style context", "err", err)

		return false
	}

	c.view.LabelClockOffset.SetTooltipText(
		"Difference between time of your machine and network time (NTP) " +
			"for synchronization.",
	)

	c.setClockOffset(styleContext, snapshot.clockOffset, snapshot.clockOffsetErr)

	c.view.LabelCommitteeSize.SetText(fmt.Sprintf("%v", snapshot.committeeSize))
	c.view.LabelActiveValidator.SetText(fmt.Sprintf("%v", snapshot.activeValidators))
	c.view.LabelCommitteeStake.SetText(snapshot.committeeStake.String())
	c.view.LabelTotalStake.SetText(snapshot.totalStake.String())
	c.setInCommittee(snapshot.inCommittee)
	c.view.LabelNumConnections.SetText(snapshot.numConnections)
	c.view.LabelReachability.SetText(snapshot.reachability)

	return false
}

func (c *NodeWidgetController) setClockOffset(styleContext *gtk.StyleContext, offset time.Duration, offsetErr error) {
	if offsetErr != nil {
		styleContext.AddClass("warning")
		c.view.LabelClockOffset.SetText("N/A")

		return
	}

	o := math.Round(offset.Seconds())
	if o == 0 {
		o = math.Abs(o) // To fix "-0 second(s)" issue
	}
	c.view.LabelClockOffset.SetText(fmt.Sprintf("%v second(s)", o))

	if offset > clockOutOfSyncThreshold || offset < -clockOutOfSyncThreshold {
		styleContext.AddClass("warning")

		return
	}
	styleContext.RemoveClass("warning")
}

func (c *NodeWidgetController) setInCommittee(inCommittee bool) {
	if inCommittee {
		c.view.LabelInCommittee.SetMarkup("<span foreground=\"#10c92f\">Yes</span>")

		return
	}

	c.view.LabelInCommittee.SetText("No")
}
