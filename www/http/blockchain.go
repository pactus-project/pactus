package http

import (
	"net/http"

	"github.com/fxamacker/cbor/v2"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/version"
	"github.com/zarbchain/zarb-go/www/capnp"
)

func (s *Server) BlockchainHandler(w http.ResponseWriter, r *http.Request) {
	res := s.server.GetBlockchainInfo(s.ctx, func(p capnp.ZarbServer_getBlockchainInfo_Params) error {
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
	res := s.server.GetNetworkInfo(s.ctx, func(p capnp.ZarbServer_getNetworkInfo_Params) error {
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
		peerID, _ := peer.IDFromString(id)
		peer := peerset.NewPeer(peerID)
		moniker, _ := p.Moniker()
		peer.UpdateMoniker(moniker)
		pubStr, _ := p.PublicKey()
		pub, _ := crypto.PublicKeyFromString(pubStr)
		peer.UpdatePublicKey(pub)
		bs, _ := p.NodeVersion()
		ver := version.Version{}
		if err := cbor.Unmarshal(bs, &ver); err != nil {
			s.writeError(w, err)
			return
		}
		peer.UpdateNodeVersion(ver)
		peer.UpdateInitialBlockDownload(p.InitialBlockDownload())
		peer.UpdateHeight(int(p.Height()))
		peer.UpdateInvalidMessage(int(p.InvalidMsg()))
		peer.UpdateReceivedMessage(int(p.ReceivedMsg()))
		peer.UpdateReceivedBytes(int(p.ReceivedBytes()))

		out.Peers[i] = peer
	}
	s.writeJSON(w, out)
}
