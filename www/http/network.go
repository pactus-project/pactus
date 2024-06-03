package http

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func (s *Server) NetworkHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if s.enableAuth {
		user, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)

			return
		}

		ctx = s.basicAuth(ctx, user, password)
	}

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

	printSortedMap := func(tm *tableMaker, stats map[int32]int64) {
		keys := make([]int32, 0, len(stats))
		for k := range stats {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			return stats[keys[i]] > stats[keys[j]]
		})

		for _, key := range keys {
			tm.addRowInt(message.Type(key).String(), int(stats[key]))
		}
	}

	tm := newTableMaker()
	tm.addRowString("Network Name", res.NetworkName)
	tm.addRowInt("Total Sent Bytes", int(res.TotalSentBytes))
	tm.addRowInt("Total Received Bytes", int(res.TotalReceivedBytes))
	tm.addRowInt("Connected Peers Count", int(res.ConnectedPeersCount))

	tm.addRowString("ReceivedBytes", "---")
	printSortedMap(tm, res.ReceivedBytes)

	tm.addRowString("SentBytes", "---")
	printSortedMap(tm, res.SentBytes)

	tm.addRowString("Peers", "---")

	sort.Slice(res.ConnectedPeers, func(i, j int) bool {
		return res.ConnectedPeers[i].ReceivedBundles > res.ConnectedPeers[j].ReceivedBundles
	})

	for i, p := range res.ConnectedPeers {
		pid, _ := lp2ppeer.IDFromBytes(p.PeerId)
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
		tm.addRowInt("InvalidBundles", int(p.InvalidBundles))
		tm.addRowInt("ReceivedBundles", int(p.ReceivedBundles))

		tm.addRowString("ReceivedBytes", "---")
		printSortedMap(tm, p.ReceivedBytes)

		tm.addRowString("SentBytes", "---")
		printSortedMap(tm, p.SentBytes)

		for _, key := range p.ConsensusKeys {
			pub, _ := bls.PublicKeyFromString(key)
			tm.addRowString("  PublicKey", pub.String())
			tm.addRowValAddress("  Address", pub.ValidatorAddress().String())
		}
	}
	s.writeHTML(w, tm.html())
}

func (s *Server) NodeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if s.enableAuth {
		user, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)

			return
		}

		ctx = s.basicAuth(ctx, user, password)
	}

	res, err := s.network.GetNodeInfo(ctx,
		&pactus.GetNodeInfoRequest{})
	if err != nil {
		s.writeError(w, err)

		return
	}

	sid, _ := lp2ppeer.IDFromBytes(res.PeerId)
	tm := newTableMaker()
	tm.addRowString("Peer ID", sid.String())
	tm.addRowString("Agent", res.Agent)
	tm.addRowString("Moniker", res.Moniker)
	tm.addRowTime("Started at", int64(res.StartedAt))
	tm.addRowString("Reachability", res.Reachability)
	tm.addRowInts("Services", res.Services)
	tm.addRowStrings("Services Names", res.ServicesNames)

	tm.addRowString("Protocols", "---")
	for i, p := range res.Protocols {
		tm.addRowString(fmt.Sprint(i), p)
	}

	tm.addRowString("LocalAddress", "---")
	for i, la := range res.LocalAddrs {
		tm.addRowString(fmt.Sprint(i), la)
	}

	s.writeHTML(w, tm.html())
}
