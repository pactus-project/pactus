//go:build gtk

package controller

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ezex-io/gopkg/scheduler"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/amount"
)

type ValidatorWidgetController struct {
	view  *view.ValidatorWidgetView
	model *model.ValidatorModel
}

func NewValidatorWidgetController(
	view *view.ValidatorWidgetView, model *model.ValidatorModel,
) *ValidatorWidgetController {
	return &ValidatorWidgetController{view: view, model: model}
}

func (c *ValidatorWidgetController) BuildView(ctx context.Context) error {
	scheduler.Every(10*time.Second).Do(ctx, func(context.Context) { c.refresh() })

	// Initial refresh.
	c.refresh()

	return nil
}

func (c *ValidatorWidgetController) refresh() {
	vals, err := c.model.Validators()
	if err != nil {
		return
	}

	gtkutil.IdleAddAsync(func() {
		c.view.ClearRows()
		for i, val := range vals {
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
					strconv.Itoa(int(val.GetUnbondingHeight())),
					gtkutil.AvailabilityScorePercent(val.GetAvailabilityScore()),
				},
			)
		}
	})
}
