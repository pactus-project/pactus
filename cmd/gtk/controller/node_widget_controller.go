//go:build gtk

package controller

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/pactus-project/gopkg/scheduler"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/amount"
)

// clockOutOfSyncThreshold is the clock offset above which we show a warning.
const clockOutOfSyncThreshold = 5 * time.Second

// syncProgressWindowBlocks is the number of blocks behind that maps to 0% sync progress (~10 min at 10s/block).
// Used as a fallback when the server does not provide SyncProgress (backward compatibility).
const syncProgressWindowBlocks = 60

type NodeWidgetController struct {
	view  *view.NodeWidgetView
	model *model.NodeModel
}

func NewNodeWidgetController(view *view.NodeWidgetView, model *model.NodeModel) *NodeWidgetController {
	return &NodeWidgetController{view: view, model: model}
}

// BuildView builds the node widget. connectionLabel is either "Remote address" or "Working directory";
// connectionValue is the remote address or working directory path respectively.
func (c *NodeWidgetController) BuildView(ctx context.Context, connectionLabel, connectionValue string) error {
	nodeInfo, err := c.model.GetNodeInfo()
	if err != nil {
		return err
	}
	chainInfo, err := c.model.GetBlockchainInfo()
	if err != nil {
		return err
	}

	gtkutil.IdleAddSync(func() {
		c.view.LabelConnectionType.SetText(connectionLabel + ":")
		c.view.LabelConnectionValue.SetText(connectionValue)
		c.view.LabelNetwork.SetText(nodeInfo.NetworkName)
		c.view.LabelNetworkID.SetText(nodeInfo.PeerId)
		c.view.LabelAgent.SetText(nodeInfo.Agent)
		c.view.LabelMoniker.SetText(nodeInfo.Moniker)
		c.view.LabelIsPrune.SetText(strconv.FormatBool(chainInfo.IsPruned))
	})

	scheduler.Every(refreshNodeProgressInterval).Do(ctx, func(context.Context) { c.timeoutProgress() })
	scheduler.Every(refreshNodeInfoInterval).Do(ctx, func(context.Context) { c.timeoutInfo() })

	// Initial refresh.
	c.timeoutProgress()
	c.timeoutInfo()

	return nil
}

func (c *NodeWidgetController) timeoutProgress() {
	chainInfo, err := c.model.GetBlockchainInfo()
	if err != nil {
		return
	}
	lastBlockTime := time.Unix(chainInfo.LastBlockTime, 0)

	gtkutil.IdleAddSync(func() {
		c.view.LabelLastBlockTime.SetText(lastBlockTime.Format("02 Jan 06 15:04:05 MST"))
		c.view.LabelLastBlockHeight.SetText(strconv.FormatInt(int64(chainInfo.LastBlockHeight), 10))

		percentage := chainInfo.SyncProgress
		if percentage == 0 {
			// Backward compatibility: old servers do not send SyncProgress.
			// Fall back to the local blocks-left heuristic.
			// TODO: remove this fallback in a future release.
			blocksLeft := (time.Now().Unix() - lastBlockTime.Unix()) / 10
			if blocksLeft > syncProgressWindowBlocks {
				blocksLeft = syncProgressWindowBlocks
			}
			percentage = 1.0 - float64(blocksLeft)/float64(syncProgressWindowBlocks)
			c.view.LabelBlocksLeft.SetText(strconv.FormatInt(blocksLeft, 10))
		} else {
			c.view.LabelBlocksLeft.SetText(strconv.FormatInt(chainInfo.BlocksLeft, 10))
		}
		c.view.ProgressBarSynced.SetFraction(percentage)
		c.view.ProgressBarSynced.SetText(fmt.Sprintf("%s %%", strconv.FormatFloat(percentage*100, 'f', 2, 64)))
	})
}

func (c *NodeWidgetController) timeoutInfo() {
	chainInfo, err := c.model.GetBlockchainInfo()
	if err != nil {
		return
	}

	nodeInfo, err := c.model.GetNodeInfo()
	if err != nil {
		return
	}

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
	totalStake := amount.Amount(chainInfo.TotalPower)

	gtkutil.IdleAddSync(func() {
		c.view.LabelClockOffset.SetTooltipText(
			"Difference between time of your machine and network time (NTP) " +
				"for synchronization.",
		)

		c.setClockOffset(clockOffset, clockOffsetErr)

		c.view.LabelActiveValidator.SetText(fmt.Sprintf("%v", chainInfo.ActiveValidators))
		c.view.LabelTotalPower.SetText(totalStake.String())
		c.view.LabelAverageScore.SetText(fmt.Sprintf("%.2f", chainInfo.AverageScore))
		c.view.LabelNumConnections.SetText(numConnections)
		c.view.LabelReachability.SetText(reachability)

		c.setInCommittee(chainInfo.InCommittee)
	})
}

func (c *NodeWidgetController) setClockOffset(offset time.Duration, offsetErr error) {
	if offsetErr != nil {
		c.view.LabelClockOffset.AddCSSClass("warning")
		c.view.LabelClockOffset.SetText("N/A")

		return
	}

	o := math.Round(offset.Seconds())
	if o == 0 {
		o = math.Abs(o) // To fix "-0 second(s)" issue
	}
	c.view.LabelClockOffset.SetText(fmt.Sprintf("%v second(s)", o))

	if offset > clockOutOfSyncThreshold || offset < -clockOutOfSyncThreshold {
		c.view.LabelClockOffset.AddCSSClass("warning")

		return
	}

	c.view.LabelClockOffset.RemoveCSSClass("warning")
}

func (c *NodeWidgetController) setInCommittee(inCommittee bool) {
	if inCommittee {
		c.view.LabelInCommittee.SetMarkup("<span foreground=\"#10c92f\">Yes</span>")

		return
	}

	c.view.LabelInCommittee.SetText("No")
}
