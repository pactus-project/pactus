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

	scheduler.Every(refreshNetworkInterval).Do(ctx, func(ctx context.Context) {
		if gtkutil.IsWidgetShowing(&c.view.LabelConnectedPeers.Widget) {
			gtkutil.Logf("refreshing network info")
			c.refreshInfo(ctx)
		}

		if gtkutil.IsWidgetShowing(&c.view.ColViewPeers.Widget) {
			gtkutil.Logf("refreshing peer list")
			c.refreshList(ctx)
		}
	})

	c.refreshInfo(ctx)
	c.refreshList(ctx)

	return nil
}

func (c *NetworkWidgetController) refreshInfo(_ context.Context) {
	netInfo, err := c.model.GetNetworkInfo()
	if err != nil {
		return
	}

	gtkutil.IdleAddAsync(func() {
		c.view.LabelNetworkName.SetText(netInfo.GetNetworkName())
		c.view.LabelConnectedPeers.SetText(strconv.Itoa(int(netInfo.GetConnectedPeersCount())))
	})
}

func (c *NetworkWidgetController) refreshList(_ context.Context) {
	peersRes, err := c.model.ListPeers(false) // active peers only
	if err != nil {
		return
	}
	var rows []peerRow
	for no, peer := range peersRes.GetPeers() {
		rows = append(rows, peerRow{no: no + 1, peer: peer})
	}

	gtkutil.IdleAddAsync(func() {
		gtkutil.SyncListModel(c.lsPeers, rows)
	})
}
