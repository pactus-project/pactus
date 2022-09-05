package grpc

import (
	"context"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/sync"
	"github.com/zarbchain/zarb-go/util/logger"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

type networkServer struct {
	sync   sync.Synchronizer
	logger *logger.Logger
}

func (s *networkServer) GetNetworkInfo(ctx context.Context,
	request *zarb.NetworkInfoRequest) (*zarb.NetworkInfoResponse, error) {
	// Create response peers
	rps := make([]*zarb.PeerInfo, int32(len(s.sync.Peers())))

	// cast to response peers from synced peers
	for i, peer := range s.sync.Peers() {
		rps[i] = new(zarb.PeerInfo)
		p := rps[i]

		bs, err := cbor.Marshal(peer.Agent)
		if err != nil {
			s.logger.Error("couldn't marshal agent", "err", err)
			continue
		}
		p.Agent = string(bs)

		p.PeerId = []byte(peer.PeerID)
		p.Moniker = peer.Moniker
		p.Agent = peer.Agent
		p.PublicKey = peer.PublicKey.String()
		p.Flags = int32(peer.Flags)
		p.Height = peer.Height
		p.ReceivedMessages = int32(peer.ReceivedBundles)
		p.InvalidMessages = int32(peer.InvalidBundles)
		p.ReceivedBytes = int32(peer.ReceivedBytes)
	}

	return &zarb.NetworkInfoResponse{
		SelfId: []byte(s.sync.SelfID()),
		Peers:  rps,
	}, nil
}
