package http

import (
	"net/http"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/www/capnp"
)

func (s *Server) BlockchainHandler(w http.ResponseWriter, r *http.Request) {
	res := s.capnp.GetBlockchainInfo(s.ctx, func(p capnp.ZarbServer_getBlockchainInfo_Params) error {
		return nil
	}).Result()

	st, err := res.Struct()
	if err != nil {
		s.writeError(w, err)
		return
	}
	hash, _ := st.LastBlockHash()

	tm := newTableMaker()
	tm.addRowBytes("Hash", hash)
	tm.addRowInt("Height", int(st.LastBlockHeight()))
	s.writeHTML(w, tm.html())
}

func (s *Server) NetworkHandler(w http.ResponseWriter, r *http.Request) {
	res := s.capnp.GetNetworkInfo(s.ctx, func(p capnp.ZarbServer_getNetworkInfo_Params) error {
		return nil
	}).Result()
	st, err := res.Struct()
	if err != nil {
		s.writeError(w, err)
		return
	}
	tm := newTableMaker()

	id, _ := st.PeerID()
	selfID, err := peer.Decode(id)
	if err != nil {
		s.writeError(w, err)
		return
	}
	tm.addRowBytes("Node ID", []byte(selfID))
	tm.addRowString("Peers", "---")

	pl, _ := st.Peers()
	for i := 0; i < pl.Len(); i++ {
		// p := pl.At(i)

		// id, _ := p.PeerID()
		// pid, _ := peer.IDFromString(id)
		// peer := peerset.NewPeer(pid)
		// status := p.Status()
		// moniker, _ := p.Moniker()
		// pubStr, _ := p.PublicKey()
		// pub, _ := bls.PublicKeyFromString(pubStr)
		// ver, _ := p.Agent()

		// tm.addRowBytes("PeerID", pid)
		// tm.addRowBytes("Status", peerset.StatusCode(status))
		// tm.addRowBytes("PublicKey", *pub)
		// tm.addRowBytes("Agent", ver)
		// tm.addRowBytes("Moniker", moniker)
		// tm.addRowBytes("Height", p.Height())
		// tm.addRowBytes("InvalidBundles", int(p.InvalidMessages()))
		// tm.addRowBytes("ReceivedBundles", int(p.ReceivedMessages()))
		// tm.addRowBytes("ReceivedBytes", int(p.ReceivedBytes()))
		// tm.addRowBytes("Flags", int(p.Flags()))
	}
	s.writeHTML(w, tm.html())
}
