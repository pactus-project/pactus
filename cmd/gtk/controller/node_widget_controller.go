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
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/node"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util/logger"
)

type NodeWidgetController struct {
	view *view.NodeWidgetView
	node *node.Node

	genesisTime time.Time
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

func NewNodeWidgetController(view *view.NodeWidgetView, nde *node.Node) *NodeWidgetController {
	return &NodeWidgetController{view: view, node: nde, genesisTime: nde.State().Genesis().GenesisTime()}
}

func (c *NodeWidgetController) Bind(ctx context.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	c.view.LabelWorkingDirectory.SetText(cwd)
	c.view.LabelNetwork.SetText(c.node.State().Genesis().ChainType().String())
	c.view.LabelNetworkID.SetText(c.node.Network().SelfID().String())
	c.view.LabelMoniker.SetText(c.node.Sync().Moniker())
	c.view.LabelIsPrune.SetText(strconv.FormatBool(c.node.State().ChainInfo().IsPruned))

	c.view.ConnectSignals(map[string]any{})

	scheduler.Every(ctx, time.Second).Do(c.timeout1)
	scheduler.Every(ctx, 10*time.Second).Do(c.timeout10)

	// Initial refresh.
	c.timeout1()
	c.timeout10()

	return nil
}

func (c *NodeWidgetController) timeout1() {
	lastBlockTime := c.node.State().LastBlockTime()
	lastBlockHeight := c.node.State().LastBlockHeight()

	glib.IdleAdd(func() bool {
		c.view.LabelLastBlockTime.SetText(lastBlockTime.Format("02 Jan 06 15:04:05 MST"))
		c.view.LabelLastBlockHeight.SetText(strconv.FormatInt(int64(lastBlockHeight), 10))

		nowUnix := time.Now().Unix()
		lastBlockTimeUnix := lastBlockTime.Unix()
		genTimeUnix := c.genesisTime.Unix()

		percentage := float64(lastBlockTimeUnix-genTimeUnix) / float64(nowUnix-genTimeUnix)
		c.view.ProgressBarSynced.SetFraction(percentage)
		c.view.ProgressBarSynced.SetText(fmt.Sprintf("%s %%", strconv.FormatFloat(percentage*100, 'f', 2, 64)))

		blocksLeft := (nowUnix - lastBlockTimeUnix) / 10
		c.view.LabelBlocksLeft.SetText(strconv.FormatInt(blocksLeft, 10))

		return false
	})
}

func (c *NodeWidgetController) timeout10() {
	info := c.node.State().ChainInfo()
	offset, offsetErr := c.node.Sync().ClockOffset()

	snapshot := nodeWidgetSnapshot{
		committeeSize:    c.node.State().Params().CommitteeSize,
		committeeStake:   amount.Amount(info.CommitteePower),
		totalStake:       amount.Amount(info.TotalPower),
		activeValidators: info.ActiveValidators,
		numConnections: fmt.Sprintf("%v (Inbound: %v, Outbound %v)",
			c.node.Network().NumConnectedPeers(),
			c.node.Network().NumInbound(),
			c.node.Network().NumOutbound(),
		),
		reachability:   c.node.Network().ReachabilityStatus(),
		inCommittee:    c.node.ConsManager().HasActiveInstance(),
		clockOffset:    offset,
		clockOffsetErr: offsetErr,
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

	if c.node.Sync().IsClockOutOfSync() {
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
