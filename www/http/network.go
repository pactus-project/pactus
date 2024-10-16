package http

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"time"

	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/util"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func (s *Server) NetworkHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	onlyConnected := false

	onlyConnectedParam := r.URL.Query().Get("onlyConnected")
	if onlyConnectedParam == "true" {
		onlyConnected = true
	}

	res, err := s.network.GetNetworkInfo(ctx,
		&pactus.GetNetworkInfoRequest{
			OnlyConnected: onlyConnected,
		})
	if err != nil {
		s.writeError(w, err)

		return
	}

	tm := newTableMaker()
	tm.addRowString("Network Name", res.NetworkName)
	tm.addRowInt("Connected Peers Count", int(res.ConnectedPeersCount))
	metricToTable(tm, res.MetricInfo)

	tm.addRowString("Peers", "---")

	sort.Slice(res.ConnectedPeers, func(i, j int) bool {
		return res.ConnectedPeers[i].MetricInfo.TotalReceived.Bundles >
			res.ConnectedPeers[j].MetricInfo.TotalReceived.Bundles
	})

	for i, p := range res.ConnectedPeers {
		id, _ := hex.DecodeString(p.PeerId)
		pid, _ := lp2ppeer.IDFromBytes(id)
		tm.addRowInt("-- Peer #", i+1)
		tm.addRowString("Status", status.Status(p.Status).String())
		tm.addRowString("PeerID", pid.String())
		tm.addRowString("Services", service.Services(p.Services).String())
		tm.addRowString("Agent", p.Agent)
		tm.addRowString("Moniker", p.Moniker)
		tm.addRowString("Remote Address", p.Address)
		tm.addRowString("Direction", p.Direction)
		tm.addRowStrings("Protocols", p.Protocols)
		tm.addRowString("LastSent", time.Unix(p.LastSent, 0).String())
		tm.addRowString("LastReceived", time.Unix(p.LastReceived, 0).String())
		tm.addRowBlockHash("Last block Hash", p.LastBlockHash)
		tm.addRowInt("Height", int(p.Height))
		tm.addRowInt("TotalSessions", int(p.TotalSessions))
		tm.addRowInt("CompletedSessions", int(p.CompletedSessions))
		metricToTable(tm, p.MetricInfo)

		for _, key := range p.ConsensusKeys {
			pub, _ := bls.PublicKeyFromString(key)
			tm.addRowString("-- PublicKey", pub.String())
			tm.addRowValAddress("-- Address", pub.ValidatorAddress().String())
		}
	}
	s.writeHTML(w, tm.html())
}

func (s *Server) NodeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := s.network.GetNodeInfo(ctx,
		&pactus.GetNodeInfoRequest{})
	if err != nil {
		s.writeError(w, err)

		return
	}
	id, _ := hex.DecodeString(res.PeerId)
	sid, _ := lp2ppeer.IDFromBytes(id)
	tm := newTableMaker()
	tm.addRowString("Peer ID", sid.String())
	tm.addRowString("Agent", res.Agent)
	tm.addRowString("Moniker", res.Moniker)
	tm.addRowTime("Started at", int64(res.StartedAt))
	tm.addRowString("Reachability", res.Reachability)
	tm.addRowFloat64("Clock Offset", res.ClockOffset)
	tm.addRowInt("Services", int(res.Services))
	tm.addRowString("Services Names", res.ServicesNames)

	tm.addRowString("Connection Info", "---")
	tm.addRowInt("-- Total connections", int(res.ConnectionInfo.Connections))
	tm.addRowInt("-- Inbound connections", int(res.ConnectionInfo.InboundConnections))
	tm.addRowInt("-- Outbound connections", int(res.ConnectionInfo.OutboundConnections))

	tm.addRowString("Protocols", "---")
	for i, p := range res.Protocols {
		tm.addRowString(fmt.Sprint(i), p)
	}

	tm.addRowString("Local Addresses", "---")
	for i, la := range res.LocalAddrs {
		tm.addRowString(fmt.Sprint(i), la)
	}

	s.writeHTML(w, tm.html())
}

func metricToTable(tm *tableMaker, mi *pactus.MetricInfo) {
	printCounter := func(tm *tableMaker, name string, c *pactus.CounterInfo) {
		tm.addRowString(name,
			fmt.Sprintf("[%d, %s]", c.Bundles, util.FormatBytesToHumanReadable(c.Bytes)))
	}

	printSortedMap := func(tm *tableMaker, msgCounter map[int32]*pactus.CounterInfo) {
		keys := make([]int32, 0, len(msgCounter))
		for k := range msgCounter {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			return msgCounter[keys[i]].Bundles > msgCounter[keys[j]].Bundles
		})

		for _, key := range keys {
			printCounter(tm, message.Type(key).String(), msgCounter[key])
		}
	}

	printCounter(tm, "Total Invalid", mi.TotalInvalid)

	tm.addRowString("Sent Metric", "---")
	printCounter(tm, "Total Sent", mi.TotalSent)
	printSortedMap(tm, mi.MessageSent)

	tm.addRowString("Received Metric", "---")
	printCounter(tm, "Total Received", mi.TotalReceived)
	printSortedMap(tm, mi.MessageReceived)
}
