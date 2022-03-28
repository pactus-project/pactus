package http

import (
	"net/http"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
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
	data, _ := st.LastBlockHash()
	h, _ := hash.FromBytes(data)
	out := new(BlockchainResult)
	out.LastBlockHeight = int(st.LastBlockHeight())
	out.LastBlockHash = h
	s.writeJSON(w, out)
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
	out := new(NetworkResult)

	id, _ := st.PeerID()
	out.SelfID, err = peer.Decode(id)
	if err != nil {
		s.writeError(w, err)
		return
	}

	pl, _ := st.Peers()
	out.Peers = make([]*peerset.Peer, pl.Len())
	for i := 0; i < pl.Len(); i++ {
		p := pl.At(i)

		id, _ := p.PeerID()
		pid, _ := peer.IDFromString(id)
		peer := peerset.NewPeer(pid)
		status := p.Status()
		moniker, _ := p.Moniker()
		pubStr, _ := p.PublicKey()
		pub, _ := bls.PublicKeyFromString(pubStr)
		ver, _ := p.Agent()

		peer.PeerID = pid
		peer.Status = peerset.StatusCode(status)
		peer.PublicKey = *pub
		peer.Agent = ver
		peer.Moniker = moniker
		peer.Height = p.Height()
		peer.InvalidBundles = int(p.InvalidMessages())
		peer.ReceivedBundles = int(p.ReceivedMessages())
		peer.ReceivedBytes = int(p.ReceivedBytes())
		peer.Flags = int(p.Flags())

		out.Peers[i] = peer
	}
	s.writeJSON(w, out)
}
