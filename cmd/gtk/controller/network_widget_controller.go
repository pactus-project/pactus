//go:build gtk

package controller

import (
	"context"
	"strconv"
	"time"

	"github.com/ezex-io/gopkg/scheduler"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func peerDirectionString(d pactus.Direction) string {
	switch d {
	case pactus.Direction_DIRECTION_INBOUND:
		return "Inbound"
	case pactus.Direction_DIRECTION_OUTBOUND:
		return "Outbound"
	default:
		return "Unknown"
	}
}

type NetworkWidgetController struct {
	view  *view.NetworkWidgetView
	model *model.NetworkModel
}

func NewNetworkWidgetController(
	view *view.NetworkWidgetView, model *model.NetworkModel,
) *NetworkWidgetController {
	return &NetworkWidgetController{view: view, model: model}
}

func (c *NetworkWidgetController) BuildView(ctx context.Context) error {
	scheduler.Every(10*time.Second).Do(ctx, func(context.Context) { c.refresh() })

	c.refresh()

	return nil
}

func (c *NetworkWidgetController) refresh() {
	netInfo, err := c.model.GetNetworkInfo()
	if err != nil {
		return
	}

	peersRes, err := c.model.ListPeers(false) // active peers only
	if err != nil {
		return
	}

	gtkutil.IdleAddAsync(func() {
		c.view.LabelNetworkName.SetText(netInfo.GetNetworkName())
		c.view.LabelConnectedPeers.SetText(strconv.Itoa(int(netInfo.GetConnectedPeersCount())))

		c.view.ClearRows()
		for i, peer := range peersRes.GetPeers() {
			c.view.AppendRow(
				[]int{0, 1, 2, 3, 4, 5, 6},
				[]any{
					strconv.Itoa(i + 1),
					peer.GetMoniker(),
					peer.GetAddress(),
					peer.GetPeerId(),
					strconv.Itoa(int(peer.GetHeight())),
					peer.GetAgent(),
					peerDirectionString(peer.GetDirection()),
				},
			)
		}
	})
}
