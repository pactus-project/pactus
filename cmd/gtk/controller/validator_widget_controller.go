//go:build gtk

package controller

import (
	"context"
	"strconv"

	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/gopkg/scheduler"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/amount"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

// validatorRow represents a validator in the validator list.
type validatorRow struct {
	no  int
	val *pactus.ValidatorInfo
}

type ValidatorWidgetController struct {
	view         *view.ValidatorWidgetView
	model        *model.ValidatorModel
	lsValidators *gioutil.ListModel[validatorRow]
}

func NewValidatorWidgetController(
	view *view.ValidatorWidgetView, model *model.ValidatorModel,
) *ValidatorWidgetController {
	lsValidators := gioutil.NewListModel[validatorRow]()
	view.ColViewValidators.SetModel(gtk.NewSingleSelection(lsValidators))

	return &ValidatorWidgetController{
		view:         view,
		model:        model,
		lsValidators: lsValidators,
	}
}

func (c *ValidatorWidgetController) BuildView(ctx context.Context) error {
	gtkutil.IdleAddSync(func() {
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewValidators, "No", func(row validatorRow) string {
			return strconv.Itoa(row.no)
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewValidators, "Address", func(row validatorRow) string {
			return row.val.GetAddress()
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewValidators, "Stake", func(row validatorRow) string {
			return amount.Amount(row.val.GetStake()).String()
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewValidators, "Bonding Height", func(row validatorRow) string {
			return strconv.Itoa(int(row.val.GetLastBondingHeight()))
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewValidators, "Sortition Height", func(row validatorRow) string {
			return strconv.Itoa(int(row.val.GetLastSortitionHeight()))
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewValidators, "Unbonding Height", func(row validatorRow) string {
			return strconv.Itoa(int(row.val.GetUnbondingHeight()))
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewValidators, "Availability Score", func(row validatorRow) string {
			return gtkutil.AvailabilityScorePercent(row.val.GetAvailabilityScore())
		})
	})

	scheduler.Every(refreshValidatorsInterval).Do(ctx, func(ctx context.Context) {
		if gtkutil.IsWidgetShowing(&c.view.ColViewValidators.Widget) {
			gtkutil.Logf("refreshing validators")
			c.refresh(ctx)
		}
	})

	c.refresh(ctx)

	return nil
}

func (c *ValidatorWidgetController) refresh(_ context.Context) {
	vals, err := c.model.Validators()
	if err != nil {
		return
	}
	var rows []validatorRow
	for no, val := range vals {
		rows = append(rows, validatorRow{no: no + 1, val: val})
	}

	gtkutil.IdleAddAsync(func() {
		gtkutil.SyncListModel(c.lsValidators, rows)
	})
}
