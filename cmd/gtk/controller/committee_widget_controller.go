//go:build gtk

package controller

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/ezex-io/gopkg/scheduler"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/amount"
)

type CommitteeWidgetController struct {
	view  *view.CommitteeWidgetView
	model *model.CommitteeModel
}

func NewCommitteeWidgetController(
	view *view.CommitteeWidgetView, model *model.CommitteeModel,
) *CommitteeWidgetController {
	return &CommitteeWidgetController{view: view, model: model}
}

func (c *CommitteeWidgetController) BuildView(ctx context.Context) error {
	scheduler.Every(10*time.Second).Do(ctx, c.refresh)

	c.refresh(ctx)

	return nil
}

func (c *CommitteeWidgetController) refresh(_ context.Context) {
	res, err := c.model.GetCommitteeInfo()
	if err != nil {
		return
	}

	committeePowerStr := amount.Amount(res.CommitteePower).String()
	totalPowerStr := amount.Amount(res.TotalPower).String()

	// Protocol versions: map[int32]float64 -> "v1: 80%, v2: 20%"
	protocolLines := make([]string, 0, len(res.ProtocolVersions))
	for ver, pct := range res.ProtocolVersions {
		protocolLines = append(protocolLines, fmt.Sprintf("v%d: %.1f%%", ver, pct*100))
	}
	sort.Slice(protocolLines, func(i, j int) bool { return protocolLines[i] < protocolLines[j] })
	protocolStr := ""
	for i, s := range protocolLines {
		if i > 0 {
			protocolStr += ", "
		}
		protocolStr += s
	}
	if protocolStr == "" {
		protocolStr = "—"
	}

	gtkutil.IdleAddAsync(func() {
		c.view.LabelCommitteeSize.SetText(strconv.Itoa(int(res.CommitteeSize)))
		c.view.LabelCommitteePower.SetText(committeePowerStr)
		c.view.LabelTotalPower.SetText(totalPowerStr)
		c.view.LabelProtocolVersions.SetText(protocolStr)

		c.view.ClearRows()
		for i, val := range res.Validators {
			stakeStr := amount.Amount(val.GetStake()).String()
			c.view.AppendRow(
				[]int{0, 1, 2, 3, 4, 5, 6, 7},
				[]any{
					strconv.Itoa(i + 1),
					val.GetAddress(),
					strconv.Itoa(int(val.GetNumber())),
					stakeStr,
					strconv.Itoa(int(val.GetLastBondingHeight())),
					strconv.Itoa(int(val.GetLastSortitionHeight())),
					strconv.Itoa(int(val.GetProtocolVersion())),
					gtkutil.AvailabilityScorePercent(val.GetAvailabilityScore()),
				},
			)
		}
	})
}
