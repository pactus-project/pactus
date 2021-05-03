package http

import (
	"net/http"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto"
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
	out := new(BlockchainResult)
	out.Height = int(st.Height())
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
	out.ID, err = peer.Decode(id)
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
		pub, _ := crypto.PublicKeyFromString(pubStr)
		ver, _ := p.NodeVersion()

		peer.UpdateStatus(peerset.StatusCode(status))
		peer.UpdatePublicKey(pub)
		peer.UpdateMoniker(moniker)
		peer.UpdateNodeVersion(ver)
		peer.UpdateMoniker(moniker)
		peer.UpdatePublicKey(pub)
		peer.UpdateInitialBlockDownload(p.InitialBlockDownload())
		peer.UpdateHeight(int(p.Height()))
		peer.UpdateInvalidMessage(int(p.InvalidMessages()))
		peer.UpdateReceivedMessage(int(p.ReceivedMessages()))
		peer.UpdateReceivedBytes(int(p.ReceivedBytes()))

		out.Peers[i] = peer
	}
	s.writeJSON(w, out)
}
