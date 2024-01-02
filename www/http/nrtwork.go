package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/sync/peerset/service"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func (s *Server) PeersHandler(w http.ResponseWriter, _ *http.Request) {
	res, err := s.network.GetPeersInfo(s.ctx,
		&pactus.GetPeersInfoRequest{})
	if err != nil {
		s.writeError(w, err)
		return
	}

	tm := newTableMaker()
	tm.addRowTime("Started at", res.StartedAt)
	tm.addRowInt("Total Sent Bytes", int(res.TotalSentBytes))
	tm.addRowInt("Total Received Bytes", int(res.TotalReceivedBytes))
	tm.addRowString("Peers", "---")

	for i, p := range res.Peers {
		pid, _ := peer.IDFromBytes(p.PeerId)
		tm.addRowInt("-- Peer #", i+1)
		tm.addRowString("Status", peerset.StatusCode(p.Status).String())
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
		tm.addRowInt("InvalidBundles", int(p.InvalidMessages))
		tm.addRowInt("ReceivedBundles", int(p.ReceivedMessages))
		tm.addRowString("ReceivedBytes", "---")
		for key, value := range p.ReceivedBytes {
			tm.addRowInt(message.Type(key).String(), int(value))
		}
		tm.addRowString("SentBytes", "---")
		for key, value := range p.SentBytes {
			tm.addRowInt(message.Type(key).String(), int(value))
		}
		for _, key := range p.ConsensusKeys {
			pub, _ := bls.PublicKeyFromString(key)
			tm.addRowString("  PublicKey", pub.String())
			tm.addRowValAddress("  Address", pub.ValidatorAddress().String())
		}
	}
	s.writeHTML(w, tm.html())
}

func (s *Server) NodeHandler(w http.ResponseWriter, _ *http.Request) {
	res, err := s.network.GetNodeInfo(s.ctx,
		&pactus.GetNodeInfoRequest{})
	if err != nil {
		s.writeError(w, err)
		return
	}

	sid, _ := peer.IDFromBytes(res.PeerId)
	tm := newTableMaker()
	tm.addRowString("Peer ID", sid.String())
	tm.addRowString("Agent", res.Agent)
	tm.addRowString("Moniker", res.Moniker)
	tm.addRowString("Reachability", res.Reachability)
	tm.addRowStrings("Addrs", res.Addrs)

	s.writeHTML(w, tm.html())
}

func (s *Server) NetworkHandler(w http.ResponseWriter, _ *http.Request) {
	res, err := s.network.GetNetworkInfo(s.ctx,
		&pactus.GetNetworkInfoRequest{})
	if err != nil {
		s.writeError(w, err)
		return
	}

	tm := newTableMaker()
	tm.addRowInt("Version", int(res.ProtocolVersion))
	tm.addRowInt("Number of connected peers", int(res.ConnectedPeers))
	tm.addRowInt("Network DHT size", int(res.NetworkDhtSize))

	tm.addRowString("Protocols", "---")
	for i, p := range res.Protocols {
		tm.addRowString(fmt.Sprint(i), p)
	}

	tm.addRowString("LocalAddress", "---")
	for i, la := range res.LocalAddress {
		tm.addRowString(fmt.Sprint(i), la)
	}

	s.writeHTML(w, tm.html())
}
