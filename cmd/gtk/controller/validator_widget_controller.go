//go:build gtk

package controller

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/node"
	"github.com/pactus-project/pactus/types/validator"
)

type ValidatorWidgetController struct {
	view *view.ValidatorWidgetView
	node *node.Node

	timeoutID glib.SourceHandle
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewValidatorWidgetController(view *view.ValidatorWidgetView, nde *node.Node) *ValidatorWidgetController {
	return &ValidatorWidgetController{view: view, node: nde}
}

func (c *ValidatorWidgetController) Bind() error {
	// Reset lifecycle context.
	if c.cancel != nil {
		c.cancel()
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())

	c.timeoutID = glib.TimeoutAdd(10000, func() bool {
		c.refresh()

		return true
	})

	// Initial refresh.
	c.refresh()

	return nil
}

func (c *ValidatorWidgetController) refresh() {
	ctx := c.ctx
	go func() {
		if gtkutil.IsContextDone(ctx) {
			return
		}

		vals := make([]*validator.Validator, 0)
		for _, instance := range c.node.ConsManager().Instances() {
			addr := instance.ConsensusKey().ValidatorAddress()
			val, _ := c.node.State().ValidatorByAddress(addr)
			if val != nil {
				vals = append(vals, val)
			}
		}

		glib.IdleAdd(func() bool {
			if gtkutil.IsContextDone(ctx) {
				return false
			}

			c.view.ClearRows()
			for i, val := range vals {
				score := c.node.State().AvailabilityScore(val.Number())
				c.view.AppendRow(
					[]int{0, 1, 2, 3, 4, 5, 6},
					[]any{
						strconv.Itoa(i + 1),
						val.Address().String(),
						strconv.Itoa(int(val.Number())),
						val.Stake().String(),
						strconv.Itoa(int(val.LastBondingHeight())),
						strconv.Itoa(int(val.LastSortitionHeight())),
						fmt.Sprintf("%.2f", score),
					},
				)
			}

			return false
		})
	}()
}

func (c *ValidatorWidgetController) Cleanup() {
	if c.timeoutID != 0 {
		glib.SourceRemove(c.timeoutID)
		c.timeoutID = 0
	}
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
	}
}
