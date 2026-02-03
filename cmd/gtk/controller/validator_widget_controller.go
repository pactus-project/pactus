//go:build gtk

package controller

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ezex-io/gopkg/scheduler"
	"github.com/gotk3/gotk3/glib"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/amount"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

// ValidatorWidgetModel is the interface used by the validator widget controller.
type ValidatorWidgetModel interface {
	Validators() ([]*pactus.ValidatorInfo, error)
}

type ValidatorWidgetController struct {
	view  *view.ValidatorWidgetView
	model ValidatorWidgetModel
}

func NewValidatorWidgetController(view *view.ValidatorWidgetView, m ValidatorWidgetModel) *ValidatorWidgetController {
	return &ValidatorWidgetController{view: view, model: m}
}

func (c *ValidatorWidgetController) Bind(ctx context.Context) error {
	scheduler.Every(ctx, 10*time.Second).Do(c.refresh)

	// Initial refresh.
	c.refresh()

	return nil
}

func (c *ValidatorWidgetController) refresh() {
	vals, err := c.model.Validators()
	if err != nil {
		return
	}

	glib.IdleAdd(func() bool {
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
					fmt.Sprintf("%.2f", val.GetAvailabilityScore()),
				},
			)
		}

		return false
	})
}
