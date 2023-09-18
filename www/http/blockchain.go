package http

import (
	"net/http"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/sync/services"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func (s *Server) BlockchainHandler(w http.ResponseWriter, _ *http.Request) {
	res, err := s.blockchain.GetBlockchainInfo(s.ctx,
		&pactus.GetBlockchainInfoRequest{})
	if err != nil {
		s.writeError(w, err)
		return
	}

	tm := newTableMaker()
	tm.addRowBlockHash("Hash", res.LastBlockHash)
	tm.addRowInt("Height", int(res.LastBlockHeight))
	tm.addRowString("--- Committee", "---")
	tm.addRowAmount("Total Power", res.TotalPower)
	tm.addRowAmount("Committee Power", res.CommitteePower)
	for i, val := range res.CommitteeValidators {
		tm.addRowInt("--- Validator", i+1)
		s.writeValidatorTable(w, val)
	}

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
	tm.addRowTime("Started at", res.StartedAt)
	tm.addRowInt("Total Sent Bytes", int(res.TotalSentBytes))
	tm.addRowInt("Total Received Bytes", int(res.TotalReceivedBytes))
	tm.addRowString("Peers", "---")

	for i, p := range res.Peers {
		pid, _ := peer.IDFromBytes(p.PeerId)
		tm.addRowInt("-- Peer #", i+1)
		tm.addRowString("Status", peerset.StatusCode(p.Status).String())
		tm.addRowString("PeerID", pid.String())
		tm.addRowString("Services", services.Services(p.Services).String())
		for _, key := range p.ConsensusKeys {
			pub, _ := bls.PublicKeyFromString(key)
			tm.addRowString("  PublicKey", pub.String())
			tm.addRowValAddress("  Address", pub.Address().String())
		}
		tm.addRowString("Agent", p.Agent)
		tm.addRowString("Moniker", p.Moniker)
		tm.addRowString("LastSent", time.Unix(p.LastSent, 0).String())
		tm.addRowString("LastReceived", time.Unix(p.LastReceived, 0).String())
		tm.addRowBlockHash("Last block Hash", p.LastBlockHash)
		tm.addRowInt("Height", int(p.Height))
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

	s.writeHTML(w, tm.html())
}
