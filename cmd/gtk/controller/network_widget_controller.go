//go:build gtk

package controller

import (
	"context"
	"strconv"
	"time"

	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/ezex-io/gopkg/scheduler"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func peerDirectionString(d pactus.Direction) string {
	switch d {
	case pactus.Direction_DIRECTION_UNKNOWN:
		return "Unknown"
	case pactus.Direction_DIRECTION_INBOUND:
		return "Inbound"
	case pactus.Direction_DIRECTION_OUTBOUND:
		return "Outbound"

	default:
		return "Unknown"
	}
}

// peerRow represents a peer in the network peers list.
type peerRow struct {
	no   int
	peer *pactus.PeerInfo
}

type NetworkWidgetController struct {
	view    *view.NetworkWidgetView
	model   *model.NetworkModel
	lsPeers *gioutil.ListModel[peerRow]
}

func NewNetworkWidgetController(
	view *view.NetworkWidgetView, model *model.NetworkModel,
) *NetworkWidgetController {
	lsPeers := gioutil.NewListModel[peerRow]()
	view.ColViewPeers.SetModel(gtk.NewSingleSelection(lsPeers))

	return &NetworkWidgetController{
		view:    view,
		model:   model,
		lsPeers: lsPeers,
	}
}

func (c *NetworkWidgetController) BuildView(ctx context.Context) error {
	gtkutil.IdleAddSync(func() {
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewPeers, "No", func(row peerRow) string {
			return strconv.Itoa(row.no)
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewPeers, "Moniker", func(row peerRow) string {
			return row.peer.GetMoniker()
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewPeers, "Address", func(row peerRow) string {
			return row.peer.GetAddress()
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewPeers, "Peer ID", func(row peerRow) string {
			return row.peer.GetPeerId()
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewPeers, "Height", func(row peerRow) string {
			return strconv.Itoa(int(row.peer.GetHeight()))
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewPeers, "Agent", func(row peerRow) string {
			return row.peer.GetAgent()
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewPeers, "Direction", func(row peerRow) string {
			return peerDirectionString(row.peer.GetDirection())
		})
	})

	scheduler.Every(10*time.Second).Do(ctx, c.refresh)

	c.refresh(ctx)

	return nil
}

func (c *NetworkWidgetController) refresh(_ context.Context) {
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

		gtkutil.ClearListModel(c.lsPeers)

		// Add new peers to the list
		for i, peer := range peersRes.GetPeers() {
			row := peerRow{
				no:   i + 1,
				peer: peer,
			}

			c.lsPeers.Append(row)
		}
	})
}
