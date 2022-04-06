package http

import (
	"net/http"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/sync/peerset"
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
	tm.addRowBlockHash("Hash", hash)
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
		p := pl.At(i)

		id, _ := p.PeerID()
		pid, _ := peer.IDFromString(id)
		status := p.Status()
		moniker, _ := p.Moniker()
		pubStr, _ := p.PublicKey()
		lastSeen := time.Unix(int64(p.LastSeen()), 0)
		pub, _ := bls.PublicKeyFromString(pubStr)
		agent, _ := p.Agent()

		tm.addRowInt("Peer #", i+1)
		tm.addRowBytes("PeerID", []byte(pid))
		tm.addRowString("Status", peerset.StatusCode(status).String())
		tm.addRowBytes("PublicKey", pub.Bytes())
		tm.addRowString("Agent", agent)
		tm.addRowString("Moniker", moniker)
		tm.addRowString("LastSeen", lastSeen.String())
		tm.addRowInt("Height", int(p.Height()))
		tm.addRowInt("InvalidBundles", int(p.InvalidMessages()))
		tm.addRowInt("ReceivedBundles", int(p.ReceivedMessages()))
		tm.addRowInt("ReceivedBytes", int(p.ReceivedBytes()))
		tm.addRowInt("Flags", int(p.Flags()))
	}
	s.writeHTML(w, tm.html())
}
