package http

import (
	"net/http"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/peerset"
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
		tm.addRowString("Public Key", val.PublicKey)
		tm.addRowValAddress("Address", val.Address)
		tm.addRowInt("Number", int(val.Number))
		tm.addRowInt("Sequence", int(val.Sequence))
		tm.addRowAmount("Stake", val.Stake)
		tm.addRowInt("LastBondingHeight", int(val.LastBondingHeight))
		tm.addRowInt("LastJoinedHeight", int(val.LastJoinedHeight))
		tm.addRowInt("UnbondingHeight", int(val.UnbondingHeight))
		tm.addRowBytes("Hash", val.Hash)
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

	sid, _ := peer.IDFromBytes(res.SelfId)
	tm := newTableMaker()
	tm.addRowString("Peer ID", sid.String())
	tm.addRowString("Peers", "---")

	for i, p := range res.Peers {
		pid, _ := peer.IDFromBytes(p.PeerId)
		tm.addRowInt("-- Peer #", i+1)
		tm.addRowString("PeerID", pid.String())
		tm.addRowString("Status", peerset.StatusCode(p.Status).String())
		for _, key := range p.Keys {
			pub, _ := bls.PublicKeyFromString(key)
			tm.addRowString("  PublicKey", pub.String())
			tm.addRowValAddress("  Address", pub.Address().String())
		}
		tm.addRowString("Agent", p.Agent)
		tm.addRowString("Moniker", p.Moniker)
		tm.addRowString("LastSeen", time.Unix(p.LastSeen, 0).String())
		tm.addRowInt("Height", int(p.Height))
		tm.addRowInt("InvalidBundles", int(p.InvalidMessages))
		tm.addRowInt("ReceivedBundles", int(p.ReceivedMessages))
		tm.addRowInt("ReceivedBytes", int(p.ReceivedBytes))
		tm.addRowInt("SendSuccess", int(p.SendSuccess))
		tm.addRowInt("SendFailed", int(p.SendFailed))
		tm.addRowInt("Flags", int(p.Flags))
	}
	s.writeHTML(w, tm.html())
}
