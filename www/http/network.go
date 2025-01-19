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

	tmk := newTableMaker()
	tmk.addRowString("Network Name", res.NetworkName)
	tmk.addRowInt("Connected Peers Count", int(res.ConnectedPeersCount))
	metricToTable(tmk, res.MetricInfo)

	tmk.addRowString("Peers", "---")

	sort.Slice(res.ConnectedPeers, func(i, j int) bool {
		return res.ConnectedPeers[i].MetricInfo.TotalReceived.Bundles >
			res.ConnectedPeers[j].MetricInfo.TotalReceived.Bundles
	})

	for index, peer := range res.ConnectedPeers {
		id, _ := hex.DecodeString(peer.PeerId)
		pid, _ := lp2ppeer.IDFromBytes(id)
		tmk.addRowInt("-- Peer #", index+1)
		tmk.addRowString("Status", status.Status(peer.Status).String())
		tmk.addRowString("PeerID", pid.String())
		tmk.addRowString("Services", service.Services(peer.Services).String())
		tmk.addRowString("Agent", peer.Agent)
		tmk.addRowString("Moniker", peer.Moniker)
		tmk.addRowString("Remote Address", peer.Address)
		tmk.addRowString("Direction", peer.Direction)
		tmk.addRowStrings("Protocols", peer.Protocols)
		tmk.addRowString("LastSent", time.Unix(peer.LastSent, 0).String())
		tmk.addRowString("LastReceived", time.Unix(peer.LastReceived, 0).String())
		tmk.addRowBlockHash("Last block Hash", peer.LastBlockHash)
		tmk.addRowInt("Height", int(peer.Height))
		tmk.addRowInt("TotalSessions", int(peer.TotalSessions))
		tmk.addRowInt("CompletedSessions", int(peer.CompletedSessions))
		metricToTable(tmk, peer.MetricInfo)

		for _, key := range peer.ConsensusKeys {
			pub, _ := bls.PublicKeyFromString(key)
			tmk.addRowString("-- PublicKey", pub.String())
			tmk.addRowValAddress("-- Address", pub.ValidatorAddress().String())
		}
	}
	s.writeHTML(w, tmk.html())
}

func (s *Server) NodeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := s.network.GetNodeInfo(ctx,
		&pactus.GetNodeInfoRequest{})
	if err != nil {
		s.writeError(w, err)

		return
	}
	pid, _ := hex.DecodeString(res.PeerId)
	sid, _ := lp2ppeer.IDFromBytes(pid)
	tmk := newTableMaker()
	tmk.addRowString("Peer ID", sid.String())
	tmk.addRowString("Agent", res.Agent)
	tmk.addRowString("Moniker", res.Moniker)
	tmk.addRowTime("Started at", int64(res.StartedAt))
	tmk.addRowString("Reachability", res.Reachability)
	tmk.addRowFloat64("Clock Offset", res.ClockOffset)
	tmk.addRowInt("Services", int(res.Services))
	tmk.addRowString("Services Names", res.ServicesNames)

	tmk.addRowString("Connection Info", "---")
	tmk.addRowInt("-- Total connections", int(res.ConnectionInfo.Connections))
	tmk.addRowInt("-- Inbound connections", int(res.ConnectionInfo.InboundConnections))
	tmk.addRowInt("-- Outbound connections", int(res.ConnectionInfo.OutboundConnections))

	tmk.addRowString("Protocols", "---")
	for i, p := range res.Protocols {
		tmk.addRowString(fmt.Sprint(i), p)
	}

	tmk.addRowString("ZeroMQ Publishers", "---")
	for i, p := range res.ZmqPublishers {
		tmk.addRowString(fmt.Sprint(i), fmt.Sprintf("%s - %s - %d", p.Topic, p.Address, p.Hwm))
	}

	tmk.addRowString("Local Addresses", "---")
	for i, la := range res.LocalAddrs {
		tmk.addRowString(fmt.Sprint(i), la)
	}

	s.writeHTML(w, tmk.html())
}

func metricToTable(tmk *tableMaker, metricInfo *pactus.MetricInfo) {
	printCounter := func(tmk *tableMaker, name string, counterInfo *pactus.CounterInfo) {
		tmk.addRowString(name,
			fmt.Sprintf("[%d, %s]", counterInfo.Bundles, util.FormatBytesToHumanReadable(counterInfo.Bytes)))
	}

	printSortedMap := func(tmk *tableMaker, msgCounter map[int32]*pactus.CounterInfo) {
		keys := make([]int32, 0, len(msgCounter))
		for k := range msgCounter {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			return msgCounter[keys[i]].Bundles > msgCounter[keys[j]].Bundles
		})

		for _, key := range keys {
			printCounter(tmk, message.Type(key).String(), msgCounter[key])
		}
	}

	printCounter(tmk, "Total Invalid", metricInfo.TotalInvalid)

	tmk.addRowString("Sent Metric", "---")
	printCounter(tmk, "Total Sent", metricInfo.TotalSent)
	printSortedMap(tmk, metricInfo.MessageSent)

	tmk.addRowString("Received Metric", "---")
	printCounter(tmk, "Total Received", metricInfo.TotalReceived)
	printSortedMap(tmk, metricInfo.MessageReceived)
}
