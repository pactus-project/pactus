//go:build gtk

package controller

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/gopkg/scheduler"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/protocol"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

// committeeRow represents a validator in the committee list.
type committeeRow struct {
	no  int
	val *pactus.ValidatorInfo
}

type CommitteeWidgetController struct {
	view  *view.CommitteeWidgetView
	model *model.CommitteeModel

	lsMembers *gioutil.ListModel[committeeRow]
}

func NewCommitteeWidgetController(
	view *view.CommitteeWidgetView, model *model.CommitteeModel,
) *CommitteeWidgetController {
	lsMembers := gioutil.NewListModel[committeeRow]()
	view.ColViewMembers.SetModel(gtk.NewSingleSelection(lsMembers))

	return &CommitteeWidgetController{
		view:      view,
		model:     model,
		lsMembers: lsMembers,
	}
}

func (c *CommitteeWidgetController) BuildView(ctx context.Context) error {
	gtkutil.IdleAddSync(func() {
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewMembers, "No", func(row committeeRow) string {
			return strconv.Itoa(row.no)
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewMembers, "Address", func(row committeeRow) string {
			return row.val.GetAddress()
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewMembers, "Stake", func(row committeeRow) string {
			return amount.Amount(row.val.GetStake()).String()
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewMembers, "Bonding Height", func(row committeeRow) string {
			return strconv.Itoa(int(row.val.GetLastBondingHeight()))
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewMembers, "Sortition Height", func(row committeeRow) string {
			return strconv.Itoa(int(row.val.GetLastSortitionHeight()))
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewMembers, "Protocol Version", func(row committeeRow) string {
			return strconv.Itoa(int(row.val.GetProtocolVersion()))
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewMembers, "Availability Score", func(row committeeRow) string {
			return gtkutil.AvailabilityScorePercent(row.val.GetAvailabilityScore())
		})
	})

	scheduler.Every(refreshCommitteeInterval).Do(ctx, func(ctx context.Context) {
		if gtkutil.IsWidgetShowing(&c.view.Box.Widget) {
			gtkutil.Logf("refreshing committee")
			c.refresh(ctx)
		}
	})

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
	for ver, percentage := range res.ProtocolVersions {
		protocolLines = append(protocolLines, fmt.Sprintf("%s: %.2f%%", protocol.Version(ver), percentage))
	}
	slices.SortFunc(protocolLines, strings.Compare)
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

	var rows []committeeRow
	for no, val := range res.Validators {
		rows = append(rows, committeeRow{no: no + 1, val: val})
	}

	gtkutil.IdleAddAsync(func() {
		c.view.LabelCommitteeSize.SetText(strconv.Itoa(int(res.CommitteeSize)))
		c.view.LabelCommitteePower.SetText(committeePowerStr)
		c.view.LabelTotalPower.SetText(totalPowerStr)
		c.view.LabelProtocolVersions.SetText(protocolStr)

		gtkutil.SyncListModel(c.lsMembers, rows)
	})
}
